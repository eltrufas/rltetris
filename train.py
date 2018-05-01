from gymenv import TetrisEnv
from baselines import deepq

if __name__ == '__main__':
    env = TetrisEnv()
    model = deepq.models.cnn_to_mlp(
        convs=[(32, 3, 1), (64, 3, 1), (64, 3, 1)],
        hiddens=[4096],
    )
    act = deepq.learn(
        env,
        q_func=model,
        lr=1e-4,
        max_timesteps=2000000,
        buffer_size=10000,
        exploration_fraction=0.1,
        exploration_final_eps=0.01,
        train_freq=4,
        learning_starts=10000,
        target_network_update_freq=1000,
        gamma=0.99,
        prioritized_replay=False
    )
    print("Saving model to tetris.pkl")
    act.save("tetris.pkl")

from gymenv import TetrisEnv
from baselines import deepq

if __name__ == '__main__':
    env = TetrisEnv()
    model = deepq.models.mlp([4096])
    act = deepq.learn(
        env,
        q_func=model,
        lr=1e-3,
        max_timesteps=1000000,
        buffer_size=50000,
        exploration_fraction=0.1,
        gamma=0.99,
        exploration_final_eps=0.02,
        print_freq=10
    )
    print("Saving model to tetris.pkl")
    act.save("tetris.pkl")

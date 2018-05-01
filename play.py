from baselines import deepq
from gym import TetrisEnv

env = TetrisEnv(5050)
act = deepq.load("tetris.pkl")

while True:
    obs, done = env.reset(), False
    ep_reward = 0

    while not done:
        obs, reward, done, _ = env.step(act(obs[None])[0])

        ep_reward += reward

import socket
import struct
import numpy as np
from gym import Env, spaces


def reshape_obs(arr):
    return arr.reshape((10, 22, 3))


class TetrisEnv(Env):
    def __init__(self, port=3030):
        self.action_space = spaces.Discrete(7)
        self.sock = None
        self.port = port

        self.sock = socket.socket()
        self.sock.connect(('localhost', port))

        self.observation_space = spaces.Box(low=0, high=1, shape=(10, 22, 3), dtype=np.uint8)

        self.sock = None

    def reset(self):
        if self.sock:
            self.sock.close()
        self.sock = socket.socket()
        self.sock.connect(('localhost', self.port))

        obs = list(self.sock.recv(4096))
        obs = np.array(obs, dtype=np.int8)
        obs = reshape_obs(obs)

        return obs

    def step(self, action):
        self.sock.send(bytes([action]))

        data = self.sock.recv(4096)
        reward, done = struct.unpack('dB', data)
        done = bool(done)
        self.sock.send(bytes([0]))
        obs = list(self.sock.recv(4096))
        obs = np.asarray(obs, dtype=np.int8)
        obs = reshape_obs(obs)

        return obs, reward, done, {}

    def _close(self):
        if self.sock:
            self.sock.close()

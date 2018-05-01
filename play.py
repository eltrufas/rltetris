from baselines import deepq
import numpy as np
from socketserver import TCPServer, BaseRequestHandler

act = deepq.load("tetris.pkl")


class ActionHandler(BaseRequestHandler):
    def handle(self):
        while True:
            data = [x for x in bytes(self.request.recv(1024))]
            print(data)
            obs = np.asarray([data], dtype=np.int8)
            action = bytes([act(obs)[0]])
            self.request.sendall(action)


while True:
    with TCPServer(('localhost', 5050), ActionHandler) as server:
        server.serve_forever()

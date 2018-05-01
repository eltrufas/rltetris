import socket

sock = socket.socket()
sock.connect(('localhost', 3030))

# print(sock.recv(4096))
# print(sock.recv(4096))
print(list(sock.recv(4096)))

sock.send(bytes([1]))

print(len(sock.recv(4096)))
sock.close()

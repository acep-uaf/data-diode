""" 
This is an example Python script provided to aid in the integration of a Fend Data Diode in a TCP Client/Server application.

For this example a simple socket is opened and listens for the diode's TCP client to connect. When data is received it is printed to
the console window.

For these scripts to function, the diode will need to be setup in TCP client/server mode. Please see the user manual
on instruction for this setup process.

These scripts will continuously run until aborted. 

The following is the expected step by step process of what the script does:
1. The script creates a socket object named "server".
2. The socket object begins listening for a connection from the diode's TCP client
3. If a connection is made, the script enters a receiving loop.
4. Messages sent from the diode will be printed from the screen.
5. The script repeats the receiving loop.
"""

import socket
 
# Change this to the Target TCP Server IP Address in your diode's Output Side TCP Client Settings.
targetTcpServerIP = "192.168.1.20"

# Change this to the Target TCP Server Port in your diode's Output Side TCP Client Settings.
targetTcpServerPort = 503 

# Set up a TCP server
server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.bind((targetTcpServerIP, targetTcpServerPort))
 
# Begin listening
server.listen(1)
 
while True:
    print("Waiting for connection")
    connection, client = server.accept()
 
    try:
        # Print to console the connected client IP address
        print("Connected to client IP: {}".format(client))
         
        # Receive and print data 10 Kbytes at a time, as long as the client is sending something
        while True:
            data = connection.recv(10240)
            print(f"Received data: {data.decode()}")
 
            if not data:
                break
 
    finally:
        connection.close()
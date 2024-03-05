""" 
This is an example Python script provided to aid in the integration of a Fend Data Diode in a UDP Client/Server application.

For this example a simple socket is opened and listens for the diode's UDP client to connect. When data is received it is printed to
the console window.

For these scripts to function, the diode will need to be setup in UDP client/server mode. Please see the user manual
on instruction for this setup process.

These scripts will continuously run until aborted. 

The following is the expected step by step process of what the script does:
1. The script creates a socket object named "server".
2. The socket object begins listening for datagrams from the diode's UDP client
3. Messages sent from the diode's UDP client will be printed from the screen.
4. The script repeats the receiving loop.
"""

import socket
 
# Change this to the Target UDP Server IP Address in your diode's Output Side UDP Client Settings.
targetUDPServerIP = "192.168.1.20"

# Change this to the Target UDP Server Port in your diode's Output Side UDP Client Settings.
targetUDPServerPort = 503 

# Set up a UDP server
server = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
server.bind((targetUDPServerIP, targetUDPServerPort))
 
# Receive and print data 10 Kbytes at a time, as long as the client is sending something
try:
    while True:
        data, addr = server.recvfrom(10240)
        print(f"Received data from {addr}: {data}")

        if not data:
            break
finally:
    server.close()
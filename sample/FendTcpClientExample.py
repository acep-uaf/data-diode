# Fend TCP and UDP Client/Server Example Python Test Scripts
# Copyright (c) 2023 Fend Incorporated. All rights reserved.

"""
This is an example Python script provided to aid in the integration of a Fend Data Diode in a TCP Client/Server application.

For this example a simple socket is extabilshed with the diode's TCP server and a test message is continuously sent to the input side
of the diode.

For the TCP Client/Server feature to function properly, data must be sent in data chunks no larger than 1460 Bytes. After
each chunk of data is sent you must wait for the diode to acknowledge (ACK) the send. This ACK will come in the form of
a text string 'OK\r\n'. When this ACK is received new data is free to be sent.

For these scripts to function, the diode will need to be setup in TCP Client/Server mode. Please see the user manual
on instruction for this setup process.

These scripts will continuously run until aborted.

The following is the expected step by step process of what the script does:
1. The script creates a socket object named "client".
2. A connection attempt to the diode is made.
3. If the connection is made, the script enters a sending loop.
4. This loop will check to see if the data is larger than 1460 bytes. If so, it will get a chunk and send that.
5. The script will then wait for the diode to send its ACK before proceeding.
6. The loop will continue until all the data has been chunked and transmitted.
7. The script waits 1 second.
8. The script repeats the send loop.

"""

import socket
from time import sleep
from timeit import default_timer

# Change this to the IP address of your diode's Input Side IP Address.
diodeInputSideIP = "192.168.1.99"

# Change this to the Diode TCP Server Port in your diode's Input Side TCP Passthrough Settings.
diodeTcpPassthroughPort = 50000

# create a socket
client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.settimeout(5)

# connect to the server
client.connect((diodeInputSideIP, diodeTcpPassthroughPort))

# Data to be chunked and sent to the diode's TCP server
data = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

try:
    while True:
        # Send data to diode
        if len(data) > 1460: # If the data you wish to send is larger than 1460 bytes, you need to chunk.
            index = 0
            while index < len(data):
                # Create chunk of 1460 chars
                chunk = data[index : index + 1460]

                # Send chunk to the diode's TCP server
                client.send(chunk.encode())

                #Set timer to wait for ACK from diode
                start = default_timer()

                # Wait for ACK from diode
                while True:
                    response = client.recv(4096)
                    # You must wait for the diode to respond with its ACK "OK\r\n"
                    if response.decode() == "OK\r\n":
                        break
                    elif default_timer() - start > 5:
                        raise TimeoutError

                # Print debug statement
                print(f'Successfully sent message to diode: {chunk}')

                # Adjust index
                index += 1460

        # This is for testing purposes only. Having a 1 second delay makes testing the cross diode connection easier to see visually.
        # For maximum performance, remove this delay.
        sleep(1)
except TimeoutError:
    print("No response was received from the diode. Please check your settings and try again.")
finally:
    client.close()

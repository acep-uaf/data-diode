""" 
This is an example Python script provided to aid in the integration of a Fend Data Diode in a TCP Passthrough application.

For this example a simple socket is extabilshed with the diode and a test message is continuously sent to the input side
of the diode.

For the TCP Passthrough protocol to function properly, data must be sent in data chunks no larger than 1460 Bytes. After
each chunk of data is sent you must wait for the diode to acknowledge (ACK) the send. This ACK will come in the form of
a text string 'OK\r\n'. When this ACK is received new data is free to be sent.

For these scripts to function, the diode will need to be setup in TCP passthrough mode. Please see the user manual
on instruction for this setup process.

These scripts will continuously run until aborted. 

The following is the expected step by step process of what the script does:
1. The script creates a socket object named "client".
2. A connection attempt to the diode is made.
3. If the connection is made, the script enters a sending loop.
4. A test message is sent to the diode. This test message includes a number that increments each send.
5. The script will then wait for the diode to send its ACK before proceeding.
6. The number of sends variable is incremented.
7. The script waits 1 second.
8. The script repeats the send loop.

"""

import socket
from time import sleep

# Change this to the IP address of your diode's Input Side IP Address.
diodeInputSideIP = "192.168.1.99" 

# Change this to the Diode TCP Server Port in your diode's Input Side TCP Passthrough Settings.
diodeTcpPassthroughPort = 50000 

# create a socket
client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.settimeout(1)

# connect to the server
client.connect((diodeInputSideIP, diodeTcpPassthroughPort))

try:
    numberOfSends = 1

    while True:
        # Send data to diode
        sendMessage = f"This is tcp passthrough test message number: {numberOfSends}"
        client.send(sendMessage.encode())

        # Wait for ACK from diode
        while True:
            response = client.recv(4096)
            # You must wait for the diode to respond with its ACK "OK\r\n"
            if response.decode() == "OK\r\n":
                break

        # Print debug statement
        print(f'Successfully sent message to diode: {sendMessage}')
        
        # Update number of sends index
        numberOfSends += 1

        # This is for testing purposes only. Having a 1 second delay makes testing the cross diode connection easier to see visually.
        # For maximum performance, remove this delay.
        sleep(1)
finally:
    client.close()
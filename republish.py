# Data Diode TCP Stream Republisher

import socket
import paho.mqtt.client as mqtt

# Configuration Settings

targetTcpServerIP = "192.168.1.20"
targetTcpServerPort = 503
mqttBrokerIP = "test.mosquitto.org"
mqttBrokerPort = 1883
mqttTopic = "test/message"

# TCP Socket

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.bind((targetTcpServerIP, targetTcpServerPort))

# Begin Listening

server.listen(1)

# MQTT Client

mqtt_client = mqtt.Client()
mqtt_client.connect(mqttBrokerIP, mqttBrokerPort)

while True:
    print(">> Waiting for connection...")
    connection, client = server.accept()

    try:
        print(">> Connected to client IP: {}".format(client))

        while True:
            data = connection.recv(10240)
            print(f">> Received data: {data.decode()}")

            if not data:
                break

            mqtt_client.publish(mqttTopic, data.decode())

    finally:
        connection.close()

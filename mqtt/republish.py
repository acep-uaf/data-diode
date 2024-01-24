# Data Diode TCP Stream Republisher

import socket
import paho.mqtt.client as mqtt
import json
from datetime import datetime

# Configuration Settings

targetTcpServerIP = "192.168.1.20"
targetTcpServerPort = 503

mqttBrokerIP = "test.mosquitto.org"
mqttBrokerPort = 1883
mqttTopic = "test/message"

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.bind((targetTcpServerIP, targetTcpServerPort))
server.listen(1)

mqtt_client = mqtt.Client()
mqtt_client.connect(mqttBrokerIP, mqttBrokerPort)

while True:
    print(">> Waiting for connection...")
    connection, client = server.accept()

    try:
        print(">> Connected to client IP: {}".format(client))

        while True:
            data = connection.recv(10240)
            if not data:
                break

            # print(f">> Received data: {data.decode()}")

            message = {
                "timestamp": datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S"),
                "topic": mqttTopic,
                "data": data.decode()
            }

            json_message = json.dumps(message)

            mqtt_client.publish(mqttTopic, json_message)

    finally:
        connection.close()

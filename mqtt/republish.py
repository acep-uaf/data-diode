# Data Diode TCP Stream Republisher

import socket
import paho.mqtt.client as mqtt
import json
import yaml
from datetime import datetime

def stream():
    """Republish TCP stream to MQTT broker."""
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

if __name__ == "__main__":
    print("🚀 Data Diode TCP Stream Republisher")

    with open("config.yaml", "r") as file:
        config = yaml.safe_load(file)

    # Configuration Settings

    targetTCPServerIP = config["diode"]["topology"][0]["output"]["ip"]
    targetTCPServerPort = config["diode"]["topology"][0]["output"]["port"]

    mqttBrokerIP = config["diode"]["topology"][2]["broker"]["ip"]
    mqttBrokerPort = config["diode"]["topology"][2]["broker"]["port"]
    mqttTopic = config["diode"]["topology"][2]["broker"]["topic"]

    try:
        server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        server.bind((targetTCPServerIP, targetTCPServerPort))
        server.listen(1)

        mqtt_client = mqtt.Client()
        mqtt_client.connect(mqttBrokerIP, mqttBrokerPort)

        while True:
            stream()
    except OSError as e:
        print("⚠️ {}".format(e))

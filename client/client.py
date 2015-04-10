import paho.mqtt.client as mqtt

# topic to subscribe to 
SUBTOPIC = "/netCloudDash/control/+"

def on_connect(client, userdata, flags, rc):
    print("Connected to server")

    client.subscribe(SUBTOPIC)

def on_message(client, userdata, msg):
    print(msg.payload)

def main():
    client = mqtt.Client()

    client.on_connect = on_connect
    client.on_message = on_message

    client.connect("iot.eclipse.org", 1883, 60)

    client.loop_forever()


if __name__ == '__main__':
    main()

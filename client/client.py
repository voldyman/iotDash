import paho.mqtt.client as mqtt
import time

# topic to subscribe to 
SUBTOPIC = "/netCloudDash/control/+"

# led state topic
STATETOPIC = "/netCloudDash/control/led-state"

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

    client.loop_start()

    while True:
        client.publish(STATETOPIC, "on")
        print "on sent"
        time.sleep(3)

        client.publish(STATETOPIC, "off")
        print "off sent"
        time.sleep(2)


if __name__ == '__main__':
    main()

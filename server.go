package main

import (
	"fmt"
	"net/http"
	"text/template"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/gorilla/mux"
)

const (
	SERVER   = "tcp://iot.eclipse.org:1883" // thanks eclicpse!
	SUBTOPIC = "/netCloudDash/control/+"
	PUBTOPIC = "/netCloudDash/control/%s"
)

func main() {
	name := "Overlord"
	// connect to mqtt
	client := connectMQTT(name)

	// create server
	serveWeb(name, client)
	// enable comm
}

func connectMQTT(name string) *mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(SERVER).SetClientID(name).SetCleanSession(true)

	opts.OnConnect = func(c *mqtt.Client) {
		if token := c.Subscribe(SUBTOPIC, 2, messageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func messageReceived(client *mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message received at: %s\nMessage:%s", msg.Topic(), msg.Payload())
}

// HTTP Server Code Below
func serveWeb(name string, client *mqtt.Client) {
	router := mux.NewRouter()
	router.HandleFunc("/", webHome)
	router.HandleFunc("/action/{state}", webAction(name, client))
	http.ListenAndServe(":8081", router)
}

func webHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("dash").Parse(basicPage)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	var data struct {
		Title   string
		Clients int
	}
	data.Title = "Dashboard"
	data.Clients = 1
	t.ExecuteTemplate(w, "dash", data)
}

func webAction(name string, client *mqtt.Client) func(w http.ResponseWriter, r *http.Request) {
	pubTopic := fmt.Sprintf(PUBTOPIC, name)

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		val := vars["state"]
		send := "off"

		fmt.Println("Got Action Request")
		fmt.Println("state: ", val)

		if val == "on" {
			send = "on"
		}
		go func() {
			if token := client.Publish(pubTopic, 1, false, send); token.Wait() && token.Error() != nil {
				fmt.Println("Error occured during publish")
			}
			fmt.Println("Sent Action")
		}()
		w.Write([]byte("ok"))
	}
}

var basicPage = `
<html>
<head><title>{{ .Title }}</title></head>
<body>
<h1>NetPlug Cloud Dashboard</h1>
<p>Number of clients: {{ .Clients }}</p>
<p><button id='on'>On</button></p>
<p><button id='off'>Off</button></p>
<script>
// get DOM element
function $(elName) { return document.getElementById(elName); }

// get request
function get(url) {
    var r = new XMLHttpRequest();
    r.open("GET", url, true);
    r.onreadystatechange = function () {
      if (r.readyState != 4 || r.status != 200) return;
    };
    r.send();
}

$('on').onclick = function() {
    get('/action/on');
};

$('off').onclick = function() {
    get('/action/off');
};

</script>
</body>

`

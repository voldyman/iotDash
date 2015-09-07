package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/gorilla/mux"
	eventsource "gopkg.in/antage/eventsource.v1"
)

const (
	SERVER     = "tcp://iot.eclipse.org:1883" // thanks eclicpse!
	SUBTOPIC   = "/netCloudDash/control/+"
	STATETOPIC = "/netCloudDash/control/led-state"
	PUBTOPIC   = "/netCloudDash/control/%s"
)

var (
	tpls *template.Template
)

func main() {
	name := "Overlord"

	// sketup templates
	var err error
	tpls, err = template.ParseGlob("./tmpl/*.tpl")
	if err != nil {
		panic("Unable to parse templates")
	}

	// connect to mqtt
	client := connectMQTT(name)

	// create server
	serveWeb(name, client)
	// enable comm
}

func connectMQTT(name string) *mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(SERVER).SetClientID(name).SetCleanSession(true)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

// HTTP Server Code Below
func serveWeb(name string, client *mqtt.Client) {
	es := eventsource.New(nil, nil)
	defer es.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", webHome)
	//	router.HandleFunc("/events", webEvents(es))
	router.HandleFunc("/action/{state}", webAction(name, client))
	router.Handle("/events", es)
	router.HandleFunc("/public/{path}", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Listenting on port ", port)

	startEvents(es, client)
	http.ListenAndServe(":"+port, router)
}

func webEvents(es eventsource.EventSource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request")
		es.ServeHTTP(w, r)
	}
}
func startEvents(es eventsource.EventSource, client *mqtt.Client) {
	client.Subscribe(STATETOPIC, 0, func(c *mqtt.Client, m mqtt.Message) {
		fmt.Println("Got message", string(m.Payload()))
		es.SendEventMessage(string(m.Payload()), "led-state", "")
	})
}

func webHome(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title   string
		Clients int
	}
	data.Title = "Dashboard"
	data.Clients = 1
	tpls.ExecuteTemplate(w, "dash.tpl", data)
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

package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	// connect to mqtt
	// create server
	serveWeb()
	// enable comm
}

// HTTP Server Code Below
func serveWeb() {
	router := mux.NewRouter()
	router.HandleFunc("/", webHome)
	router.HandleFunc("/action/{state}", webAction)
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

func webAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Got Action Request")
	fmt.Println("state: ", vars["state"])

	w.Write([]byte("ok"))
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
      alert("Success: " + r.responseText);
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

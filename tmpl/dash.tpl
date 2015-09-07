<html>
    <head>
        <title>{{ .Title }}</title>
        <style>
         body {
             width: 80%;
             margin: 0 auto;
         }
         .container {
             width: 300px;
             margin: 0 auto;
         }
         h1 {
             text-align:center;
         }
        </style>
    </head>
    <body>
        <div class="container" >

        <h1>NetPlug Cloud Dashboard</h1>
        <p>Press the buttons to send command to your iot Device</p>
        <img id='led-status' src="public/grey.png">
        <p>
            <button style="float:left"  id='on'>On</button>
            <button style="float:right" id='off'>Off</button>
        </p>
        </div>
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

         function ledNeutral() {
             $('led-status').src = '/public/grey.png';
         }

         function ledOn() {
             $('led-status').src = '/public/blue.png';
         }

         function ledOff() {
             $('led-status').src = '/public/black.png';
         }

         function getEvents() {
             var evts = new EventSource("/events");

             evts.addEventListener('open', function() {
                 console.log("Connected to server");
             }, false);

             evts.addEventListener('led-state', function(ev) {
                 console.log(ev.data);
                 if (ev.data === 'on') {
                     ledOn();
                 } else if (ev.data === 'off') {
                     ledOff();
                 } else {
                     ledNeutral();
                 }

             }, false);

             evts.addEventListener('error', function(err) {
                 console.log(err);
                 console.log("Error occured");
             }, false);
         }

         $('on').onclick = function() {
             get('/action/on');
         };

         $('off').onclick = function() {
             get('/action/off');
         };

         ledNeutral();

         getEvents();

        </script>
    </body>
</html>

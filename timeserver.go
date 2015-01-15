// Author: Tyler Hills

package main

import (
	"fmt"
	"time"
	"flag"		
	"os"
	"net/http"
	)

// Version Number
const AppVersion = "timeserver version: 1.0"

// Handler for timeserver, prints the current time to the second
func timeserver(w http.ResponseWriter, r *http.Request) {
	
	// if url is off then 404
	if r.URL.Path != "/time/" {
	NotFoundHandler(w, r)
	return
	}
	// time formatting
	const layout = "3:04:05 PM"
	
	// get current time
	t := time.Now().Local()
	
	// html formatting and displaying current time
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "<head>")
	fmt.Fprintln(w, "<style>")
	fmt.Fprintln(w, "p {font-size: xx-large}")
	fmt.Fprintln(w, "span.time {color: red}")
	fmt.Fprintln(w, "</style>")
	fmt.Fprintln(w, "</head>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, "<p>The time is now <span class=\"time\">" + t.Format(layout) + "</span>.</p>")
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}

// Login Handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index/" {
		NotFoundHandler(w, r)
		return
		}
		
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, "<form action=\"login\">")
	fmt.Fprintln(w, "What is your name, Earthling?")
	fmt.Fprintln(w, "<input type=\"text\" name=\"name\" size=\"50\">")
	fmt.Fprintln(w, "<input type=\"submit\">")
	fmt.Fprintln(w, "</form>")
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}

// 404 error handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, "<p>These are not the URLs you're looking for.</p>")
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}

func main() {
	// version flag -V
	version := flag.Bool("V", false, "prints current version")
	
	// port flag -p port_number, default to 8080
	port := flag.String("port", "8080", "sets service port number")
	
	// parse flags, if -V print version number and exit
	flag.Parse()
	if *version {
      		fmt.Println(AppVersion)
      		os.Exit(0)
    	}
	
	// add handlers to the DefaultServeMux
	http.HandleFunc("/", LoginHandler)
	http.HandleFunc("/index/", LoginHandler)
	http.HandleFunc("/time/", timeserver)
	
	// Start the server, print error message if any problem
	err := http.ListenAndServe("localhost:" + *port, nil)
	if err != nil {
		fmt.Println("Server Error: %s", err)
		os.Exit(1)
	}
}

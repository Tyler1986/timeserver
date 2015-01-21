// Author: Tyler Hills

package main

import (
	"fmt"
	"time"
	"flag"		
	"os"
	"net/http"
	"strings"
	"sync"
	)

// Version Number
const AppVersion = "timeserver version: 2.0"
		
//var names map[string]string

var names = struct{
    sync.RWMutex
    m map[string]string
}{m: make(map[string]string)}

var UID string = "A"

// Handler for timeserver, prints the current time to the second
func timeserver(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	// if url is off then 404
	if r.URL.Path != "/time/" {
	NotFoundHandler(w, r)
	return
	}
	
	// time formatting
	const local = "3:04:05 PM"
	const UTC = "15:04:05 UTC"
	
	timeString := ""
	
	// get current time
	t := time.Now().Local()
	
	// check for valid cookie which determines the message to show
	sessionCookie, err := r.Cookie("login")
	if err != nil {
		timeString = "<p>The time is now <span class=\"time\">" + t.Format(local) + "</span>.</p>" 	  
	} else {
		names.RLock()
		timeString = "<p>The time is now <span class=\"time\">" + t.Format(local) + "</span>," + names.m[sessionCookie.Value] + "</p>" 
		names.RUnlock()
	}
	
	// html formatting and displaying current time
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "<head>")
	fmt.Fprintln(w, "<style>")
	fmt.Fprintln(w, "p {font-size: xx-large}")
	fmt.Fprintln(w, "span.time {color: red}")
	fmt.Fprintln(w, "</style>")
	fmt.Fprintln(w, "</head>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, timeString)
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}

// Login Handler
func LoginForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	if r.URL.Path != "/" && r.URL.Path != "/index/" {
		NotFoundHandler(w, r)
		return
	}
	
	// display login screen or welcome message depending on if cookie is present
	sessionCookie, err := r.Cookie("login")
	if err != nil {
		fmt.Fprintln(w, "<html>")
		fmt.Fprintln(w, "<body>")
		fmt.Fprintln(w, "<form action=\"login\">")
		fmt.Fprintln(w, "What is your name, Earthling?")
		fmt.Fprintln(w, "<input type=\"text\" name=\"name\" size=\"50\">")
		fmt.Fprintln(w, "<input type=\"submit\">")
		fmt.Fprintln(w, "</form>")
		fmt.Fprintln(w, "</body>")
		fmt.Fprintln(w, "</html>")
	} else {
		if  len(strings.TrimSpace(sessionCookie.Value)) != 0 {
			names.RLock()
			fmt.Fprintln(w, "Greetings, " + names.m[sessionCookie.Value])
			names.RUnlock()
		} else {
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
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	name := r.URL.Query().Get("name")
	
	// create and set the cookie, also add the name and cookie UID to the map
	if len(strings.TrimSpace(name)) != 0 {
		cookie := http.Cookie{Name: "login", Value: UID}
		names.Lock()
		names.m[UID] = name
		names.Unlock()
		UID = UID + UID
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		fmt.Fprint(w, "C'mon, I need a name.")
	}
}

// Logout Handler
func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	if r.URL.Path != "/logout/" {
		NotFoundHandler(w, r)
		return
	}
	
	//http.SetCookie(w, &http.Cookie{Name: "login", Value: "", MaxAge: -1, Path: "/"})
	deleteCookie, _ := r.Cookie("login")
	deleteCookie.Value = "" 
	deleteCookie.MaxAge = -1
	deleteCookie.Path = "/"
	http.SetCookie(w, deleteCookie)
	
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "<head>")
	fmt.Fprintln(w, "<META http-equiv=\"refresh\" content=\"10;URL=/\">")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, "<p>Good-bye.</p>")
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}

// 404 error handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
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
	
	//names = make(map[string]string)
	
	// add handlers to the DefaultServeMux
	http.HandleFunc("/", LoginForm)
	http.HandleFunc("/index/", LoginForm)
	http.HandleFunc("/time/", timeserver)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout/", Logout)
	
	
	// Start the server, print error message if any problem
	err := http.ListenAndServe("localhost:" + *port, nil)
	if err != nil {
		fmt.Println("Server Error: %s", err)
		os.Exit(1)
	}
}

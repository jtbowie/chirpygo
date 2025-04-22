package main

import "net/http"

func main() {
	myServeMux := http.NewServeMux()
	var myHTTPServer http.Server

	myServeMux.Handle("/", http.FileServer(http.Dir(".")))
	myHTTPServer.Handler = myServeMux
	myHTTPServer.Addr = ":8080"
	myHTTPServer.ListenAndServe()
}

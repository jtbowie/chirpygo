package main

import "net/http"

func httpHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	httpWriter.Header().Set("Content-type", "text/plain; charset=utf-8")
	httpWriter.WriteHeader(200)
	_, err := httpWriter.Write([]byte("OK"))
	if err != nil {
		// TODO: Implement Error message/handling
		// This should never fail, but it's worth handling the error case"
		return
	}

}

func main() {
	myServeMux := http.NewServeMux()
	var myHTTPServer http.Server

	defer myHTTPServer.Close()

	httpFileServer := http.FileServer(http.Dir("."))
	myServeMux.HandleFunc("/healthz", httpHandler)
	myServeMux.Handle("/app/", http.StripPrefix("/app", httpFileServer))
	myHTTPServer.Handler = myServeMux
	myHTTPServer.Addr = ":8080"
	err := myHTTPServer.ListenAndServe()
	if err != nil {
		return
	}
}

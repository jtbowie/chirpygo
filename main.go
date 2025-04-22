package main

import "net/http"
import "sync/atomic"
import "fmt"

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (apicfg *apiConfig) resetFileServerHitCounter(httpWriter http.ResponseWriter, httpReq *http.Request) {
	apicfg.fileserverHits.Store(0)
}

func setTextHeader(httpWriter http.ResponseWriter) {
	httpWriter.Header().Set("Content-type", "text/plain; charset=utf-8")
}

func (apicfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apicfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (apicfg *apiConfig) serveMetrics(httpWriter http.ResponseWriter, httpReq *http.Request) {
	setTextHeader(httpWriter)
	httpWriter.WriteHeader(200)
	response := fmt.Sprintf("Hits: %d\n", apicfg.fileserverHits.Load())
	httpWriter.Write([]byte(response))

}

func httpHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	setTextHeader(httpWriter)
	httpWriter.WriteHeader(200)
	_, err := httpWriter.Write([]byte("OK\n"))
	if err != nil {
		// TODO: Implement Error message/handling
		// This should never fail, but it's worth handling the error case"
		return
	}

}

func main() {
	myServeMux := http.NewServeMux()
	var myHTTPServer http.Server
	var apicfg apiConfig

	defer myHTTPServer.Close()

	httpFileServer := http.FileServer(http.Dir("."))
	myServeMux.HandleFunc("GET /healthz", httpHandler)
	myServeMux.HandleFunc("/reset", apicfg.resetFileServerHitCounter)
	baseAppHandler := http.StripPrefix("/app", httpFileServer)
	myServeMux.Handle("GET /app/", apicfg.middlewareMetricsInc(baseAppHandler))
	myServeMux.HandleFunc("/metrics", apicfg.serveMetrics)
	myHTTPServer.Handler = myServeMux
	myHTTPServer.Addr = ":8080"
	err := myHTTPServer.ListenAndServe()
	if err != nil {
		return
	}
}

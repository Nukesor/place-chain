package app

import (
	"fmt"
	"net/http"
)

type WebServer struct {
	app KVStoreApplication
}

func (*WebServer) setPixel(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", r.URL.Query())
}

func (*WebServer) getPixels(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", r.URL.Query())
}

func (self *WebServer) LaunchHTTP() {
	http.HandleFunc("/pixel/", self.setPixel)
	http.HandleFunc("/pixels/", self.getPixels)
	http.ListenAndServe(":8080", nil)
}

package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebServer struct {
	App *KVStoreApplication
}

type CreationResponse struct {
	message string `json:"message"`
	color   string `json:"status"`
}

func (*WebServer) setPixel(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	var created CreationResponse
	if data["x"] != nil && data["y"] != nil {
		created = CreationResponse{
			"Placed successfully!",
			"Rainbow",
		}
	} else {
		created = CreationResponse{
			"Invalid request!",
			"",
		}
	}
	b, _ := json.Marshal(created)
	fmt.Fprintf(w, "%s", b)
}

func (*WebServer) getPixels(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", r.URL.Query())
}

func (self *WebServer) LaunchHTTP() {
	http.HandleFunc("/pixel/", self.setPixel)
	http.HandleFunc("/pixels/", self.getPixels)
	port := "8080"
	fmt.Printf("Listening on http://localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

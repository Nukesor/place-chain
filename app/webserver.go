package app

import (
	"../types"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebServer struct {
	App *KVStoreApplication
}

func (self *WebServer) setPixel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var pr types.PixelRequest
	err := decoder.Decode(&pr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Bad request: %s", err)
		return
	}
	res, err := self.App.SetPixel(pr.X, pr.Y)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}
	defer r.Body.Close()
	fmt.Fprintf(w, "Response: %s", res)
}

func (self *WebServer) getPixels(w http.ResponseWriter, r *http.Request) {
	data := self.App.GetGrid()
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error returning pixels: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (self *WebServer) LaunchHTTP() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/script.js")
	})
	http.HandleFunc("/style.less", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/style.less")
	})
	http.HandleFunc("/pixel/", self.setPixel)
	http.HandleFunc("/pixels/", self.getPixels)
	port := "8080"
	fmt.Printf("Listening on http://localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

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

type CreationResponse struct {
	message string `json:"message"`
	color   string `json:"status"`
}

func (*WebServer) setPixel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var pr types.PixelRequest
	err := decoder.Decode(&pr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Bad request: %s", err)
		return
	}
	defer r.Body.Close()
	fmt.Fprintf(w, "You sent this: %v", pr)
}

func (*WebServer) getPixels(w http.ResponseWriter, r *http.Request) {
	nums := []int{1, 2, 3}
	b, err := json.Marshal(nums)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error returning pixels: ", err)
	}
	w.Write(b)
}

func (self *WebServer) LaunchHTTP() {
	http.HandleFunc("/pixel/", self.setPixel)
	http.HandleFunc("/pixels/", self.getPixels)
	port := "8080"
	fmt.Printf("Listening on http://localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

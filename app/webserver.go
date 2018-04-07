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

func (*WebServer) getPixels(w http.ResponseWriter, r *http.Request) {
	nums := []int{1, 2, 3}
	b, err := json.Marshal(nums)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error returning pixels: ", err)
		return
	}
	w.Write(b)
}

func (self *WebServer) LaunchHTTP() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	http.HandleFunc("/pixel/", self.setPixel)
	http.HandleFunc("/pixels/", self.getPixels)
	port := "8080"
	fmt.Printf("Listening on http://localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

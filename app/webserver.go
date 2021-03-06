package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"place-chain/types"
)

type WebServer struct {
	PlacechainApp *PlacechainApp
	TwitterCache  *TwitterCache
}

func (self *WebServer) setPixel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var pr types.PixelRequest
	err := decoder.Decode(&pr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Bad request: %s\n", err)
		return
	}
	defer r.Body.Close()
	if !self.PlacechainApp.IsTransactionValid(pr.ToTransaction()) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Printf("Unprocessable pixel request")
		return
	}
	_, err = self.PlacechainApp.PublishTx(pr.ToTransaction())
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ws *WebServer) twitterUser(w http.ResponseWriter, r *http.Request) {
	name := path.Base(r.URL.Path)
	user := ws.TwitterCache.getUser(name)
	responseData, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error returning user: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func (self *WebServer) getPixels(w http.ResponseWriter, r *http.Request) {
	data := self.PlacechainApp.GetGrid()
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error returning pixels: %v ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (self *WebServer) register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var rr types.RegisterRequest

	if err := decoder.Decode(&rr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Bad request: %s\n", err)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	defer r.Body.Close()

	if !rr.IsValid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Printf("Unprocessable pixel request")
		return
	}

	if err := self.PlacechainApp.RegisterUser(rr); err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (webServer *WebServer) isRegistered(w http.ResponseWriter, r *http.Request) {
	twitterHandle := r.URL.Query().Get("twitterHandle")
	if twitterHandle == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := webServer.PlacechainApp.GetPubKey(twitterHandle)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	return
}

func (self *WebServer) LaunchHTTP() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/script.js")
	})
	http.HandleFunc("/blank_profile_100.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/blank_profile_100.png")
	})
	http.HandleFunc("/style.less", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/style.less")
	})
	http.HandleFunc("/bundle.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/bundle.js")
	})

	http.HandleFunc("/pixel", self.setPixel)
	http.HandleFunc("/pixel/", self.setPixel)

	http.HandleFunc("/pixels", self.getPixels)
	http.HandleFunc("/pixels/", self.getPixels)

	http.HandleFunc("/register", self.register)
	http.HandleFunc("/register/", self.register)

	http.HandleFunc("/profile", self.twitterUser)
	http.HandleFunc("/profile/", self.twitterUser)

	http.HandleFunc("/isRegistered", self.isRegistered)
	http.HandleFunc("/isRegistered/", self.isRegistered)
	port := "8080"
	fmt.Printf("Listening on http://localhost:%s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		fmt.Printf("Could not serve via http: %s", err)
	}
}

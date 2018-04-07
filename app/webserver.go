package app

import (
    "fmt"
    "net/http"
)

func setPixel(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s", r.URL.Query())
}

func getPixels(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s", r.URL.Query())
}

func LaunchHTTP() {
    http.HandleFunc("/pixel/", setPixel)
    http.HandleFunc("/pixels/", getPixels)
    http.ListenAndServe(":8080", nil)
}

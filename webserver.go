package main

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

func main() {
    http.HandleFunc("/pixel", setPixel)
    http.HandleFunc("/pixels", getPixels)
    http.ListenAndServe(":8080", nil)
}

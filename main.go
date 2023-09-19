package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mymain/geo"
	"net/http"
	"os"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func main(){
	geo.API_KEY = os.Getenv("ROUTES_KEY")
	http.HandleFunc("/route", routeHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	port := os.Getenv("port")
	if port == "" {
		port = "8080"
	}
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}


func routeHandler(w http.ResponseWriter, r *http.Request) {
	params, err := geo.ParseRouteQueryParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	if err := geo.CheckParams(params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	
	fmt.Println("params: ", params)

	resp, err := geo.GetRouteData(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

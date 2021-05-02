package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CallResponseJSON struct {
	Message string
	Status string
}


func Other(w http.ResponseWriter, r *http.Request) {
	if res, err := w.Write([]byte("[PONG] Hello World\n")); err != nil {
		panic(err)
	} else {
		fmt.Print(res)
	}
}


func OAuthHttpCallback(res http.ResponseWriter, req *http.Request) {
	// print the body that comes with the request.
	fmt.Println(req.Body)

	responseBody := CallResponseJSON{"Authentication success", "ok"}

	data, err := json.Marshal(responseBody)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(200)
	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}


func FileServer(res http.ResponseWriter, req *http.Request) http.Handler {
	return http.FileServer(http.Dir("./app"))
}
package main

import (
	"log"
	"net/http"
	"paredao/src/poll"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()
	route.HandleFunc("/vote", poll.CreateVote).Methods("POST")
	route.HandleFunc("/vote", poll.GetVotes).Methods("GET")

	http.Handle("/", route)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

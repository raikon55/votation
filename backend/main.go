package main

import (
	"log"
	"net/http"
	"os"
	poll "paredao/src/poll"
	queue "paredao/src/queue"

	"github.com/gorilla/mux"
)

func main() {
	queueName := os.Getenv("RABBITMQ_QUEUE")
	rmq := queue.InitRabbitMQ(queueName)
	p := poll.Init(rmq)

	route := mux.NewRouter()
	route.HandleFunc("/vote", p.CreateVote).Methods("POST")
	route.HandleFunc("/vote", p.GetVotes).Methods("GET")

	http.Handle("/", route)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

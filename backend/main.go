package main

import (
	"log"
	"net/http"
	"os"
	poll "paredao/src/poll"
	queue "paredao/src/queue"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	queueName := os.Getenv("RABBITMQ_QUEUE")
	rmq := queue.InitRabbitMQ(queueName)
	p := poll.Init(rmq)

	route := mux.NewRouter()
	route.HandleFunc("/vote", p.CreateVote).Methods("POST")
	route.HandleFunc("/vote/candidate", p.GetVotesByCandidate).Methods("GET")
	route.HandleFunc("/vote/total", p.GetVotes).Methods("GET")

	http.Handle("/", route)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

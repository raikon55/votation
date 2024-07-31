package poll

import (
	"fmt"
	"net/http"
	queue "paredao/src/queue"
)

type PollServer interface {
	CreateVote(response http.ResponseWriter, request *http.Request)
	GetVotes(response http.ResponseWriter, request *http.Request)
}

type Poll struct {
	rm queue.RabbitMQ
	ps PollServer
}

func (p *Poll) CreateVote(response http.ResponseWriter, request *http.Request) {
	p.rm.Enqueue("", "") // refactor
	response.WriteHeader(http.StatusCreated)
	fmt.Fprint(response, "Voto registrado")
}

func (p *Poll) GetVotes(response http.ResponseWriter, request *http.Request) {

}

func Init(rm queue.RabbitMQ) (p Poll) {
	p.rm = rm
	return
}

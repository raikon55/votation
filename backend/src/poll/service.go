package poll

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	queue "paredao/src/queue"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PollServer interface {
	CreateVote(response http.ResponseWriter, request *http.Request)
	GetVotes(response http.ResponseWriter, request *http.Request)
	save(c candidate)
	initMongo()
	consumer()
}

type Poll struct {
	rm         queue.RabbitMQ
	conn       *mongo.Client
	collection *mongo.Collection
}

type candidate struct {
	Name string `json:"name"`
	Vote int    `json:"vote"`
}

type pollResult struct {
	result []candidate
}

func (p *Poll) CreateVote(response http.ResponseWriter, request *http.Request) {
	var c candidate
	message := "{\"message\": \"Voto registrado\", \"candidate\": %q}"
	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&c)
	if err != nil {
		log.Println("Invalid body")
	}

	byt, err := json.Marshal(c)
	if err != nil {
		log.Println("Invalid body")
	}

	p.rm.Enqueue(string(byt))
	response.Header().Add("content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	fmt.Fprintf(response, message, c.Name)
}

func (p *Poll) GetVotes(response http.ResponseWriter, request *http.Request) {

}

func Init(rm queue.RabbitMQ) (p Poll) {
	p.rm = rm
	p.initMongo()

	enabled, err := strconv.ParseBool(os.Getenv("ENABLE_CONSUMER"))
	if err != nil {
		enabled = true
	}

	if enabled {
		p.consumer()
	}

	return
}

func (p *Poll) consumer() {
	go func() {
		votes := p.rm.Consume()
		forever := make(chan bool)

		for v := range votes {
			var c candidate
			err := json.Unmarshal(v.Body, &c)
			if err != nil {
				p.rm.Enqueue(string(v.Body))
			}
			p.save(c)
		}

		<-forever
	}()
}

func (p *Poll) save(c candidate) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := p.collection.InsertOne(ctx, c)
	if err != nil {
		byt, _ := json.Marshal(c)
		p.rm.Enqueue(string(byt))
	}
}

func (p *Poll) initMongo() {
	mongoUrl := os.Getenv("MONGO_URL")
	database := os.Getenv("MONGO_DATABASE")
	collection := os.Getenv("MONGO_COLLECTION")

	connOptions := options.Client().ApplyURI(mongoUrl)
	conn, err := mongo.Connect(context.Background(), connOptions)

	if err != nil {
		log.Printf("%q", err)
	}

	client := conn.Database(database).Collection(collection)

	p.conn = conn
	p.collection = client
}

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

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
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
	metric     prometheus.Counter
}

type candidate struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type PollResult struct {
	Name  string `json:"name"`
	Votes int    `json:"votes"`
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
	var results []bson.M

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := p.collection.Aggregate(ctx, bson.A{
		bson.D{
			{"$count", "total"},
		},
	})

	if err != nil {
		log.Printf("%q", err)
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		log.Printf("%q", err)
	}

	response.Header().Add("content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	req, _ := json.Marshal(results[0])
	fmt.Fprintln(response, string(req))
}

func (p *Poll) GetVotesByCandidate(response http.ResponseWriter, request *http.Request) {
	var results []bson.M
	var result []PollResult

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := p.collection.Aggregate(ctx, bson.A{
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$name"},
					{"total", bson.D{{"$sum", 1}}},
				},
			},
		},
	})
	if err != nil {
		log.Printf("%q", err)
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		log.Printf("%q", err)
	}

	for _, c := range results {
		v := fmt.Sprint(c["total"])
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Println("Error converting value from database")
		}
		r := PollResult{Name: fmt.Sprint(c["_id"]), Votes: int(value)}
		result = append(result, r)
	}

	response.Header().Add("content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	req, _ := json.Marshal(result)
	fmt.Fprintln(response, string(req))
}

func Init(rm queue.RabbitMQ) (p Poll) {
	p.rm = rm
	p.initMongo()

	p.metric = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "votes_count",
			Help: "No of request handled by poll server",
		},
	)

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
			c.CreatedAt = time.Now().Local()
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

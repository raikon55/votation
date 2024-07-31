package poll

import (
	"net/http"
	"net/http/httptest"
	"os"
	"paredao/src/poll"
	"paredao/src/queue"
	"testing"
)

func TestCreateVote(t *testing.T) {
	t.Run("create a vote", func(t *testing.T) {
		os.Setenv("RABBITMQ_URL", "amqp://test:test@localhost:5672/paredao")

		queueName := "test-queue"
		rmq := queue.InitRabbitMQ(queueName)
		p := poll.Init(rmq)

		request, _ := http.NewRequest(http.MethodPost, "/vote", nil)
		response := httptest.NewRecorder()

		p.CreateVote(response, request)
		rmq.Close()

		got := response.Result().StatusCode
		want := 201

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

package queue

import (
	"os"
	"paredao/src/queue"
	"runtime"
	"sync"
	"testing"
)

func TestEnqueue(t *testing.T) {
	t.Run("enqueue a message", func(t *testing.T) {
		os.Setenv("RABBITMQ_URL", "amqp://test:test@localhost:5672/paredao")

		message := "{message: Hello World!}"
		queueName := "test-queue"

		runtime.GOMAXPROCS(8)
		var wg sync.WaitGroup

		for i := 0; i < 2_000; i++ {
			wg.Add(1)
			go func() {
				queue.Enqueue(message, queueName)
				wg.Done()
			}()
		}

		wg.Wait()
	})
}

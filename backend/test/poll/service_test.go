package poll

import (
	"net/http"
	"net/http/httptest"
	"paredao/src/poll"
	"testing"
)

func TestCreateVote(t *testing.T) {
	t.Run("create a vote", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/vote", nil)
		response := httptest.NewRecorder()

		poll.CreateVote(response, request)

		got := response.Result().StatusCode
		want := 201

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

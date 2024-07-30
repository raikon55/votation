package poll

import (
	"net/http"
)

func CreateVote(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusCreated)
}

func GetVotes(response http.ResponseWriter, request *http.Request) {

}

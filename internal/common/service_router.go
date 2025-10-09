package common

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mbilarusdev/durak_network/network"
)

func NewServiceRouter() *mux.Router {
	return mux.NewRouter()
}

func AddRoute(
	router *mux.Router,
	path string,
	handler func(w http.ResponseWriter, r *http.Request) *network.Result,
	methods ...string,
) {
	router.HandleFunc(path, network.Handler(handler)).Methods(methods...)
}

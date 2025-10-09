package endpoint

import (
	"net/http"

	"github.com/mbilarusdev/durak_network/network"
)

type Endpoint interface {
	Call(w http.ResponseWriter, r *http.Request) *network.Result
}

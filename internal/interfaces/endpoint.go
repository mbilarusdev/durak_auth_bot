package interfaces

import (
	"net/http"
)

type Endpoint interface {
	Call(w http.ResponseWriter, r *http.Request) error
}

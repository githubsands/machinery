package observability

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewObserver() *Observer {
	return &Observer{
		h: promhttp.Handler(),
	}
}

type Observer struct {
	h http.Handler
}

func (o *Observer) Run() {
	http.ListenAndServe(":2112", nil)
}

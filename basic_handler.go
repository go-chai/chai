package chai

import (
	"net/http"
)

func NewBasicHandler(h http.HandlerFunc) *BasicHandler {
	return &BasicHandler{
		f: h,
	}
}

type BasicHandler struct {
	f http.HandlerFunc
}

func (h *BasicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.f.ServeHTTP(w, r)
}

func (h *BasicHandler) Handler() any {
	return h.f
}

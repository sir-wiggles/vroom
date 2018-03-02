package http

import (
	"net/http"

	"github.com/sir-wiggles/arc/pkg/webstore"
)

type Handler struct {
	UserService *webstore.UserService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

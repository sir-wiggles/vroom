package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sir-wiggles/arc/pkg/webstore"
)

type Handler struct {
	mux *chi.Mux

	PhoneService webstore.PhoneService
	CacheService webstore.CacheService
}

func NewHandler() *Handler {
	h := &Handler{mux: chi.NewMux()}
	h.mux.Use(middleware.Recoverer)
	h.mux.Use(middleware.RequestID)

	h.mux.Post("/phone/sms/code", h.phoneSMSCode)
	h.mux.Post("/phone/sms/verify", h.phoneSMSVerify)

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) phoneSMSCode(w http.ResponseWriter, r *http.Request) {
	h.PhoneService.SendCode("asdf")
}

func (h *Handler) phoneSMSVerify(w http.ResponseWriter, r *http.Request) {
	h.PhoneService.VerifyCode("asdf", 1234)
}

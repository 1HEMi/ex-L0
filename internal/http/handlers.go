package http

import (
	"Ex-L0/internal/cache"
	"Ex-L0/internal/logger"
	"Ex-L0/internal/repository"
	"Ex-L0/internal/service"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type Handlers struct {
	Log           *logger.Logger
	OrdersService *service.OrdersService
	Cache         *cache.Cache
}

func (h *Handlers) HandlerGetOrderByID(w http.ResponseWriter, r *http.Request) {

	uid := mux.Vars(r)["order_uid"]

	if uid == "" {
		writeError(w, http.StatusBadRequest, "empty order id")
		return

	}

	if o, ok := h.Cache.Get(uid); ok {
		writeJSON(w, http.StatusOK, o)
	}

	o, err := h.OrdersService.GetByID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "order not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "server error")
		return
	}
	h.Cache.Set(uid, o)
	writeJSON(w, http.StatusOK, o)

}

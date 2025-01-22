package message

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/DamirAbd/HLA-HW1/services/auth"
	"github.com/DamirAbd/HLA-HW1/types"
	"github.com/DamirAbd/HLA-HW1/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.MessageStore
	userStore types.UserStore
}

func NewHandler(store types.MessageStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/dialog/{userID}/send", auth.WithJWTAuth(h.handleCreateMessage, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/dialog/{userID}/list", auth.WithJWTAuth(h.handleGetMessages, h.userStore)).Methods(http.MethodGet)
}

func (h *Handler) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	var payload types.PostPayload

	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var Message types.Message

	Message.From = auth.GetUserIDFromContext(r.Context())
	Message.To = str
	Message.Message = payload.Post //use same posts

	errm := h.store.CreateMessage(Message)
	if errm != nil {
		utils.WriteError(w, http.StatusInternalServerError, errm)
		return
	}

}

func (h *Handler) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	var Message types.Message

	Message.From = auth.GetUserIDFromContext(r.Context())
	Message.To = str

	msgs, errm := h.store.GetMessages(Message.From, Message.To)
	if errm != nil {
		utils.WriteError(w, http.StatusInternalServerError, errm)
		return
	}
	utils.WriteJSON(w, http.StatusOK, msgs)
}

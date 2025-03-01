package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/DamirAbd/HLA-HW1/services/auth"
	"github.com/DamirAbd/HLA-HW1/types"
	"github.com/DamirAbd/HLA-HW1/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/user/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/user/get/{userID}", auth.WithJWTAuth(h.handleGetUser, h.store)).Methods(http.MethodGet)
	router.HandleFunc("/user/search", auth.WithJWTAuth(h.handleSearchUser, h.store)).Methods(http.MethodGet)
	router.HandleFunc("/friend/set/{userID}", auth.WithJWTAuth(h.handlePutFriend, h.store)).Methods(http.MethodPut)
	router.HandleFunc("/friend/delete/{userID}", auth.WithJWTAuth(h.handleDeleteFriend, h.store)).Methods(http.MethodPut)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err := h.store.GetUserByID(user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid ID or password"))
		return

	}

	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid ID or password"))
		return
	}

	secret := []byte("very-secret-string")
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// how to check if user exists ?

	// hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	ID := uuid.New().String()

	err = h.store.CreateUser(types.User{
		ID:         ID,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		BirthDate:  user.BirthDate,
		Biography:  user.Biography,
		City:       user.City,
		Password:   hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"user_id": ID})

}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	user, err := h.store.GetUserByID(str)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) handleSearchUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fname, ok := query["first_name"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing FirstName"))
		return
	}

	lname, ok := query["last_name"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing LastName"))
		return
	}

	user, err := h.store.GetUsersByName(strings.Join(fname, ""), strings.Join(lname, ""))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) handlePutFriend(w http.ResponseWriter, r *http.Request) {

	//TODO: add unick pairs user+friend;

	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	_, err := h.store.GetUserByID(str)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("friend not found"))
	}

	ctx := r.Context()
	strid := auth.GetUserIDFromContext(ctx)

	//write to table friends

	err1 := h.store.SetFriend(strid, str)
	if err1 != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("friend not found"))
	}

}

func (h *Handler) handleDeleteFriend(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	ctx := r.Context()
	strid := auth.GetUserIDFromContext(ctx)

	//write to table friends

	err1 := h.store.DeleteFriend(strid, str)
	if err1 != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("friend not found"))
	}

}

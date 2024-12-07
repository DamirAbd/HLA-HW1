package post

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/DamirAbd/HLA-HW1/services/auth"
	"github.com/DamirAbd/HLA-HW1/types"
	"github.com/DamirAbd/HLA-HW1/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.PostStore
	userStore types.UserStore
}

func NewHandler(store types.PostStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/post/create", auth.WithJWTAuth(h.handleCreatePost, h.userStore)).Methods((http.MethodPost))
	router.HandleFunc("/post/update", auth.WithJWTAuth(h.handleUpdatePost, h.userStore)).Methods((http.MethodPut))
	router.HandleFunc("/post/delete/{postID}", auth.WithJWTAuth(h.handleDeletePost, h.userStore)).Methods((http.MethodPut))
	router.HandleFunc("/post/get/{postID}", auth.WithJWTAuth(h.handleGetPost, h.userStore)).Methods((http.MethodGet))
	router.HandleFunc("/post/feed", auth.WithJWTAuth(h.handleFeed, h.userStore)).Methods((http.MethodGet))

}

func (h *Handler) handleCreatePost(w http.ResponseWriter, r *http.Request) {

	var post types.Post
	var payload types.PostPayload

	post.AutorId = auth.GetUserIDFromContext(r.Context())

	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	post.ID = uuid.New().String()

	post.Post = payload.Post

	err = h.store.CreatePost(types.Post{
		ID:      post.ID,
		AutorId: post.AutorId,
		Post:    post.Post,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, post.ID)
}

func (h *Handler) handleGetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["postID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing post ID"))
		return
	}

	post, err := h.store.GetPostByID(str)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, post)

}

func (h *Handler) handleUpdatePost(w http.ResponseWriter, r *http.Request) {

	type payload struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}

	var postUpdate payload

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &postUpdate)

	err1 := h.store.UpdatePost(postUpdate.ID, postUpdate.Text)
	if err1 != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "OK")

}

func (h *Handler) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["postID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing post ID"))
		return
	}

	err := h.store.DeletePost(str)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNotFound, "")
}

func (h *Handler) handleFeed(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the user's friends
	friends, err := h.userStore.GetFriends(userID)
	if err != nil {
		http.Error(w, "Failed to fetch friends", http.StatusInternalServerError)
		return
	}

	// Collect IDs of friends for batch fetching posts
	var friendIDs []string
	for _, friend := range friends {
		friendIDs = append(friendIDs, friend.ID)
	}

	// Get posts for all friends in one query
	feedPosts, err := h.store.GetPostsByUsers(friendIDs) // Optimized to batch fetch posts
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	/*
		var feedPosts []*types.Post
		for _, value := range friends {
			feedPost, _ := h.store.GetPostsByUser(value.ID)
			for _, post := range feedPost {
				feedPosts = append(feedPosts, post)
			}
		}
	*/

	utils.WriteJSON(w, http.StatusOK, feedPosts)
}

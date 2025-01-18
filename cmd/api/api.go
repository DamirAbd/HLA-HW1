package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/DamirAbd/HLA-HW1/cache"
	"github.com/DamirAbd/HLA-HW1/services/post"
	"github.com/DamirAbd/HLA-HW1/services/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	//initialize cache
	redisCache := cache.NewRedisCache()

	postStore := post.NewStore(s.db)
	postStoreHandler := post.NewHandler(postStore, userStore, redisCache)
	postStoreHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}

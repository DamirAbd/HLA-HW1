package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/DamirAbd/HLA-HW1/cache"
	"github.com/DamirAbd/HLA-HW1/services/message"
	"github.com/DamirAbd/HLA-HW1/services/post"
	"github.com/DamirAbd/HLA-HW1/services/user"
	"github.com/DamirAbd/HLA-HW1/websockets"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
	cdb  *sql.DB
}

func NewAPIServer(addr string, db *sql.DB, cdb *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
		cdb:  cdb,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	manager := websockets.NewManager()
	router.PathPrefix("/feed").Handler(http.StripPrefix("/feed", http.FileServer(http.Dir("/go/src/api/frontend"))))
	router.Handle("/ws", http.HandlerFunc(manager.ServeWS))
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	//initialize cache
	redisCache := cache.NewRedisCache()

	postStore := post.NewStore(s.db)
	postStoreHandler := post.NewHandler(postStore, userStore, redisCache)
	postStoreHandler.RegisterRoutes(subrouter)

	messageStore := message.NewStore(s.cdb)
	messageStoreHandler := message.NewHandler(messageStore, userStore)
	messageStoreHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}

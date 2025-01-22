package message

import (
	"database/sql"

	"github.com/DamirAbd/HLA-HW1/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateMessage(message types.Message) error {
	_, err := s.db.Exec("INSERT INTO messages (sender, recipient, message) VALUES ($1, $2, $3)", message.From, message.To, message.Message)
	if err != nil {
		return err
	}

	return nil
}

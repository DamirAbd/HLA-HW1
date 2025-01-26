package message

import (
	"database/sql"
	"log"
	"sort"
	"strings"

	"github.com/DamirAbd/HLA-HW1/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateMessage(message types.Message) error {
	var key = []string{message.From, message.To}
	sort.Strings(key)
	c_key := strings.Join(key, "")
	_, err := s.db.Exec("INSERT INTO messages (sender, recipient, message, citus_key) VALUES ($1, $2, $3, $4)", message.From, message.To, message.Message, c_key)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetMessages(sender string, recipient string) ([]*types.Message, error) {
	rows, err := s.db.Query(`
	SELECT m.sender, m.recipient, m.message 
	FROM messages m
	WHERE m.sender = $1
	AND m.recipient = $2`, sender, recipient)

	var msgs []*types.Message

	for rows.Next() {
		m := new(types.Message)
		if err := rows.Scan(&m.From, &m.To, &m.Message); err != nil {
			log.Fatal(err)
		}
		msgs = append(msgs, m)
	}

	if err != nil {
		return nil, err
	}

	return msgs, nil
}

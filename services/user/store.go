package user

import (
	"database/sql"
	"log"

	"github.com/DamirAbd/HLA-HW1/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (ID, FirstName, SecondName, BirthDate, Biography, City,Password) VALUES ( $1, $2, $3, $4, $5, $6, $7)", user.ID, user.FirstName, user.SecondName, user.BirthDate, user.Biography, user.City, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByID(id string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = $1", id)
	u := new(types.User)

	for rows.Next() {
		var userid int64

		if err := rows.Scan(&userid, &u.ID, &u.FirstName, &u.SecondName, &u.BirthDate, &u.Biography, &u.City, &u.Password, &u.CreatedAt); err != nil {
			log.Fatal(err)
		}
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

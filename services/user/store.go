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
	rows, err := s.db.Query(`
	SELECT u.userid
	,u.ID
	,u.FirstName
	,u.SecondName
	,u.BirthDate
	,COALESCE(u.Biography,'')
	,COALESCE(u.City,'')
	,u.password
	,u.createdat
	FROM users u
	WHERE id = $1`, id)
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

func (s *Store) GetUsersByName(fname string, lname string) ([]*types.UserForm, error) {
	rows, err := s.db.Query(`
		SELECT u.ID
			,u.FirstName
			,u.SecondName
			,u.BirthDate
			,COALESCE(u.Biography,'')
			,COALESCE(u.City,'')
		FROM users u
		WHERE FirstName LIKE $1 || '%' AND secondname LIKE $2 || '%'
		ORDER BY ID DESC`, fname, lname)
	users := make([]*types.UserForm, 0)

	for rows.Next() {
		u := new(types.UserForm)

		err := rows.Scan(&u.ID, &u.FirstName, &u.SecondName, &u.BirthDate, &u.Biography, &u.City)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Store) SetFriend(userId string, friendId string) error {
	_, err := s.db.Exec(`INSERT INTO friends (user_id, friend_id) VALUES ( $1, $2)`, userId, friendId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteFriend(userId string, friendId string) error {
	_, err := s.db.Exec(`DELETE FROM friends WHERE user_id = $1 AND friend_id = $2`, userId, friendId)
	if err != nil {
		return err
	}

	return nil
}

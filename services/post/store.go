package post

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

func (s *Store) GetPostByID(id string) (*types.Post, error) {
	rows, err := s.db.Query(`
	SELECT p.post_id, p.post, p.author_id 
	FROM posts p
	WHERE p.post_id = $1`, id)
	p := new(types.Post)

	for rows.Next() {
		if err := rows.Scan(&p.ID, &p.Post, &p.AutorId); err != nil {
			log.Fatal(err)
		}
	}

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Store) UpdatePost(id string, post string) error {
	_, err := s.db.Exec(`
	UPDATE posts
	SET post = $1
	WHERE post_id = $2`, post, id)
	if err != nil {
		return err
	}

	return nil

}

func (s *Store) CreatePost(post types.Post) error {
	_, err := s.db.Exec("INSERT INTO posts (post_id, post, author_id) VALUES ($1, $2, $3)", post.ID, post.Post, post.AutorId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeletePost(id string) error {
	_, err := s.db.Exec(`
	DELETE FROM posts
	WHERE post_id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

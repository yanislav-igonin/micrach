package repositories

import (
	"context"
	Db "micrach/db"

	"github.com/jackc/pgx/v4"
)

type PostsRepository struct{}

var Posts PostsRepository

func (r *PostsRepository) Get(limit, offset int) ([]Post, error) {
	return PostsDb, nil
}

func (r *PostsRepository) Create(p Post) (int, error) {
	sql := `
		INSERT INTO posts (is_parent, parent_id, title, text, is_sage)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	conn, err := Db.Pool.Acquire(context.TODO())
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	var row pgx.Row
	if p.IsParent {
		row = conn.QueryRow(
			context.TODO(), sql, p.IsParent, nil, p.Title, p.Text, p.IsSage,
		)
	} else {
		row = conn.QueryRow(
			context.TODO(), sql, p.IsParent, p.ParentID, p.Title, p.Text, p.IsSage,
		)
	}
	if err != nil {
		return 0, nil
	}

	createdPost := new(Post)
	row.Scan(&createdPost.ID)

	return createdPost.ID, nil
}

// func (r *PostsRepository) GetByID() int {
// }

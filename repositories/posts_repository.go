package repositories

import (
	"context"
	Db "micrach/db"

	"github.com/jackc/pgx/v4"
)

type PostsRepository struct{}

var Posts PostsRepository

func (r *PostsRepository) Get(limit, offset int) ([]Post, error) {
	if Db.Pool == nil {
		return PostsDb, nil
	}

	sql := `
		SELECT id, title, text, created_at
		FROM posts
		WHERE
			is_parent = true
			AND is_deleted = false
		ORDER BY updated_at DESC
		LIMIT $1
	`

	conn, err := Db.Pool.Acquire(context.TODO())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.TODO(), sql, limit)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Title, &post.Text, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
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

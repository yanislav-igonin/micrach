package repositories

import (
	"context"
	Config "micrach/config"
	Db "micrach/db"
	"time"

	"github.com/jackc/pgx/v4"
)

type PostsRepository struct{}

var Posts PostsRepository

func (r *PostsRepository) Get(limit, offset int) ([]Post, error) {
	sql := `
		SELECT id, title, text, created_at
		FROM posts
		WHERE
			is_parent = true
			AND is_deleted != true
		ORDER BY updated_at DESC
		OFFSET $1
		LIMIT $2
	`

	rows, err := Db.Pool.Query(context.TODO(), sql, offset, limit)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	postsMap := make(map[int]Post)
	var postIDs []int
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Title, &post.Text, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		postsMap[post.ID] = post
		postIDs = append(postIDs, post.ID)
	}

	filesMap, err := Files.GetByPostIDs(postIDs)
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, postID := range postIDs {
		post := postsMap[postID]
		post.Files = filesMap[postID]
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostsRepository) GetCount() (int, error) {
	sql := `
		SELECT COUNT(*)
		FROM posts
		WHERE 
			is_parent = true
			AND is_deleted != true
			AND is_archived != true
	`

	row := Db.Pool.QueryRow(context.TODO(), sql)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostsRepository) Create(p Post) (int, error) {
	sql := `
		INSERT INTO posts (is_parent, parent_id, title, text, is_sage, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var row pgx.Row
	if p.IsParent {
		row = Db.Pool.QueryRow(
			context.TODO(), sql, p.IsParent, nil, p.Title, p.Text, p.IsSage, time.Now(),
		)
	} else {
		row = Db.Pool.QueryRow(
			context.TODO(), sql, p.IsParent, p.ParentID, p.Title, p.Text, p.IsSage, nil,
		)
	}

	createdPost := new(Post)
	err := row.Scan(&createdPost.ID)
	if err != nil {
		return 0, err
	}

	return createdPost.ID, nil
}

func (r *PostsRepository) GetThreadByPostID(ID int) ([]Post, error) {
	sql := `
		SELECT 
			id,
			title,
			text,
			is_sage,
			created_at,
			is_parent,
			parent_id
		FROM posts
		WHERE
				(id = $1 AND is_parent = true AND is_deleted != true)
				OR (parent_id = $1 AND is_deleted != true)
		ORDER BY created_at ASC
	`

	rows, err := Db.Pool.Query(context.TODO(), sql, ID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	postsMap := make(map[int]Post)
	var postIDs []int
	for rows.Next() {
		var post Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Text,
			&post.IsSage,
			&post.CreatedAt,
			&post.IsParent,
			&post.ParentID,
		)
		if err != nil {
			return nil, err
		}

		postsMap[post.ID] = post
		postIDs = append(postIDs, post.ID)
	}

	filesMap, err := Files.GetByPostIDs(postIDs)
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, postID := range postIDs {
		post := postsMap[postID]
		post.Files = filesMap[postID]
		posts = append(posts, post)
	}

	return posts, nil
}

// Check if thread is archived
func (r *PostsRepository) GetIfThreadIsArchived(ID int) (bool, error) {
	sql := `
		SELECT is_archived
		FROM posts
		WHERE id = $1
		LIMIT 1
	`

	row := Db.Pool.QueryRow(context.TODO(), sql, ID)
	var isArchived bool
	err := row.Scan(&isArchived)
	if err != nil {
		return false, err
	}
	return isArchived, nil
}

func (r *PostsRepository) CreateInTx(tx pgx.Tx, p Post) (int, error) {
	sql := `
		INSERT INTO posts (is_parent, parent_id, title, text, is_sage, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var row pgx.Row
	if p.IsParent {
		row = tx.QueryRow(
			context.TODO(), sql, p.IsParent, nil, p.Title, p.Text, p.IsSage, time.Now(),
		)
	} else {
		row = tx.QueryRow(
			context.TODO(), sql, p.IsParent, p.ParentID, p.Title, p.Text, p.IsSage, nil,
		)
	}

	createdPost := new(Post)
	err := row.Scan(&createdPost.ID)
	if err != nil {
		return 0, err
	}

	return createdPost.ID, nil
}

func (r *PostsRepository) GetOldestThreadUpdatedAt() (time.Time, error) {
	sql := `
		SELECT updated_at
		FROM posts
		WHERE 
			is_parent = true
			AND is_deleted != true
			AND is_archived != true
		ORDER BY updated_at DESC
		OFFSET $1 - 1
		LIMIT 1
	`

	row := Db.Pool.QueryRow(context.TODO(), sql, Config.App.ThreadsMaxCount)
	var updatedAt time.Time
	err := row.Scan(&updatedAt)
	if err != nil {
		return time.Time{}, err
	}
	return updatedAt, nil
}

func (r *PostsRepository) ArchiveThreadsFrom(t time.Time) error {
	sql := `
		UPDATE posts
		SET is_archived = true
		WHERE 
			is_parent = true
			AND is_archived != true
			AND updated_at <= $1
	`

	_, err := Db.Pool.Exec(context.TODO(), sql, t)
	return err
}

// Returns count of posts in thread by thread ID
func (r *PostsRepository) GetThreadPostsCount(id int) (int, error) {
	sql := `
		SELECT COUNT(*)
		FROM posts
		WHERE
			(id = $1 AND is_parent = true AND is_deleted != true)
			OR (parent_id = $1 AND is_deleted != true)
	`

	row := Db.Pool.QueryRow(context.TODO(), sql, id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Updates threads updated at time by thread ID
func (r *PostsRepository) BumpThreadInTx(tx pgx.Tx, id int) error {
	sql := `
		UPDATE posts
		SET updated_at = now()
		WHERE id = $1
	`

	_, err := tx.Query(context.TODO(), sql, id)
	return err
}

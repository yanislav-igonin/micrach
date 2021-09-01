package repositories

import (
	"context"
	Db "micrach/db"
)

type FilesRepository struct{}

var Files FilesRepository

func (r *FilesRepository) Create(f File) error {
	conn, err := Db.Pool.Acquire(context.TODO())
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := `
		INSERT INTO files (post_id, name, ext, size)
		VALUES ($1, $2, $3, $4)
	`

	conn.QueryRow(
		context.TODO(), sql, f.PostID, f.Name, f.Ext, f.Size,
	)

	return nil
}

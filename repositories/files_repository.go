package repositories

import (
	"context"
	Db "micrach/db"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
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

	row := conn.QueryRow(
		context.TODO(), sql, f.PostID, f.Name, f.Ext, f.Size,
	)

	err = row.Scan()
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	return nil
}

func (r *FilesRepository) GetByPostIDs(postIDs []int) (map[int][]File, error) {
	sql := `
		SELECT post_id, name, size
		FROM files
		WHERE post_id = ANY ($1)
		ORDER BY id ASC
	`

	ids := &pgtype.Int4Array{}
	ids.Set(postIDs)

	rows, err := Db.Pool.Query(context.TODO(), sql, ids)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	filesMapByPostId := make(map[int][]File)
	for rows.Next() {
		var file File
		err = rows.Scan(&file.PostID, &file.Name, &file.Size)
		if err != nil {
			return nil, err
		}

		if filesMapByPostId[file.PostID] == nil {
			filesMapByPostId[file.PostID] = []File{}
		}

		filesMapByPostId[file.PostID] = append(
			filesMapByPostId[file.PostID],
			file,
		)
	}

	return filesMapByPostId, nil
}

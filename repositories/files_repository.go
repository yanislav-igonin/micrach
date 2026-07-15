package repositories

import (
	"context"
	Db "micrach/db"

	"github.com/jackc/pgx/v5"
)

type FilesRepository struct{}

var Files FilesRepository

func toInt32s(values []int) []int32 {
	if values == nil {
		return nil
	}

	result := make([]int32, len(values))
	for i, value := range values {
		result[i] = int32(value)
	}
	return result
}

func (r *FilesRepository) Create(ctx context.Context, f File) error {
	sql := `
		INSERT INTO files (post_id, name, ext, size)
		VALUES ($1, $2, $3, $4)
	`

	row := Db.Pool.QueryRow(
		ctx, sql, f.PostID, f.Name, f.Ext, f.Size,
	)

	err := row.Scan()
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	return nil
}

func (r *FilesRepository) GetByPostIDs(ctx context.Context, postIDs []int) (map[int][]File, error) {
	sql := `
		SELECT id, post_id, name, size, ext
		FROM files
		WHERE post_id = ANY ($1)
		ORDER BY id ASC
	`

	rows, err := Db.Pool.Query(ctx, sql, toInt32s(postIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	filesMapByPostId := make(map[int][]File)
	for rows.Next() {
		var file File
		err = rows.Scan(&file.ID, &file.PostID, &file.Name, &file.Size, &file.Ext)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return filesMapByPostId, nil
}

func (r *FilesRepository) CreateInTx(ctx context.Context, tx pgx.Tx, f File) (int, error) {
	sql := `
		INSERT INTO files (post_id, name, ext, size)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	row := tx.QueryRow(
		ctx, sql, f.PostID, f.Name, f.Ext, f.Size,
	)

	createdFile := new(File)
	err := row.Scan(&createdFile.ID)
	if err != nil && err != pgx.ErrNoRows {
		return 0, err
	}

	return createdFile.ID, nil
}

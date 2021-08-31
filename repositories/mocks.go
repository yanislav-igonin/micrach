package repositories

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getFile(id, postId int, name string) File {
	return File{
		ID:        id,
		PostID:    postId,
		Name:      name,
		Ext:       "image/jpeg",
		Size:      10000,
		CreatedAt: time.Now(),
	}
}

func getPost(id int) Post {
	return Post{
		ID:        id,
		IsParent:  true,
		ParentID:  0,
		IsDeleted: false,
		Title:     randSeq(rand.Intn(100)),
		Text:      randSeq(rand.Intn(100)),
		IsSage:    false,
		Files: []File{
			getFile(2, id, "https://memepedia.ru/wp-content/uploads/2018/03/ebanyy-rot-etogo-kazino.png"),
			getFile(1, id, "https://memepedia.ru/wp-content/uploads/2018/03/ebanyy-rot-etogo-kazino.png"),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

var PostsDb = []Post{}

func SeedMocks() {
	rand.Seed(time.Now().UnixNano())

	for i := 1; i < 10; i++ {
		PostsDb = append(PostsDb, getPost(i))
	}
}

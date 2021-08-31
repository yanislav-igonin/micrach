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

func getPost(id, threadID int) Post {
	return Post{
		ID:       id,
		ThreadID: threadID,
		Title:    randSeq(rand.Intn(100)),
		Text:     randSeq(rand.Intn(100)),
		IsSage:   false,
		Files: []File{
			getFile(2, id, "https://memepedia.ru/wp-content/uploads/2018/03/ebanyy-rot-etogo-kazino.png"),
			getFile(1, id, "https://memepedia.ru/wp-content/uploads/2018/03/ebanyy-rot-etogo-kazino.png"),
		},
		CreatedAt: time.Now(),
	}
}

func getThread(id int) Thread {
	return Thread{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Posts: []Post{
			getPost(1, id),
			getPost(1, id),
		},
	}
}

var ThreadsDb = []Thread{}

func SeedMocks() {
	rand.Seed(time.Now().UnixNano())

	for i := 1; i < 100; i++ {
		ThreadsDb = append(ThreadsDb, getThread(i))
	}
}

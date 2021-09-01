package repositories

import (
	"log"
	"math/rand"
	Db "micrach/db"
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

func getPost(id int, pid *int) Post {
	var parentID int
	if pid == nil {
		parentID = 0
	} else {
		parentID = *pid
	}

	var isParent bool
	if parentID == 0 {
		isParent = true
	} else {
		isParent = false
	}
	return Post{
		ID:        id,
		IsParent:  isParent,
		ParentID:  parentID,
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

func seedLocal() {
	rand.Seed(time.Now().UnixNano())

	for i := 1; i < 10; i++ {
		PostsDb = append(PostsDb, getPost(i, nil))
	}
}

func seedDb() {
	var posts []Post
	for i := 1; i < 10; i++ {
		post := getPost(i, nil)
		posts = append(posts, post)
	}

	for _, parentPost := range posts {
		for i := 0; i < 10; i++ {
			childPost := getPost(parentPost.ID*10+i, &parentPost.ID)
			posts = append(posts, childPost)
		}
	}

	for _, post := range posts {
		_, err := Posts.Create(post)
		if err != nil {
			log.Panicln(err)
		}
	}
}

func Seed() {
	if Db.Pool != nil {
		seedDb()
	} else {
		seedLocal()
	}
	log.Println("mocks - online")
}

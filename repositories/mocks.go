package repositories

import "time"

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
		Title:    "Basic Title",
		Text:     "Basic Text",
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
	for i := 1; i < 10; i++ {
		ThreadsDb = append(ThreadsDb, getThread(i))
	}
}

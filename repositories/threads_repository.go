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
			getFile(2, id, "Screenshot 2020-08-14 at 23.17.29.png"),
			getFile(1, id, "maxresdefault.jpg"),
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

var ThreadsDb = []Thread{
	getThread(1),
	getThread(2),
	getThread(3),
}

type ThreadsRepository struct{}

var Threads ThreadsRepository

func (r *ThreadsRepository) Get() []Thread {
	return ThreadsDb
}

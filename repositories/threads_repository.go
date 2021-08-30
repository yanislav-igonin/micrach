package repositories

type ThreadsRepository struct{}

var Threads ThreadsRepository

func (r *ThreadsRepository) Get(limit, offset int) ([]Thread, error) {
	// conn, err := Db.Pool.Acquire(context.TODO())
	// if err != nil {
	// 	return nil, err
	// }
	// defer conn.Release()

	// sql := `
	// 	SELECT
	// 		id
	// 	FROM threads
	// 	WHERE is_deleted != true
	// 	ORDER BY updated_at DESC
	// `

	// rows, err := conn.Query(context.TODO(), sql)
	// if err != nil {
	// 	return nil, err
	// }

	// var threads []Thread
	// var threadsIDs []int
	// for rows.Next() {
	// 	var thread Thread
	// 	err = rows.Scan(&thread.ID)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	threadsIDs = append(threadsIDs, thread.ID)
	// 	threads = append(threads, thread)
	// }
	// rows.Close()

	// sql = `
	// 	SELECT
	// 		id,
	// 		thread_id,
	// 		title,
	// 		text,
	// 		is_sage,
	// 		created_at
	// 	FROM posts
	// 	WHERE thread_id IN (
	// `
	// buf := bytes.NewBufferString(sql)
	// for i, tid := range threadsIDs {
	// 	buf.WriteString(strconv.Itoa(tid))
	// 	if i != len(threadsIDs)-1 {
	// 		buf.WriteString(",")
	// 	}
	// }
	// sql = buf.String() + ") ORDER BY created_at ASC"
	// rows, err = conn.Query(context.TODO(), sql)
	// if err != nil {
	// 	return nil, err
	// }

	// var posts []Post
	// var postIDs []int
	// for rows.Next() {
	// 	var post Post
	// 	err = rows.Scan(&post.ID, &post.ThreadID, &post.Title, &post.Text, &post.IsSage, &post.CreatedAt)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	postIDs = append(postIDs, post.ID)
	// 	posts = append(posts, post)
	// }
	// rows.Close()

	return ThreadsDb, nil
	// return nil, nil
}

// func (r *ThreadsRepository) Create(post Post) int {

// }

// func (r *ThreadsRepository) GetByID() int {
// 	newThreadID := time.Now().Second()
// 	ThreadsDb = append(ThreadsDb, getThread(newThreadID))
// 	return newThreadID
// }

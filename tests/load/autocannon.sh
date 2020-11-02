# Create Thread
# POST /threads
autocannon -c 1000 -d 60 -r 1000 -m 'POST' -b 'post.json' -H 'Content-Type: application/json' localhost:3000/threads

sleep 10

# Create Post
# POST /threads/:id
autocannon -c 1000 -d 60 -r 1000 -m 'POST' -b 'post.json' -H 'Content-Type: application/json' localhost:3000/threads/1

sleep 10

# Get Threads
# GET /threads
autocannon -c 1000 -d 60 -r 1000 localhost:3000/threads

sleep 10

# Get Thread
# GET /threads/:id
autocannon -c 1000 -d 60 -r 1000 localhost:3000/threads/1

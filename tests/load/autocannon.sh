# GET /threads
autocannon -c 1000 -d 60 -r 1000 localhost:3000/threads

# GET /threads/:id
autocannon -c 1000 -d 60 -r 1000 localhost:3000/threads/19

# POST /threads
autocannon -c 1000 -d 60 -r 1000 -m 'POST' -b 'thread.json' -H 'Content-Type: application/json' localhost:3000/threads
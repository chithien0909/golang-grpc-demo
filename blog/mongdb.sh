# shellcheck disable=SC2046
docker run --name mongod-blog-golang -p 27017:27017 -v $(pwd)/data:/data/db -d mongo

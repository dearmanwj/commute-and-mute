echo "Building image..."
docker build --platform linux/amd64 -t docker-image:test .

docker run -d -p 9000:8080 --name commute-and-mute docker-image:test
enterredis:
	docker exec -it redis_oc redis-cli

compile:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app

composerestart:
	docker-compose down -v
	docker-compose up -d

dockerrun:
	docker build . -t oc_be:1.0.0
	docker run --rm --name oc_be oc_be:1.0.0

dockerclear:
	docker stop oc_be
	docker rm oc_be
	docker rmi oc_be:1.0.0

gotest:
	go test -v ./app

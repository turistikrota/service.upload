build:
	docker build --build-arg GITHUB_USER=${TR_GIT_USER} --build-arg GITHUB_TOKEN=${TR_GIT_TOKEN} -t github.com/turistikrota/service.upload . 

run:
	docker service create --name upload-api-turistikrota-com --network turistikrota --secret jwt_private_key --secret jwt_public_key --env-file .env --publish 6018:6018 github.com/turistikrota/service.upload:latest

remove:
	docker service rm upload-api-turistikrota-com

stop:
	docker service scale upload-api-turistikrota-com=0

start:
	docker service scale upload-api-turistikrota-com=1

restart: remove build run
	
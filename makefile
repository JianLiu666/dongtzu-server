docker_build:
	docker build -t jianliu0616/dongtzu-server:latest .

docker_push:
	docker push jianliu0616/dongtzu-server:latest

deploy_gcp:
	docker build -t asia.gcr.io/dongtzu-20210630/dongtzu-server:latest .
	docker push asia.gcr.io/dongtzu-20210630/dongtzu-server:latest
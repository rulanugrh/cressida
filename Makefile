install_depedency_ubuntu:
	sudo apt install docker docker.io -y

install_depedency_mac:
	sudo brew install docker docker.io -y

build_image:
	sudo docker build -t app-cressida .

run:
	sudo docker compose up -d db-cressida
	sudo docker compose up -d

down:
	sudo docker compose down
install_depedency_ubuntu:
	sudo gpg -k && sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69 && echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list && sudo apt-get update && sudo apt-get install k6 docker docker.io -y

install_depedency_mac:
	sudo brew install docker docker.io k6 -y

build_image:
	sudo docker build -t app-cressida .

load_test_user_register:
	k6 run test/load_test/build/src/user-register.js -a localhost:6556

load_test_post_vehicle:
	k6 run test/load_test/build/src/post-vehicle.js -a localhost:6556

load_test_post_transporter:
	k6 run test/load_test/build/src/post-transporter.js -a localhost:6556

load_test_get_vehicle_by_id:
	k6 run test/load_test/build/src/get-vehicle-by-id.js -a localhost:6556

load_test_get_all_vehicle:
	k6 run test/load_test/build/src/get-all-vehicle.js -a localhost:6556

load_test_get_transporter_by_id:
	k6 run test/load_test/build/src/get-transporter-by-id.js -a localhost:6556

load_test_get_all_transporter:
	k6 run test/load_test/build/src/get-all-transporter.js -a localhost:6556

run:
	sudo docker compose up -d db-cressida
	sudo docker compose up -d

down:
	sudo docker compose down
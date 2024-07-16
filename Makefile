install_depedency_ubuntu:
	sudo apt install docker docker.io -y

install_depedency_mac:
	sudo brew install docker docker.io -y

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
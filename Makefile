run_app:
	./run.sh

run_ipe:
	cd ./ipe
	./ipe

migrateup:
	echo "migrate up"
	soda migrate

migrate_add:
	soda generate sql createHostsTable

migrate_fizz_add:
	soda generate fizz createHostsTable

build_ws:
	docker build -t ipe:1.0.0 -f ./ipe/Dockerfile ./ipe

run_ipe:
	docker run -d --name ipe -p 4001:4001 ipe:1.0.0

run_mailhog:
	mailhog
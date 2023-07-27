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
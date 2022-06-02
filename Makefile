export DBNAME = postgres
export DBUSER = postgres
export DBPASS = postgres

export PGPASSWORD = $(DBPASS)

build-network:
	-docker network create dev

build-devdb:
	docker build -t devdb -f Dockerfile.devdb \
	--build-arg DBNAME=$(DBNAME) \
	--build-arg DBUSER=$(DBUSER) \
	--build-arg DBPASS=$(DBPASS) \
	.

docker-run-devdb: build-network build-devdb
	docker run -d --rm \
	--name devdb \
	--network dev \
	-p 5432:5432 \
	devdb

build-etl:
	docker build -t etl -f Dockerfile.etl .

docker-run-etl: build-network
	docker run --rm \
	--name etl \
	--network dev \
	-e DBHOST=devdb \
	-e DBUSER=$(DBUSER) \
	-e DBPASS=$(DBPASS) \
	etl

build-webserver:
	docker build -t webserver -f Dockerfile.webserver .

build-webserver-heroku:
	docker build -t registry.heroku.com/webserver/web -f Dockerfile.webserver .

docker-run-webserver: build-webserver
	docker run -d --rm \
	--name webserver \
	--network dev \
	-e DBHOST=devdb \
	-e DBPORT=5432 \
	-e DBNAME=$(DBNAME) \
	-e DBUSER=$(DBUSER) \
	-e DBPASS=$(DBPASS) \
	-p 8000:8000 \
	webserver

get-file:
	curl -O https://gist.githubusercontent.com/seahyc/b9ebbe264f8633a1bf167cc6a90d4b57/raw/021d2e0d2c56217bad524119d1c31419b2938505/restaurant_with_menu.json
	curl -O https://gist.githubusercontent.com/seahyc/de33162db680c3d595e955752178d57d/raw/785007bc91c543f847b87d705499e86e16961379/users_with_purchase_history.json

clean-json:
	rm *.json

clean-csv:
	rm *.csv

clean-file: clean-json clean-csv

migrate:
	psql -h $(DBHOST) -U $(DBUSER) \
	-c "\COPY restaurants FROM restaurants.csv DELIMITER ',' csv" \
	-c "\COPY opening_hours FROM opening_hours.csv DELIMITER ',' csv" \
	-c "\COPY menus FROM menus.csv DELIMITER ',' csv" \
	-c "\COPY users FROM users.csv DELIMITER ',' csv" \
	-c "\COPY purchase_histories FROM purchase_histories.csv DELIMITER ',' csv"

run-migrate-csv: get-file start-etl clean-json

run-migrate: get-file start-etl migrate clean-file

docker-run-migrate: build-etl docker-run-etl

start-etl:
	go run cmd/etl/main.go

start-webserver:
	go run cmd/webserver/main.go
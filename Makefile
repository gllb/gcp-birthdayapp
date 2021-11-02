tf-validate: terraform/*
	cd ./terraform; terraform init && terraform validate

tf-clean:
	rm -fr ./terraform/.terraform

start-local-db:
	source ./.env; docker run -p 5432:5432 --name postgres-test -e POSTGRES_PASSWORD=test -d postgres

stop-local-db:
	docker rm postgres-test
go-test:
	source ./.env; \
	cd ./src; \
docker run --rm -v "${PWD}":/usr/src/helloworld \
-w /usr/src/helloworld \
-e DBHOST -e DBPORT -e DBUSER \
-e DBPASSWORD -e DBSSLMODE \
golang:1.17 go env -w GO111MODULE=auto && go test

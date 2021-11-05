tf-validate: terraform/*
	cd ./terraform; terraform init && terraform validate

tf-clean:
	rm -fr ./terraform/.terraform

start-local-db:
	source ./.env; docker run -p 5432:5432 --name postgres-test -e POSTGRES_PASSWORD=test -d postgres

stop-local-db:
	docker stop postgres-test && docker rm postgres-test
go-test:
	docker run --rm -v $(CURDIR)/src:/go/src/birthdayapp \
-w /go/src/birthdayapp \
-e DBHOST -e DBPORT -e DBUSER -e DBNAME \
-e DBPASSWORD -e DBSSLMODE -e GO111MODULE=auto \
golang:1.17 bash -c "go get && go test"

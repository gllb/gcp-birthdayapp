tf-validate: terraform/*
	cd ./terraform; terraform init && terraform validate

tf-clean:
	rm -fr ./terraform/.terraform

go-test:
	cd ./src; docker run --rm -v "${PWD}":/usr/src/helloworld -w /usr/src/helloworld golang:1.17 go env -w GO111MODULE=auto && go test

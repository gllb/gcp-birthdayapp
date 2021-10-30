validate: terraform/*
	cd ./terraform; terraform init && terraform validate

clean:
	rm -fr ./terraform/.terraform

* Terraform
- Create 2 VPC app/db https://github.com/terraform-google-modules/terraform-google-network
- Create a VPC peering https://registry.terraform.io/modules/terraform-google-modules/network/google/latest/submodules/network-peering
- Create a gke in app VPC https://registry.terraform.io/modules/terraform-google-modules/kubernetes-engine/google/latest
- Create a google sql db (pgsql) https://github.com/terraform-google-modules/terraform-google-sql-db/tree/master/modules/postgresql
- Allow connection from app to db on pgsql port

* K8S
Create a deployment of helloworld app able to connect to the db, with a Loadbalancer type service.

* Helloworld app
- Create a test suite
- Implement the failing test
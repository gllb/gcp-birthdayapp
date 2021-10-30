module "vpc-app" {
  source  = "terraform-google-modules/network/google"
  version = "~> 3.0"

  project_id   = var.project_id
  network_name = "app"
  routing_mode = "GLOBAL"

  subnets = [
    {
      subnet_name           = "subnet-app"
      subnet_ip             = "10.10.10.0/24"
      subnet_region         = var.region
    }
  ]

  routes = []
}


module "vpc-db" {
  source  = "terraform-google-modules/network/google"
  version = "~> 3.0"

  project_id   = var.project_id
  network_name = "db"
  routing_mode = "GLOBAL"

  subnets = [
    {
      subnet_name           = "subnet-db"
      subnet_ip             = "10.10.20.0/24"
      subnet_region         = var.region
    }
  ]

  routes = []
}

module "vpc-peering" {
  source        = "terraform-google-modules/network/google//modules/network-peering"
  version       = "~> 3.2.1"

  local_network = module.vpc-app.network_name
  peer_network  = module.vpc-db.network_name
}

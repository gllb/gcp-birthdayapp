module "sql-db" {
  source  = "GoogleCloudPlatform/sql-db/google//modules/postgresql"
  version = "8.0.0"

  name                 = "${var.db_name}-db"
  random_instance_name = true
  project_id           = var.project_id
  database_version     = "POSTGRES_9_6"
  region               = "us-central1"

  // Master configurations
  tier                            = "db-custom-2-13312"
  zone                            = "us-central1-c"
  availability_type               = "REGIONAL"
  maintenance_window_day          = 7
  maintenance_window_hour         = 12
  maintenance_window_update_track = "stable"

  // for testing purpose
  deletion_protection = false

  database_flags = [{ name = "autovacuum", value = "off" }]

  ip_configuration = {
    ipv4_enabled    = true
    require_ssl     = true
    private_network = module.vpc-db.network_id
    authorized_networks = [
      {
        name  = "${var.project_id}-app"
        value = module.vpc-app.subnets_ips
      },
    ]
  }

  backup_configuration = {
    enabled                        = true
    start_time                     = "20:55"
    location                       = null
    point_in_time_recovery_enabled = false
    transaction_log_retention_days = null
    retained_backups               = 10
    retention_unit                 = "COUNT"
  }

  // Read replica configurations
  # read_replica_name_suffix = "-${var.db_name}"
  # read_replicas = [
  #   {
  #     name                = "0"
  #     zone                = "us-central1-a"
  #     tier                = "db-custom-2-13312"
  #     ip_configuration    = local.read_replica_ip_configuration
  #     database_flags      = [{ name = "autovacuum", value = "off" }]
  #     disk_autoresize     = null
  #     disk_size           = null
  #     disk_type           = "PD_HDD"
  #     encryption_key_name = null
  #   },
  #   {
  #     name                = "1"
  #     zone                = "us-central1-b"
  #     tier                = "db-custom-2-13312"
  #     ip_configuration    = local.read_replica_ip_configuration
  #     database_flags      = [{ name = "autovacuum", value = "off" }]
  #     disk_autoresize     = null
  #     disk_size           = null
  #     disk_type           = "PD_HDD"
  #     encryption_key_name = null
  #   },
  #   {
  #     name                = "2"
  #     zone                = "us-central1-c"
  #     tier                = "db-custom-2-13312"
  #     ip_configuration    = local.read_replica_ip_configuration
  #     database_flags      = [{ name = "autovacuum", value = "off" }]
  #     disk_autoresize     = null
  #     disk_size           = null
  #     disk_type           = "PD_HDD"
  #     encryption_key_name = null
  #   },
  # ]

  db_name      = var.db_name
  db_charset   = "UTF8"
  db_collation = "en_US.UTF8"

  user_name     = var.db_user
  user_password = var.db_password
}

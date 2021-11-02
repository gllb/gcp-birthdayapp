resource "kubernetes_ingress" "birthdayapp" {
  wait_for_load_balancer = true
  metadata {
    name = "birthdayapp-ingress"
  }
  spec {
    defaultBackend {
      service {
	name = "birthdayapp-service"
	port {
	  number = 8080
	}
      }
    }
  }
}

resource "kubernetes_service" "birthdayapp" {
  metadata {
    name = "birthdayapp-service"
  }
  spec {
    port {
      port = 8080
      target_port = 80
    }
    selector {
      app = "birthdayapp"
    }
    type = "LoadBalancer"
  }
}

resource "kubernetes_deployment" "birthdayapp" {
  metadata {
    name = "birthdayapp-deployment"
    labels = {
      app = "birthdayapp"
    }
  }

  spec {
    replicas = 3

    selector {
      match_labels = {
	app = "birthdayapp"
      }

      template {
	metadata {
	  labels = {
	    app = "birthdayapp"
	  }
	}

	spec {
	  container {
	    image = ""
	    imagePullPolicy = "IfNotPresent"
	    name = "birthdayapp"
	    ports {
	      containerPort = 8080
	      protocol = "TCP"
	    }
	  }

	  resources {
	    
	  }
	}
      }
    }
  }
}


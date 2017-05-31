# Manage the Choria discovery_proxy service
class choria_discovery_proxy::service {
  service { $choria_discovery_proxy::service_name:
    ensure   => $choria_discovery_proxy::service_ensure,
    enable   => true,
  }
}
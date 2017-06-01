# Manage the Choria discovery_proxy service
class discovery_proxy::service {
  service { $discovery_proxy::service_name:
    ensure   => $discovery_proxy::service_ensure,
    enable   => true,
  }
}

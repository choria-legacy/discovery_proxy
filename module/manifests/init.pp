class discovery_proxy (
    String $listen = "0.0.0.0",
    Integer $tls_port = 8085,
    Optional[Integer] $port = 0,
    String $puppetdb_host = "puppet",
    Integer $puppetdb_port = 8081,
    String $cert_file = "/etc/puppetlabs/puppet/ssl/certs/${trusted['certname']}.pem",
    String $key_file = "/etc/puppetlabs/puppet/ssl/private_keys/${trusted['certname']}.pem",
    String $ca_file = "/etc/puppetlabs/puppet/ssl/certs/ca.pem",
    Boolean $manage_user = false,
    Boolean $manage_group = false,
    String $user = "root",
    String $group = "root",
    String $binpath = "/usr/bin/discovery_proxy",
    String $binary_source = "puppet:///modules/discovery_proxy/discovery_proxy-0.1.0",
    String $service_name = "discovery_proxy",
    Enum["running", "stopped"] $service_ensure = "running",
    String $db_dir = "/var/lib/choria",
    Boolean $debug = false,
    Boolean $manage_service = true,
) {
    include discovery_proxy::install

    if $manage_service {
        include discovery_proxy::service
    }
}

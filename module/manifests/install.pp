# Manage the Choria discovery_proxy binary
class choria_discovery_proxy::install {
    if $choria_discovery_proxy::manage_group and $choria_discovery_proxy::group != "root" {
        group { $choria_discovery_proxy::group:
            ensure    => present,
            allowdupe => false,
        }
    }

    if $choria_discovery_proxy::manage_user and $choria_discovery_proxy::choria_discovery_proxy != "root" {
        user { $choria_discovery_proxy::user:
            ensure    => present,
            gid       => $choria_discovery_proxy::group,
            allowdupe => false,
        }
    }

    file {
        default:
            owner => $choria_discovery_proxy::user,
            group => $choria_discovery_proxy::group,
            mode  => "0755";

        $choria_discovery_proxy::binpath:
            source => $choria_discovery_proxy::binary_source;

        $choria_discovery_proxy::db_dir:
            ensure => "directory";
    }

    contain "::systemd"

    systemd::unit_file { "${choria_discovery_proxy::service_name}.service":
        content => epp("choria_discovery_proxy/systemd_service.epp"),
    }

    Class[$name] ~> Class["choria_discovery_proxy::service"]
}
# Manage the Choria discovery_proxy binary
class discovery_proxy::install {
    if $discovery_proxy::manage_group and $discovery_proxy::group != "root" {
        group { $discovery_proxy::group:
            ensure    => present,
            allowdupe => false,
        }
    }

    if $discovery_proxy::manage_user and $discovery_proxy::discovery_proxy != "root" {
        user { $discovery_proxy::user:
            ensure    => present,
            gid       => $discovery_proxy::group,
            allowdupe => false,
        }
    }

    file {
        default:
            owner => $discovery_proxy::user,
            group => $discovery_proxy::group,
            mode  => "0755";

        $discovery_proxy::binpath:
            source => $discovery_proxy::binary_source;

        $discovery_proxy::db_dir:
            ensure => "directory";
    }

    contain "::systemd"

    systemd::unit_file { "${discovery_proxy::service_name}.service":
        content => epp("discovery_proxy/systemd_service.epp"),
    }

    Class[$name] ~> Class["discovery_proxy::service"]
}

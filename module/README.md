# Choria Discovery Proxy


A basic module to install the Choria Discovery Proxy service and it's cli.

## Overview

When configuring Choria to use PuppetDB as a discovery source it's required to expose the PuppetDB query interface to all users.  This can be a source of potential secrets leak due to the vast amount of data stored in PuppetDB.

This proxy sits in front of PuppetDB and exposes a HTTPS secure REST service that Choria uses to do discovery.  This service will only return certnames thus greatly reducing the possibility of sensitive information leaking.  You now only have to allow this proxy to communicate with PuppetDB directly.

Additionally it allows named *sets* to be created that can later be referenced by name in Choria discovery.

## Usage

By default this module sets up the proxy to listen on `0.0.0.0:8085` for incoming HTTPS requests from clients using HTTPS client certificates signed by the Puppet CA.

```puppet
class{"choria_discovery_proxy": }
```

This sets it all up working with your PuppetDB on `puppet:8081`.

There are many customizations available, you can specify custom PuppetDB location and ports for example:

```puppet
class{"choria_discovery_proxy": 
    tls_port => 9292,
    puppetdb_host => "puppetdb.example.net",
}
```

See the module source for other options available.

In the above cases the server process will be started, you can install just the client using this:

```puppet
class{"choria_discovery_proxy": 
    manage_service => false
}
```

Then look at the `discovery_proxy sets --help` output to see about maintaining sets.

## MCollective

At present integration with MCollective is not yet released, but eventually to integrate with this service you'd add SRV records like:

```
_mcollective-discovery._tcp   IN  SRV 10  0 8085  puppetdb1.example.net.
```

And then enable use of the proxy by setting `plugin.choria.discovery_proxy` to `true`.

The host and port can also be set using `plugin.choria.discovery_host` and `plugin.choria.discovery_port`.

Further documentation will be written in the main choria docs about this integration.

A created set can be discovered using something like `mco package update foo -I set:bobs_machines` where `bobs_machines` were made using `discovery_proxy sets create bobs_machines`.

## Compatibility

This is a early release of the module and the proxy so for now the module embeds the compiled proxy as a binary and will only support Linux distributions using SystemD.
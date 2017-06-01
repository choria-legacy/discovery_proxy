# Choria Discovery Proxy for PuppetDB

When configuring the PuppetDB discovery for Choria it is required to open the PuppetDB query port to all clients.

This is a security problem because a lot of sensitive information lives in PuppetDB and it's nearly impossible to sanitise.

This project provides a proxy service that listens on HTTP and HTTPS for discovery requests from MCollective and performs the PQL query on its behalf

This way the PuppetDB query interface only have to be opened to the proxy and not to everyone.

Additionally it provides a way to store PQL queries and give them names.  You can thus create a PQL query that match a subset of machines and store them as `bobs_machines`, later discover will be able to refer to this set in identity filters with `-I set:bobs_machines`.

## Starting

This proxy needs certificates signed by your Puppet CA, use `mco choria request_cert pdb_proxy.example` to create them.  They have to match your hostname.

```
$ discovery_proxy server --help
usage: discovery_proxy server --ca=CA --cert=CERT --key=KEY --db=DB [<flags>]

Runs a Proxy Server

Flags:
      --help                    Show context-sensitive help (also try --help-long and --help-man).
      --version                 Show application version.
  -d, --debug                   Enable debug logging
  -l, --listen="0.0.0.0"        address to bind to for client requests
  -p, --port=0                  HTTP port to bind to for client requests
      --tlsport=8085            HTTPS port to bind to for client requests
  -H, --puppetdb-host="puppet"  PuppetDB host
  -P, --puppetdb-port=8081      PuppetDB port
      --ca=CA                   Certificate Authority file
      --cert=CERT               Public certificate file
      --key=KEY                 Public key file
      --logfile=LOGFILE         File to log to, STDOUT when not set
      --db=DB                   Path to the database file to write
```

## Data Storage

Data is stored in a [BoltDB](https://github.com/boltdb/bolt) instance in the path you specify.  BoltDB locks the store so you can only have single access to it.  To make a backup of this file you have to hit */v1/backup* with a simple GET request, this will make `your.db.bak`, you can safely back this file up.

This backup file may be smaller than the main data base that's because the backup also compacts the store.  There are some CLI tools listed in the BoltDB README that you can use to view the data in the backup.

In future other backends - probably consul - will be supported to allow this service to be made Highly Available.

## Client

A client is included to create and edit sets, here's some of it in use:

```
$ discovery_proxy sets
Found 6 set(s)

   test                   test_set3              test_set4
   test_set5              test_set6              test_set7

$ discovery_proxy sets create mt_hosts
Please enter a PQL query, you can scroll back for history and use normal shell editing short cuts:
pql> inventory { facts.country = "mt" }
Matched Nodes:

   dev1.example.net           dev10.example.net          dev11.example.net
   dev12.example.net          dev13.example.net          dev2.example.net
   dev3.example.net           dev4.example.net           dev5.example.net
   dev6.example.net           dev7.example.net           dev8.example.net
   dev9.example.net           nuc1.example.net           nuc2.example.net

Do you want to store this query (y/n)> y

$ discovery_proxy set view mt_hosts
Details for the 'mt_hosts' set

Query:

    inventory { facts.country = "mt" }


Use --discover to view matching nodes

$ discovery_proxy sets view mt_hosts --discover
Details for the 'mt_hosts' set

Query:

    inventory { facts.country = "mt" }

Matched Nodes:

   dev1.example.net           dev10.example.net          dev11.example.net
   dev12.example.net          dev13.example.net          dev2.example.net
   dev3.example.net           dev4.example.net           dev5.example.net
   dev6.example.net           dev7.example.net           dev8.example.net
   dev9.example.net           nuc1.example.net           nuc2.example.net

$ discovery_proxy sets rm mt_hosts
Details for the 'mt_hosts' set

Query:

    inventory { facts.country = "mt" }


Are you sure you wish to delete this node set? (y/n)> y
$
```

To use this you should have in mcollective *client.cfg*

```ini
plugin.choria.discovery_host = proxy.example.net
plugin.choria.discovery_port = 9293
```

Access to the HTTPS port is secured with your normal choria SSL certificate.

## API Format

The API is defined in a Swagger API found in `schema.yaml` and can be previewed on the [Swagger UI](http://petstore.swagger.io/?url=https://raw.githubusercontent.com/choria-io/discovery_proxy/master/schema.yaml).


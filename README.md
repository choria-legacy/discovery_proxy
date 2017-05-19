# Choria Discovery Proxy for PuppetDB

When configuring the PuppetDB discovery for Choria it is required to open the PuppetDB query port to all clients.

This is a security problem because a lot of sensitive information lives in PuppetDB and it's nearly impossible to sanitise.

This project provides a proxy service that listens on HTTP and HTTPS for discovery requests from MCollective and performs the PQL query on its behalf

This way the PuppetDB query interface only have to be opened to the proxy and not to everyone.

Additionally it provides a way to store PQL queries and give them names.  You can thus create a PQL query that match a subset of machines and store them as `bobs_machines`, later discover will be able to refer to this set in identity filters with `-I set: bobs_machines`.

## Starting

This proxy needs certificates sign by your Puppet CA, use `mco choria request_cert pdb_proxy.example` to create them.

```
$ pdbproxy --help
usage: pdbproxy --ca=CA --cert=CERT --key=KEY --db=DB [<flags>]

Service performing Choria discovery requests securely against PuppetDB

Required flags:
  --ca=CA      Certificate Authority file
  --cert=CERT  Public certificate file
  --key=KEY    Public key file
  --db=DB      Path to the database file to write

Optional flags:
      --help                    Show context-sensitive help (also try --help-long and --help-man).
      --version                 Show application version.
  -d, --debug                   enable debug logging
  -l, --listen="0.0.0.0"        address to bind to for client requests
  -p, --port=0                  port to bind to for client requests
      --tlsport=8081            port to bind to for client HTTPS requests
  -H, --puppetdb-host="puppet"  PuppetDB host
  -P, --puppetdb-port=8081      PuppetDB port
      --logfile=LOGFILE         File to log to, STDOUT when not set
```

## API Format

The API is defined in a Swagger API found in `schema.yaml` and can be previewed on the [Swagger UI](http://petstore.swagger.io/?url=https://raw.githubusercontent.com/choria-io/pdbproxy/master/schema.yaml).

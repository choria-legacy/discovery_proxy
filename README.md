# Choria Discovery Proxy for PuppetDB

When configuring the PuppetDB discovery for Choria it is required to open the PuppetDB query port to all clients.

This is a security problem because a lot of sensitive information lives in PuppetDB and it's nearly impossible to sanatise.

This project provides a proxy service that listens on HTTP for discovery requests from MCollective and performs the PQL query on its behalf

This way the PuppetDB query interface only have to be opened to the proxy and not to everyone.

## Starting

This proxy needs certificates sign by your Puppet CA, use `mco choria request_cert pdb_proxy.example` to create them.

```
$ pdbproxy --help
usage: pdbproxy --ca=CA --cert=CERT --key=KEY [<flags>]

Service performing Choria discovery requests securely against PuppetDB

Required flags:
  --ca=CA      Certificate Authority file
  --cert=CERT  Public certificate file
  --key=KEY    Public key file

Optional flags:
      --help                    Show context-sensitive help (also try --help-long and --help-man).
      --version                 Show application version.
  -d, --debug                   enable debug logging
  -l, --listen="0.0.0.0"        address to bind to for client requests
  -p, --port=8080               port to bind to for client requests
  -H, --puppetdb-host="puppet"  PuppetDB host
  -P, --puppetdb-port=8081      PuppetDB port
      --logfile=LOGFILE         File to log to, STDOUT when not set
```

## API Format

The API is defined in a Swagger API found in `schema.yaml` and can be previewed on the [Swagger UI](petstore.swagger.io/?url=https://raw.githubusercontent.com/choria-io/pdbproxy/master/schema.yaml).

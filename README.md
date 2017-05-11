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

## Request Format

MCollective will submit to the service a JSON document with the following structure, all fields are optional, a empty
hash will discover all nodes:

```json
{
	"facts": [
		{ "fact": "country", "operator": "==", "value": "mt" }
	],
	"classes": ["docker", "/registrator/"],
	"agents": ["rpcutil", "weather"],
	"identities": ["example.com", "/another/"],
	"collective": "mt_collective"
}
```

## TODO

 - [ ] Listen for clients using HTTPS
 - [ ] Listen on the middleware to facilitate easy scaling and avoid opening listen ports
 - [ ] Export stats using `expvar`

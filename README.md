# Choria Discovery Proxy for PuppetDB

When configuring the PuppetDB discovery for Choria it is required to open the PuppetDB query port to all clients.

This is a security problem because a lot of sensitive information lives in PuppetDB and it's nearly impossible to sanitise.

This project provides a proxy service that listens on HTTP and HTTPS for discovery requests from MCollective and performs the PQL query on its behalf

This way the PuppetDB query interface only have to be opened to the proxy and not to everyone.

Additionally it provides a way to store PQL queries and give them names.  You can thus create a PQL query that match a subset of machines and store them as `bobs_machines`, later discover will be able to refer to this set in identity filters with `-I set:bobs_machines`.

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

## Data Storage

Data is stored in a [BoltDB](https://github.com/boltdb/bolt) instance in the path you specify.  BoltDB locks the store so you can only have single access to it.  To make a backup of this file you have to hit */v1/backup* with a simple GET request, this will make `your.db.bak`, you can safely back this file up.

This backup file may be smaller than the main data base that's because the backup also compacts the store.  There are some CLI tools listed in the BoltDB README that you can use to view the data in the backup.

## API Format

The API is defined in a Swagger API found in `schema.yaml` and can be previewed on the [Swagger UI](http://petstore.swagger.io/?url=https://raw.githubusercontent.com/choria-io/pdbproxy/master/schema.yaml).

## Sample Requests

### Discovery
The discovery requests look a lot like those that MCollective `Util::empty_filter` produce, full details are in the Swagger schema but here some examples:

Filter that do not apply to your request can be left out, you can combine them along the same rules as any MCollective discovery would work - generally they get boolean `AND`d together.

Discovery filters are made using *GET* to */v1/discover*, the JSON goes in the body.

#### Facts

<details>
<summary>-F country=mt</summary>

```json
{
	"facts": [
		{
			"fact": "country",
			"operator": "==",
			"value": "mt"
		}
	]
}
```
</details>

<details>
<summary>-F lsbdistdescription=/centos/</summary>

```json
{
	"facts": [
		{
			"fact": "lsbdistdescription",
			"operator": "=~",
			"value": "/centos/"
		}
	]
}
```
</details>

#### Classes

<details>
<summary>-C nats -C /puppetdb/</summary>

```json
{
	"classes": [
        "nats",
        "/puppetdb/"
	]
}
```
</details>

#### Agents

<details>
<summary>-A apache -A /weather/</summary>

```json
{
	"agents": [
        "apache",
        "/weather/"
	]
}
```
</details>

#### Identities

<details>
<summary>-I /web/</summary>

```json
{
	"identities": [
        "/web/"
	]
}
```
</details>

<details>
<summary>-I pql:nodes {}</summary>

```json
{
	"identities": [
        "pql:nodes {}"
	]
}
```
</details>

<details>
<summary>-I set:bobs_machines</summary>

```json
{
	"identities": [
        "set:bobs_machines"
	]
}
```
</details>

#### PQL

This is not really used by MCollective but it's handy.

<details>
<summary>arbitrary PQL</summary>

```json
{
	"queries": [
        "nodes {}"
	]
}
```
</details>

### Set maintenance

Set names must match *^[a-zA-Z0-9_\\-\\.]+$*

#### List Sets

GET request to */v1/sets*

#### Create a Set

<details>
<summary>POST /v1/set</summary>

```json
{
	"set": "test_set",
	"query": "nodes { (certname in inventory[certname] { facts.mcollective.server.collectives.match(\"\\d+\") = \"mt_collective\" }) }",
	"nodes": []
}
```
</details>

#### Retrieve a Set

Do a GET against */v1/set/test_set* you can also request the node list to be included like */v1/set/test_set?discover=true*

#### Update a Set

<details>
<summary>PUT /v1/set/test_set</summary>

```json
{
	"set": "test_set",
	"query": "nodes { (certname in inventory[certname] { facts.mcollective.server.collectives.match(\"\\d+\") = \"mt_collective\" }) }",
	"nodes": []
}
```
</details>

#### Delete a Set

Do a DELETE against */v1/set/test_set*
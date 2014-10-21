# gopherneo
=========

A lean, pragmatic and battle-tested Neo4j driver for Go. GopherNeo ("Go4Neo") conforms strictly to the latest [Neo4j HTTP API](http://docs.neo4j.org/chunked/stable/rest-api.html), supporting only the latest Neo4j 2.x features. 

For comprehensive, up to date documentation and code examples, checkout the wiki. 

# Examples
==
(todo)

# Completed Features
==
* connect to neo
* get node by label, id
* get node by label, property, value
* execute a cypher query
* create node

# Feature Roadmap
==
## HIGH
* create unique node with label and properties
* set node/relationship properties
* execute cypher query, return nodes/relationships
* execute cypher query, return list of properties
* remove node property
* delete node

## MEDIUM (convenience)
* link node to another node
* unlink node from another node
* list node's linked nodes
   - specifcy order, pagination
* list node's relationships
   - specify order, pagination

### Troubleshooting via Curl Examples

````
curl -X POST \
  -H "Accept: application/json; charset=UTF-8" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "CREATE (t:Thing { name: '915738050', age: 45 }) RETURN t",
    "params": {}
  }' \
  http://localhost:7474/db/data/cypher

curl -X POST \
  -H "Accept: application/json; charset=UTF-8" \
  -H "Content-Type: application/json" \
  -d '{
    "prop1": "val1"
  }' \
  http://localhost:7474/db/data/node


 ````



# gopherneo
Go4Neo is a Neo4j driver for the Go programming language. It's focus is on simple, practical and realistic features that 90% of users working with Neo4j would expect.

# Version Info
- v0.2, 01/22/15:  Supports authentication tokens for Neo4j 2.2.x and above
- v0.1, 12/23/14:  Alpha version. (Though it will soon be used for a real app and upgraded to "Beta" sometime in Jan)

Compatible with Neo4j 2.x and up.

# Examples
- TODO, see test for now

# Supported Features
* execute a cypher query
* get node with label, property
* get nodes with property, paginated
* create node with label, properties
* update a node
* delete node
* link nodes
* unlink nodes
* list linked nodes

# Planned Features
## High Priority
* return descriptive error if can't connect to server or creds invalid
* ability to return multiple entire nodes
* set relationship properties

## Medium Priority
* solution for passing order by option for listing nodes

## Low Priority
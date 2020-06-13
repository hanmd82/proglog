### Chapter 6 Notes

- Service discovery is the process of figuring out how to connect to a service.
- A service discovery mechanism keeps an up-to-date registry of services, and their locations and health.
- Downstream services then query this registry to discover the location of upstream services and connect to them.
- Load balancers add cost, increase latency, introduce single points of failure, and have to be updated as services scale up and down.
  - This would increase operational burden, infrastructure costs, and latency.
- Two service-discovery problems to solve
  - How will the servers in the cluster discover each other?
  - How will the clients discover the servers?

- Requirements of a service discovery tool
  - Manage a registry of servers containing information such as IP addresses, ports
  - Help services find other services using the registry
  - Perform health checks on service instances and remove them if they are not responding well
  - De-register services when they go offline

---

- Serf is a Golang library that provides decentralised cluster membership, failure detection and orchestration
  - can be used to embed service discovery into distributed services
  - so that we dont have to implement service discovery ourselves, and users dont need to run an addition service discovery cluster
- Serf maintains cluster membership using a gossip protocol to communicate between nodes - Serf does not use a central registry

Service Discovey with Serf
1. Create a Serf node on each server
2. Configure each Serf node with an address to listen on, and accept connections from other Serf nodes
3. Configure each Serf node with addresses of other Serf nodes, and join their cluster
4. Handle Serf's cluster discovery events, such as when a node joins or fails in the cluster

```bash
# build and test
cd ./internal/discovery
go test -c
./discovery.test
```
---

- Replication can make a service more resilient to failures. For example, if a node’s disk fails and its data cannot be recovered, there is a copy of the data on another disk
- Enable the servers to replicate each other when they discover each other
- Discovery is important because the discovery events trigger other processes, like replication and consensus
- When a server joins the cluster, the replicator component will connect to the server and run a loop that consumes from the discovered server and produces to the local server.
  - pull-based replication: consumer periodically polls the data source to check if it has new data to consume
  - push-based replication: the data source pushes the data to its replicas
- Lazy initialization gives structs a useful zero value, which reduces the API’s size and complexity while maintaining the same functionality
  - otherwise, may need to export constructor functions, or getter/setter functions on fields in the struct
- `chan struct{}` is typically used as a signalling channel
  - `leave chan struct{}` is used to signal that the destination server has left the cluster and current server should stop replicating from it
- use the `dynaport` library to allocate two ports: one for gRPC log connections and one for Serf service discovery connections

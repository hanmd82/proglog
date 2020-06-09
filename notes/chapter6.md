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


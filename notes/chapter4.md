### Chapter 4 Notes

- gRPC uses types to define service methods, request and response bodies. Enables type-checking of requests, responses, models, and serialization.
- Use different kinds of load balancing with gRPC based on needs, including thick client-side load balancing, proxy load balancing, look-aside balancing, or service mesh - https://grpc.io/blog/grpc-load-balancing.
- A gRPC service is essentially a group of related RPC endpoints. Creating a gRPC service involves defining it in protobuf and then compiling protocol buffers into code comprising the client and server stubs.


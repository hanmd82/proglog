### Chapter 4 Notes

- gRPC uses types to define service methods, request and response bodies. Enables type-checking of requests, responses, models, and serialization.
- Use different kinds of load balancing with gRPC based on needs, including thick client-side load balancing, proxy load balancing, look-aside balancing, or service mesh - https://grpc.io/blog/grpc-load-balancing.
- A gRPC service is essentially a group of related RPC endpoints. Creating a gRPC service involves defining it in protobuf and then compiling protocol buffers into code comprising the client and server stubs.
- Go’s gRPC implementation has a `status` package which can be used to build errors with status codes, and include other data.
- To create an error with a status code, create the `error` with the `Error` function from the `status` package, and pass the relevant code from the `codes` package that matches the type of error.
- Add a custom error type `ErrOffsetOutOfRange` that the server will send back to the client when the client tries to consume an offset that’s outside of the log. `ErrOffsetOutOfRange` includes a localized message, a status code, and an error message.

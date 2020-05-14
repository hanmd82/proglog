### Chapter 5 Notes

- Security in distributed services can be broken down into three steps:
  1. Encrypt data in-flight to protect against man-in-the-middle attacks
  2. Authenticate to identify clients
  3. Authorize to determine the permissions of the identified client

---
**Encrypt In-Flight Data**
- The most widely used technology for preventing MITM attacks and encrypting data in-flight is Transport Layer Security (TLS), the successor to SSL
- Build TLS support into the Log service to encrypt data in-flight and authenticate the server. During TLS handshake, the client and server:
  1. Specify which version of TLS they’ll use
  2. Decide which cipher suites (the set of encryption algorithms) they’ll use
  3. Authenticate the identity of the server via the server’s private key and the certificate authority’s digital signature
  4. Generate session keys for symmetric encryption after the handshake is complete

---
**Authenticate to Identify Clients**
- Authentication is the process of identifying who the client is (TLS has already handled authenticating the server)
- Often, client authentication is left to the application to work out, usually by some combination of username-password credentials and tokens
- TLS mutual authentication (or two-way authentication) in which both the server and the client validate the other’s communication, is more commonly used in machine-to-machine communications and distributed systems. Both the server and the client use a certificate to prove their identity

---
**Authorize to Determine the Permissions of Clients**
- Differentiating between authentication and authorization is necessary for resources with shared access and varying levels of ownership
- In this Log service, build access control list-based authorization to control whether a client is allowed to read from or write to (or both) the log.

---

- CloudFlare has written a toolkit called CFSSL that can be used for signing, verifying, and bundling TLS certificates so that it can act as its own CA for internal services:
  - `cfssl` to sign, verify, and bundle TLS certificates and output the results as JSON
  - `cfssljson` to take that JSON output and split them into separate key, certificate, CSR, and bundle files

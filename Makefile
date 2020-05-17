CONFIG_PATH=${HOME}/.proglog/

.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl gencert \
		-initca test/ca-csr.json | cfssljson -bare ca

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=server \
		test/server-csr.json | cfssljson -bare server
	mv *.pem *.csr ${CONFIG_PATH}

.PHONY: test
test:
	go test -race ./...

.PHONY: compile
compile:
	protoc api/v1/*.proto \
	--gogo_out=Mgogoproto/gogo.proto=github.com/gogo/protobuf/proto,plugins=grpc:. \
	--proto_path=$$(go list -f '{{ .Dir }}' -m github.com/gogo/protobuf) \
	--proto_path=.

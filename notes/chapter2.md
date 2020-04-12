### Chapter 2 Notes

- Go has two runtimes to compile protobuf into Go code.
- `gogoprotobuf` extends the original Go runtime and adds code generation features, e.g. embed fields, cast generated types, and generate non-pointer fields. Includes APIs for building plugins.
- Download and configure Protobuf 3.11.4
    ```bash
    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.11.4/protoc-3.11.4-osx-x86_64.zip
    unzip protoc-3.11.4-osx-x86_64.zip -d /usr/local/protobuf
    echo 'export PATH="/usr/local/protobuf/bin:$PATH"' >> ~/.bashrc
    source ~/.bashrc
    ```

- Download and install `protoc-gen-gogo` plugin
    ```bash
    cd $GOPATH
    git clone git@github.com:gogo/protobuf.git
    git checkout v1.3.1
    cd protobuf
    make install

    # check that protoc-gen-gogo binary is now in your $GOPATH/bin
    ls $GOPATH/bin | grep protoc-gen-gogo
    echo 'export PATH="$PATH:$GOPATH/bin"' >> ~/.bashrc
    source ~/.bashrc
    ```

- Build protobuf
    ```bash
    protoc api/v1/*.proto \
    --gogo_out=Mgogoproto/gogo.proto=github.com/gogo/protobuf/proto:. \
    --proto_path=$(go list -f '{{ .Dir }}' -m github.com/gogo/protobuf) \
    --proto_path=.
    ```

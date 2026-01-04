.PHONY: proto-gen
proto-gen:
	rm -rf ./gen
	mkdir -p ./gen
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	find proto -name "*.proto" -print0 | xargs -0 protoc \
		--go_out=gen \
        --go_opt=paths=source_relative \
        --go-grpc_out=gen \
        --go-grpc_opt=paths=source_relative
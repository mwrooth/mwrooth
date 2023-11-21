SHELL := /bin/bash

.PHONY: gen*, postgres_init

gen:
	$(MAKE) gen-grpc-stubs

gen-grpc-stubs:
	@for f in $$(find protobuf -name "*.proto"); do \
		file_base=$${f%.proto}; \
		if [ ! -f $$file_base.pb.go -o $$f -nt $$file_base.pb.go ]; then \
			$$(protoc \
				-I $$HOME/.local/include \
				-I ./protobuf \
				--go_out=./protobuf \
				--go_opt=paths=source_relative \
				--go-grpc_out=require_unimplemented_servers=false:./protobuf \
				--go-grpc_opt=paths=source_relative $$f);\
		fi \
	done

clean:
	find ./protobuf -type f -name '*.mock.go' | xargs rm -f
	find ./protobuf -type f -name '*.pb.go' | xargs rm -f

postgres_init:
	docker run -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=root -e -d --name postgres postgres
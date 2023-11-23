SHELL := /bin/bash
CONTAINER_NAME := postgres
CONTAINER_PORT := 5433

.PHONY: gen*, postgres*

gen:
	$(MAKE) gen-grpc-stubs

gen-grpc-stubs:
	@for f in $$(find pkg -name "*.proto"); do \
		file_base=$${f%.proto}; \
		if [ ! -f $$file_base.pb.go -o $$f -nt $$file_base.pb.go ]; then \
			$$(protoc \
				-I $$HOME/.local/include \
				-I ./pkg \
				--go_out=./pkg \
				--go_opt=paths=source_relative \
				--go-grpc_out=require_unimplemented_servers=false:./pkg \
				--go-grpc_opt=paths=source_relative $$f);\
		fi \
	done

clean:
	find ./pkg -type f -name '*.mock.go' | xargs rm -f
	find ./pkg -type f -name '*.pb.go' | xargs rm -f

postgres_init:
	docker run -p $(CONTAINER_PORT):5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=root -d --rm --name $(CONTAINER_NAME) postgres

postgres_clean:
	docker container rm -f $(CONTAINER_NAME)

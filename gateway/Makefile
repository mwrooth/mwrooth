SHELL := /bin/bash

.PHONY: gen*

IMAGE_NAME := "orchasystems/openapi-generator:latest"

gen:
	$(MAKE) gen-openapi

gen-openapi:
	docker run --rm -v $$(pwd):/local $(IMAGE_NAME) -- -r /local -f ./openapi/specs/openapi.yaml -o ./openapi

clean:
	find ./ -type f -name '*.mock.go' | xargs rm -f
	find ./ -type f -name '*.pb.go' | xargs rm -f

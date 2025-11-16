.PHONY: run-restapi run-grpcapi clean api_generate buf_generate go_generate generate clean preview_open_api

run-restapi:
	cd cmd && go run . restapi -s std-out -z file-writer

run-grpcapi:
	cd cmd && go run . grpcapi -s std-out -z file-writer

api_generate: api/openapi/api.yaml
	npx openapi-format api/openapi/api.yaml -s api/openapi/openapi-sort.json -f api/openapi/openapi-filter.json -o api/openapi/api.yaml
	mkdir -p internal/gen/restapigen
	go tool oapi-codegen --package restapigen -generate types $< > internal/gen/restapigen/api-types.gen.go
	go tool oapi-codegen --package restapigen -generate chi,spec $< > internal/gen/restapigen/api-server.gen.go

buf_generate:
	buf generate

go_generate:
	go generate ./...

generate: api_generate buf_generate go_generate

clean:
	rm -rf dist/* generated build vendor
	find . -name "*.mock.gen.go" -type f -delete
	find . -name "*.out" -type f -delete
	find . -name "wire_gen.go" -type f -delete
	find . -name "*.mock.gen.go" -type f -delete

preview_open_api:
	redocly preview-docs api/openapi/api.yaml
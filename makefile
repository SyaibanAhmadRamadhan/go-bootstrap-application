.PHONY: run clean 

run:
	cd cmd && go run . restapi -s std-out -z file-writer

generate_api: api/openapi/api.yaml
	npx openapi-format api/openapi/api.yaml -s api/openapi/openapi-sort.json -f api/openapi/openapi-filter.json -o api/openapi/api.yaml
	mkdir -p gen/restapigen
	go tool oapi-codegen --package restapigen -generate types $< > gen/restapigen/api-types.gen.go
	go tool oapi-codegen --package restapigen -generate chi,spec $< > gen/restapigen/api-server.gen.go

clean:
	rm -rf dist/* generated build vendor
	find . -name "*.mock.gen.go" -type f -delete
	find . -name "*.out" -type f -delete
	find . -name "wire_gen.go" -type f -delete
	find . -name "*.mock.gen.go" -type f -delete

preview_open_api:
	redocly preview-docs api/openapi/api.yaml

generate_wire:
	go tool wire ./...
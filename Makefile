.PHONY: dbuild proto

# TODO: dockerize buf+protoc for reproductive builds
# (for this demo we are using local buf and plugins)
proto:
	@cd api/proto && \
		buf lint && \
		buf generate


dbuild:
	docker build -t space-trouble-api:dev .


# TODO: dockerize db migration as well for version consistency
migration:
	@cd ./sql && \
		migrate create -ext sql -dir ./ -seq $(name)

LOCAL_BIN=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

get-deps:
	go get -u encoding/json
	go get -u github.com/pkg/errors
	go get -u google.golang.org/grpc/peer
	go get -u github.com/jackc/pgx/v4
	go get -u github.com/jackc/pgx/v4/pgxpool
	go get -u github.com/georgysavva/scany/pgxscan
	go get -u go.uber.org/zap
	go get -u github.com/sony/gobreaker
	go get -u github.com/golang-jwt/jwt/v5


	go get -u go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc

	go get -u github.com/IBM/sarama

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi

generate-api:
	mkdir -p pkg/proto/events/v1
	protoc --proto_path proto/events/v1 \
	--go_out=pkg/proto/events/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	proto/events/v1/events.proto

generate-descriptor:
	protoc -I ../ChatServer/api/chat/v1/ \
       -I ../UserServer/api/user/v1/   \
       -I ../AuthService/api/auth/v1/  \
			 -I vendor.protogen \
       --include_imports \
       --include_source_info \
       --descriptor_set_out=./infra/envoy/messanger_descriptor.pb \
       ../ChatServer/api/chat/v1/chat.proto  \
       ../UserServer/api/user/v1/user.proto  \
       ../AuthService/api/auth/v1/auth.proto


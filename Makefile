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

	go get -u go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc

	go get -u github.com/IBM/sarama

generate-api:
	mkdir -p pkg/proto/events/v1
	protoc --proto_path proto/events/v1 \
	--go_out=pkg/proto/events/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	proto/events/v1/events.proto

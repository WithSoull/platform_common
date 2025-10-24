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

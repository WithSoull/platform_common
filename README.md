# Platform Library

## Overview

This repository contains a comprehensive platform library implemented in Go, designed to provide foundational infrastructure components and utilities for building distributed and microservices-based applications. It includes modules for messaging (Kafka), metrics and telemetry (OpenTelemetry), logging, circuit breaking, rate limiting, tracing, context management, token handling, and error handling, among others.

***

## Features

- **Kafka Client and Middleware:** Producer and consumer abstractions with middleware integrations.
- **Circuit Breaker:** Fault tolerance to handle transient failures gracefully.
- **Rate Limiter:** Control request rates at different layers.
- **OpenTelemetry Integration:** Telemetry collection setup with Dockerized Prometheus and Grafana.
- **Logging:** Configurable structured logging support.
- **Metrics:** Metrics collection and configuration utilities.
- **Context Management:** Helper packages for context propagation with claims, IPs, trace IDs, and transactions.
- **Tracing:** gRPC interceptors and tracer utilities.
- **Token Management:** JWT based token generation and handling.
- **Error Handling:** Sys codes and validation utilities.
- **Protobuf Events:** Well-defined protobuf schemas for events and messages.

***

## Project Structure

```
.
├── bin/                           # Custom binaries (e.g., protoc plugins)
├── infra/                         # Infrastructure setups and configs
│   ├── kafka/                     # Kafka Docker and Makefile setups
│   └── otel/                      # OpenTelemetry stack (Collector, Grafana, Prometheus)
├── pkg/                           # Main platform library packages
│   ├── circuitbreaker/            # Circuit breaker implementation
│   ├── client/db/                 # Database client utilities
│   ├── closer/                    # Resource cleanup helper
│   ├── contextx/                  # Context-related utilities and abstractions
│   ├── kafka/                     # Kafka client and middleware
│   ├── logger/                    # Logger configuration and implementation
│   ├── metric/                    # Metrics utilities
│   ├── middleware/                # Middleware implementations for various concerns
│   ├── proto/                     # Internal protobuf events definitions
│   ├── ratelimiter/               # Rate limiting middleware
│   ├── sys/                       # System codes, error handling, validation
│   ├── tokens/                    # Token generator and JWT handling
│   └── tracing/                   # Tracing and gRPC interceptors
├── proto/                         # Protobuf definitions accessible externally
├── Makefile                       # Project build and helper commands
├── go.mod                         # Go modules file
└── go.sum                         # Go modules checksums
```

***

## Getting Started

### Prerequisites

- Go 1.XX or higher installed.
- Docker and Docker Compose for running Kafka and OpenTelemetry stack.
- Protobuf compiler (`protoc`) with Go plugin available.

### Installation

```bash
go get https://github.com/WithSoull/platform_common
```

### Running Infrastructure

```bash
make rebuild
```

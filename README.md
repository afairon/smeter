# Smeter

Smeter is a small system that records metrics and provides useful interface to monitor time series based data.

## Building Smeter

### Prerequisites

- The latest stable version of Go. Earlier releases may work, but we recommend always using the latest stable version.
At the time of writing this is Go 1.13.
- Protocol Buffers
- gRPC (needed to generate cpp sources and headers)
- Git

### Building

Build server

```
# This will generate all protobuf files and build a binary
$ make
```

Building for ARM

```
# Target ARM architecture
$ GOOS=linux GOARCH=arm make
```

Generate cpp files

```
# Generates cpp sources and headers
# Files are being generated in internal/proto/cpp
$ make proto-cpp
```

## PostgreSQL

Smeter uses PostgreSQL to store data. You need to create a database before you can use Smeter.

```
# Login as postgres
$ su - postgres

# Create database and tables
$ createdb smeter
$ psql smeter -c '\i sql/schema.pgsql'
```

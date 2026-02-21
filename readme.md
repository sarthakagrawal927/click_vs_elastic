# ClickHouse vs Elasticsearch Ingestion Benchmark

A Go benchmark project comparing bulk ingestion throughput between ClickHouse and Elasticsearch.

## What Is Benchmarked

`main.go` runs three ingestion paths:

1. ClickHouse bulk insert (`pkg/clickhouse`)
2. Elasticsearch bulk insert with official bulk indexer (`pkg/es/index.go`)
3. Elasticsearch bulk insert via raw HTTP `_bulk` (`pkg/es/v2.go`)

Workload constants (`pkg/constants/index.go`):

- `BATCH_SIZE = 5000`
- `NUM_BATCHES = 2000`

Target scale is roughly 10M records.

## Prerequisites

- Go 1.21+
- Docker
- ClickHouse running at `127.0.0.1:9000`
- Elasticsearch running at `127.0.0.1:9210` (as used by code)

## Local Setup

### 1) Start ClickHouse

```bash
docker run -d --name clickhouse \
  -p 9000:9000 -p 8123:8123 \
  clickhouse/clickhouse-server:23.8
```

Create DB + table:

```sql
CREATE DATABASE IF NOT EXISTS fp;

CREATE TABLE fp.index_test (
  id String DEFAULT generateUUIDv4(),
  name String,
  age UInt8,
  height Float32,
  weight Float32,
  is_active UInt8
) ENGINE = MergeTree()
ORDER BY id;
```

### 2) Start Elasticsearch

```bash
docker run --rm --name elastic \
  -p 9210:9200 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  docker.elastic.co/elasticsearch/elasticsearch:8.8.0
```

Note: `compose.yaml` is available for multi-node experiments, but code defaults to `127.0.0.1:9210`.

### 3) Run benchmark

```bash
go mod download
go run .
```

## Useful Validation Commands

```bash
curl 'http://127.0.0.1:9210/test_index/_count?pretty=true&q=*%3A*'
docker stats
```

## Current Observations

- ClickHouse path batches directly via prepared inserts
- ES bulk indexer flush behavior differs from ClickHouse batching model
- At high parallelism, Elasticsearch may require tuning (TCP buffers, file descriptors, connection limits)

## Next Improvements

- Parameterize batch size/worker count via flags
- Standardize flush strategy for fairer A/B comparison
- Persist benchmark runs and produce repeatable reports

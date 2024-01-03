# Clickhouse vs ElasticSearch

## Setup

```bash
go mod download
go run *.go
```

## Clickhouse

[QuickStart](https://clickhouse.com/docs/en/getting-started/quick-start)

Creating Tables:

```sql
create database if not exists fp;

create table fp.index_test (
    id String default generateUUIDv4() PRIMARY KEY,
    name String,
    age UInt8,
    height Float32,
    weight Float32,
    is_active UInt8,
) engine = MergeTree()
```

## ES

```bash
docker run --rm --name elasticsearch_container -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" -e "xpack.security.enabled=false" elasticsearch:8.8.0
```

Cluster Mode:
```bash
docker-compose up
```

[Track Writes node wise](http://localhost:9200/_cat/thread_pool/write?v)

Stats:
```bash
docker stats node1 node2 node3
```

Rest Endpoints:

To Count

```bash
curl --location 'http://localhost:9200/test_index/_count?pretty=true&q=*%3A*' \
--header 'Cookie: lang=en-US'
```

## Right now:

Clickhouse client sends the batch request every time I want to insert a new batch, whereas ES client sends the batch request only when either the batch size is reached or flush time is reached. Because of which it is hard to predict.

Result:
For inserting 10mil docs.
```go
const BATCH_SIZE = 5000
const NUM_BATCHES = 2000
```

Output:
```
[CLICKHOUSE] bulkInsertMany took 6.715145208s
Successfully indexed [9,933,450] documents in 56.973s (174,353 docs/sec)
[ES] bulkInsertMany took 56.960981709s
```

ES Peaks at 200k docs/sec, whereas Clickhouse peaks at 1.7mil docs/sec.


## Considerations
- ES is running on a single node. Will need to study more on Clickhouse's architecture.
- Clickhouse table is MergeTree engine mode, log based engine mode might be better for this use case. Need to look into other engines as well.
- ES is not sending batch request every time unlike clickhouse. Thus net doc indexes is little less than 10m.


2000 requests in parallel is too much for ES to handle. For 5000 * 20 packets it takes:

```
[CLICKHOUSE] bulkInsertMany took 179.392625ms
[ES W/O Library] bulkInsertMany took 637.056208ms
```

Going higher was causing issues in ES due to TCP Buffer size & connection limit. It ended up making ES unresponsive for the main port. [More Info](https://chat.openai.com/share/4fa92809-5cfa-4b00-81f8-be157ab9fe2b)
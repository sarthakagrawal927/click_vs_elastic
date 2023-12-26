# Clickhouse vs ElasticSearch

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

Rest Endpoints:

To Count

```bash
curl --location 'http://localhost:9200/test_index/_count?pretty=true&q=*%3A*' \
--header 'Cookie: lang=en-US'
```
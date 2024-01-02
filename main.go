package main

import (
	"fp/chvses/pkg/clickhouse"
	elastic "fp/chvses/pkg/es"
)

func main() {
	clickhouse.TestClickHouseIngestion()
	elastic.TestESIngestion()
	elastic.BulkInsertWithoutLibrary()
}

package main

import (
	elastic "fp/chvses/pkg/es"
)

func main() {
	// clickhouse.TestClickHouseIngestion()
	// elastic.TestESIngestion()
	elastic.BulkInsertWithoutLibrary()
}

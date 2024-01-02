package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"fp/chvses/pkg/constants"
	"fp/chvses/pkg/types"
	"fp/chvses/pkg/utils"
	"log"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

const ES_IDX_NAME = "test_index"

func bulkIndexHandler() error {
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9210"},
	})
	if err != nil {
		return err
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      ES_IDX_NAME, // The default index name
		Client:     esClient,    // The Elasticsearch client
		FlushBytes: 1024 * 750,
		// FlushInterval: time.Millisecond * 1,
	})
	if err != nil {
		return err
	}
	esClient.Indices.Delete([]string{ES_IDX_NAME})
	esClient.Indices.Create(ES_IDX_NAME)

	return bulkInsertMany(&bi)
}

func bulkInsertMany(bi *esutil.BulkIndexer) error {
	var ctx = context.Background()
	defer (*bi).Close(ctx)
	defer utils.Timer("[ES] bulkInsertMany")()
	start := time.Now().UTC()
	data, err := json.Marshal(&types.Person{
		ID:       1,
		Age:      20,
		Name:     "John Doe",
		Height:   1.8,
		Weight:   70,
		IsActive: 1,
	})
	if err != nil {
		return err
	}

	var pg sync.WaitGroup
	pg.Add(constants.NUM_BATCHES)
	for i := 0; i < constants.NUM_BATCHES; i++ {
		go func() {
			for i := 0; i < constants.BATCH_SIZE; i++ {
				err := (*bi).Add(
					ctx,
					esutil.BulkIndexerItem{
						Index:  ES_IDX_NAME,
						Action: "index",
						Body:   bytes.NewReader(data),
					},
				)
				if err != nil {
					fmt.Println(err)
				}
			}
			pg.Done()
		}()
	}
	pg.Wait()
	biStats := (*bi).Stats()
	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Successfully indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}
	return nil
}

func TestESIngestion() {
	err := bulkIndexHandler()
	if err != nil {
		log.Fatal(err)
	}
}

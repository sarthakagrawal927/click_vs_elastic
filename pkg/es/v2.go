package elastic

import (
	"bytes"
	json "encoding/json"
	"fmt"
	"fp/chvses/pkg/constants"
	"fp/chvses/pkg/types"
	"fp/chvses/pkg/utils"
	"sync"

	"net/http"
)

type BulkInsertMetaData struct {
	Index BulkInsertIndex `json:"index"`
}

type BulkInsertIndex struct {
	Index string `json:"_index"`
}

func BulkInsertWithoutLibrary() {
	dataMeta := BulkInsertMetaData{
		Index: BulkInsertIndex{
			Index: "test_index",
		},
	}
	data := types.Person{
		ID:       1,
		Age:      20,
		Name:     "John Doe",
		Height:   1.8,
		Weight:   70,
		IsActive: 1,
	}

	TojsBulInsertData, _ := json.Marshal(data)
	TojsBulkInsertMetaData, _ := json.Marshal(dataMeta)

	BulkMetaData := bytes.NewBuffer(append(TojsBulkInsertMetaData, []byte("\n")...))
	BulkData := bytes.NewBuffer(append(TojsBulInsertData, []byte("\n")...))

	defer utils.Timer("[ES W/O Library] bulkInsertMany")()

	var pg sync.WaitGroup
	pg.Add(constants.NUM_BATCHES)
	for i := 0; i < constants.NUM_BATCHES; i++ {
		go func() {
			combinedBulk := bytes.NewBuffer(append(BulkMetaData.Bytes(), BulkData.Bytes()...))

			for i := 1; i < constants.BATCH_SIZE-1; i++ {
				combinedBulk.Write(append(BulkMetaData.Bytes(), BulkData.Bytes()...))
			}

			_, err := http.Post("http://127.0.0.1:9210/_bulk", "application/json", combinedBulk)
			if err != nil {
				fmt.Println(err)
			}
			pg.Done()
		}()
	}
	pg.Wait()
}

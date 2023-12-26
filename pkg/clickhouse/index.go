package clickhouse

import (
	"context"
	"fmt"
	"fp/chvses/pkg/constants"
	"fp/chvses/pkg/utils"
	"log"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
)

func bulkInsertMany(conn *driver.Conn) error {
	defer utils.Timer("[CLICKHOUSE] bulkInsertMany")()
	var pg sync.WaitGroup
	pg.Add(constants.NUM_BATCHES)
	for i := 0; i < constants.NUM_BATCHES; i++ {
		go bulkInsert(conn, &pg)
	}
	pg.Wait()
	return nil
}

func bulkInsert(conn *driver.Conn, wg *sync.WaitGroup) error {
	// will be inserting into index_test table
	batch, err := (*conn).PrepareBatch(context.Background(), "INSERT INTO index_test")
	if err != nil {
		return err
	}

	for i := 0; i < constants.BATCH_SIZE; i++ {
		err := batch.Append(
			uuid.New(),
			"John",
			uint8(20),
			175,
			70,
			uint8(1),
		)
		if err != nil {
			return err
		}
	}
	err = batch.Send()
	wg.Done()
	return err
}

func connect() (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{"127.0.0.1:9000"},
			Auth: clickhouse.Auth{
				Database: "fp",
				Username: "default",
				Password: "",
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},

			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}

func logTables(conn *driver.Conn) {
	ctx := context.Background()
	rows, err := (*conn).Query(ctx, "SELECT name,toString(uuid) as uuid_str FROM system.tables")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			name, uuid string
		)
		if err := rows.Scan(
			&name,
			&uuid,
		); err != nil {
			log.Fatal(err)
		}
		log.Printf("name: %s, uuid: %s",
			name, uuid)
	}
}

func TestClickHouseIngestion() {
	conn, err := connect()
	if err != nil {
		panic((err))
	}

	err = bulkInsertMany(&conn)
	if err != nil {
		log.Fatal(err)
	}
}

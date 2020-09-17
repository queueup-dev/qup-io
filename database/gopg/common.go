package gopg

import (
	"fmt"
	"github.com/go-pg/pg/v9"
)

type GoPgConection struct {
	*pg.DB
}

type ConnectionWithTransacion interface {
	RunInTransaction(func(Connection) error) error
	Connection
}

type Connection interface {
	Insert(...interface{}) error
}

type PGRecord interface {
	GetId() int64
}

func (c GoPgConection) RunInTransaction(insert func(Connection) error) error {
	return c.DB.RunInTransaction(func(in *pg.Tx) error {
		return insert(in)
	})
}

func SetupConnection(dbUrlString string) (connection GoPgConection) {

	dbOptions, err := pg.ParseURL(dbUrlString)

	if err != nil {
		fmt.Println("error when parsing the db URL")
		panic(err)
	}

	connection.DB = pg.Connect(dbOptions)

	if _, dbConnectionError := connection.DB.Exec("SELECT 1"); dbConnectionError != nil {
		panic(fmt.Errorf("unable to perform test query on database"))
	}

	return
}

func StoreErrorRecord(db Connection, record PGRecord, err error) error {
	fmt.Println(err)
	err = db.Insert(record)

	if err != nil {
		fmt.Println("Failed to store body in database.")
		return err
	}
	fmt.Printf("Body stored in database, see record: %d\n", record.GetId())

	return nil
}

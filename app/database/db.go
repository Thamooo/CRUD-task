package database

import (
	"database/sql"
	aah "aahframe.work"
	_ "github.com/lib/pq"
  "github.com/juju/errors"
  "log"
  "context"
)


var Instance *sql.DB
var ctx context.Context
var TXconnection *sql.Tx

func ConnectTransaction() (){
  //log.Print("test")
  tx, err := Instance.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
  if(err != nil){
    //return errors.Annotate(err, "someone is already editing database")
    log.Fatal(err)
  }
  TXconnection = tx
  //return nil
}

func Connect(_ *aah.Event) {

	db, err := connect()
	if err != nil {
		panic(err)
	}
  Instance=db
}

func Disconnect(_ *aah.Event) {


	if err := Instance.Close(); err != nil {
		panic(errors.Annotate(err, "closing connection to database failed"))
    if(err!=nil){
      panic(err)
    }
	}
}

func CloseTransaction(){

  if rollbackErr := TXconnection.Rollback(); rollbackErr != nil {
		log.Print(rollbackErr)
	}

}


func connect() (*sql.DB, error) {
  var db *sql.DB
	db, err := sql.Open("postgres", "user=postgres dbname=test sslmode=disable")
	if err != nil {
		return nil, errors.Annotate(err, "connecting to database failed")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Annotate(err, "pinging database failed")
	}

	return db, errors.Annotate(err, "Connection successful")
}

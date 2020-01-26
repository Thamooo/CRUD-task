package database

import (
	"database/sql"
	"aahframe.work"
	_ "github.com/lib/pq"
  "github.com/juju/errors"
)


var Instance *sql.DB
var ActiveIDS = make(map[int]*sql.Tx) //LIST OF CURRENTLY EDITABLE ID

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

func CloseTransaction(id int){

  err := ActiveIDS[id].Rollback()
  if(err != nil){
    panic(err)
  }
  delete(ActiveIDS, id);
  return
}

func CommitTransaction(id int, tx *sql.Tx) (){
    err := tx.Commit()
  	if(err != nil){
  		panic(err)
  	}
    delete(ActiveIDS, id);
    return
}

func ConnectTransaction(id int) (*sql.Tx){

  tx, _ := Instance.Begin()

  ActiveIDS[id]=tx

  return tx

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

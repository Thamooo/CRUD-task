package database

import (
	"database/sql"
	aah "aahframe.work"
	_ "github.com/lib/pq"
  "github.com/juju/errors"
  // "log"
  "context"
  // "fmt"
)


var Instance *sql.DB
var ctx context.Context
var TXconnection *sql.Tx
var Stmt *sql.Stmt
var ActiveIDS = make(map[int]*sql.Tx)

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
}

func CommitTransaction(id int, tx *sql.Tx) (){
    err := tx.Commit()
  	if(err != nil){
  		panic(err)
  	}
    delete(ActiveIDS, id);
}

func ConnectTransaction(id int) (*sql.Tx){

  tx, _ := Instance.Begin()

  //statement := fmt.Sprintf(`UPDATE clients SET first_name = $1, last_name = $2, birth_date = $3, gender = $4, email = $5, address = $6 WHERE id = %v`, id)
  // stmt, err := TXconnection.Prepare(statement)
  // if(err != nil){
  //   panic(err)
  // }
  ActiveIDS[id]=tx

  return tx
 // if(err != nil){
 //   panic(err)
 // }
 // log.Print(stmt)
 //  Stmt = stmt

  //return nil
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

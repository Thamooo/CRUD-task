package session

import(
  "aahframe.work"
  "test-task/app/database"
  // "log"
)

type SessionController struct {
	*aah.Context
}

var session bool;

func (c *SessionController) ReserveSession() {


     database.ConnectTransaction();
    // if(err != nil){
    //   log.Print(err)
    // }
    // log.Print("Hey")
    // err :=
    // if(err != nil){
    //   c.Reply().JSON(aah.Data{
  	// 		"success": "false",
  	// 		"error":  "Someone is already editing",
  	// 	})
    // }
}

func (c *SessionController) RollbackSession() {

    database.CloseTransaction()

    // err :=
    // if(err != nil){
    //   c.Reply().JSON(aah.Data{
  	// 		"success": "false",
  	// 		"error":  "Someone is already editing",
  	// 	})
    // }
}

package session

import(
  "aahframe.work"
  "test-task/app/models"
  "test-task/app/database"
  "log"
  "encoding/json"
)

type SessionController struct {
	*aah.Context
}

func (c *SessionController) ReserveSession(id *models.ReservedID) {
  log.Print(id.ID)
    for k := range (database.ActiveIDS){
      log.Print(k)
      if(k == id.ID){
        c.Reply().JSON(aah.Data{
          "success": false,
          "errors": "Already editing" ,
        })
        return
      }
    }

    tx := database.ConnectTransaction(id.ID)

    rows, err := tx.Query(`SELECT * FROM clients WHERE id=$1`, id.ID)
    if(err != nil){
      panic(err)
    }
    defer rows.Close()

    var client models.Client
    for rows.Next() {
  		singleСlient := models.Client{}

  		err := rows.Scan(&singleСlient.ID, &singleСlient.Firstname, &singleСlient.Lastname, &singleСlient.Birthday, &singleСlient.Gender, &singleСlient.Email, &singleСlient.Address)
  		if err != nil {
  			panic(err)
  		}
  		client = singleСlient
  	}

      json_client, err := json.Marshal(client)
      if err != nil {
      }

      c.Reply().JSON(aah.Data{
        "success": true,
        "data":  string(json_client),
      })
      
      return
}

func (c *SessionController) RollbackSession(id *models.ReservedID) {

    
    for k := range (database.ActiveIDS){
      log.Print(k)
      if(k == id.ID){
        database.CloseTransaction(id.ID)
        return
      }

    }

}

package session

import(
  "aahframe.work"
  "test-task/app/models"
  "test-task/app/database"
  "log"
  "encoding/json"
  "fmt"
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
  		single_client := models.Client{}

  		err := rows.Scan(&single_client.ID, &single_client.Firstname, &single_client.Lastname, &single_client.Birthday, &single_client.Gender, &single_client.Email, &single_client.Address)
  		if err != nil {
  			return
  		}
  		client = single_client
  	}

      json_client, err := json.Marshal(client)
      if err != nil {
      }

      response := fmt.Sprintf(`{"success" : true, "data" : %v}`, string(json_client))

    	c.Reply().Ok().JSON(response)

}

func (c *SessionController) RollbackSession(id *models.ReservedID) {

    database.CloseTransaction(id.ID)

}

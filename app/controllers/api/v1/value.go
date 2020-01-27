package v1

import (
	"database/sql"
	"aahframe.work"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
	"test-task/app/models"
	"test-task/app/database"
	"regexp"
	"log"
)

type ValueController struct {
	*aah.Context
}



func (c *ValueController) EditClient(val *models.Client) {

	errors := ValidateUser(val, "edit")
	spew.Dump(errors)

	if len(errors) != 0 {
		c.Reply().JSON(aah.Data{
			"success": false,
			"errors":  errors,
		})
		return
	}

	tx := database.ActiveIDS[val.ID]

	res, err := tx.Exec(`UPDATE clients SET first_name = $1, last_name = $2, birth_date = $3, gender = $4, email = $5, address = $6 WHERE id = $7`, val.Firstname, val.Lastname, val.Birthday, val.Gender, val.Email, val.Address, val.ID)



	log.Print(res, err)
	rowsAffected, err := res.RowsAffected()
	if(err != nil){
		panic(err)
	}
	if(rowsAffected != 1){
		log.Print(rowsAffected)
	}

	database.CommitTransaction(val.ID, tx)

	c.Reply().JSON(aah.Data{
		"success": true,
		"errors":  errors,
	})

}

func (c *ValueController) DeleteClient(val *models.ToDeleteIDs) {

	var db = database.Instance

	for _, element := range val.IDs {
		_, err := db.Exec("DELETE FROM clients WHERE ID= $1",
			element)
		if(err != nil){
			c.Reply().JSON(aah.Data{
				"success": false,
			})
		}
	}



	c.Reply().JSON(aah.Data{
		"success": true,
	})
}

func (c *ValueController) AddClient(val *models.Client) {

	var db = database.Instance

	errors := ValidateUser(val, "add")
	spew.Dump(errors)

	if len(errors) != 0 {
		c.Reply().JSON(aah.Data{
			"success": false,
			"errors":  errors,
		})
		return
	}

	result, err := db.Exec("insert into clients (first_name, last_name, birth_date, gender, email, address) values ($1, $2, $3, $4, $5, $6)",
	val.Firstname, val.Lastname, val.Birthday, val.Gender, val.Email, val.Address)


	rowsAffected, err := result.RowsAffected()
	if(err != nil){

	}
	if(rowsAffected != 0){
		c.Reply().JSON(aah.Data{
			"success": false,
			"errors":  "{'main' : 'Database connection error'}",
		})
	}

	c.Reply().JSON(aah.Data{
		"success": true,
		"errors":  errors,
	})
}

func (c *ValueController) GetClients(id int, search string, sorting string) {
	var db = database.Instance

	searchQuery := " ";

	if(search != ""){
		var re = regexp.MustCompile(` `)
		s := re.ReplaceAllString(search, `%`)

		searchQuery = `AND Lower(first_name) || ' ' || Lower(last_name) LIKE '%%%v%%'`
		searchQuery = fmt.Sprintf(searchQuery, s)
	}
	if(sorting == ""){
		sorting="id ASC"
	}

	var howManyPages = fmt.Sprintf(`SELECT COUNT(*) FROM clients WHERE id > 0 %v`, searchQuery)

	pages, err := db.Query(howManyPages)
	if err != nil {
		panic(err)
	}
	defer pages.Close()

	var pagesCount int
	for pages.Next() {
	switch err := pages.Scan(&pagesCount); err {
	case sql.ErrNoRows:
		log.Print("Could not select any page")
	case nil:
		
	default:
	  panic(err)
		}
	}	

	totalPage := countPages(10, pagesCount);

	var sqlStatement = fmt.Sprintf(`SELECT t.id, t.first_name, t.last_name, t.birth_date, t.gender, t.email, t.address FROM (SELECT *, row_number() OVER(ORDER BY %v) AS row FROM clients WHERE id>0 %v) t WHERE t.row BETWEEN ($1 - 1) * 10 + 1 AND $1 * 10`, sorting, searchQuery)

	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	clients := []models.Client{}
	for rows.Next() {
		singleClient := models.Client{}
		err := rows.Scan(&singleClient.ID, &singleClient.Firstname, &singleClient.Lastname, &singleClient.Birthday, &singleClient.Gender, &singleClient.Email, &singleClient.Address)
		if err != nil {
			continue
		}
		clients = append(clients, singleClient)
	}

	jsonClients, err := json.Marshal(clients)
	if err != nil {
	}

	c.Reply().JSON(aah.Data{
		"pages": totalPage,
		"data":  string(jsonClients),
	})
	
	return
}
func countPages(totalPages int, items int) int{
	return (items + totalPages - 1) / totalPages
}

func ValidateUser(val *models.Client, validation_type string) map[string]string {

	checkData := make(map[string]bool)
	errors := make(map[string]string)

	checkData["Email"] = aah.App().ValidateValue(val.Email, "email,required")
	checkData["Firstname"] = aah.App().ValidateValue(val.Firstname, "min=2,max=100,required")
	checkData["Lastname"] = aah.App().ValidateValue(val.Lastname, "min=2,max=100,required")
	checkData["Gender"] = aah.App().ValidateValue(val.Gender, "oneof=male female|oneof=MALE FEMALE")
	checkData["Birthday"] = aah.App().ValidateValue(val.Birthday.String(), "age")
	checkData["Address"] = aah.App().ValidateValue(val.Address, "max=200")

	if(validation_type=="add"){
		checkData["EmailRegistr"] = aah.App().ValidateValue(val.Email, "emailRegistered,required")
	}

	spew.Dump(checkData)
	for k, v := range checkData {
		if v != true {
			switch k {
			case "Email":
				errors[k] = "Required, Should be valid email"
			case "Firstname":
				errors[k] = "Required, from 2 to 100 characters"
			case "Lastname":
				errors[k] = "Required, from 2 to 100 characters"
			case "Gender":
				errors[k] = "Required, allowed values are Male, Female"
			case "Address":
				errors[k] = "Optional, up to 200 characters"
			case "Birthday":
				errors[k] = "Required, from 18 till 60 years"
			case "EmailRegistr":
				errors["Email"] = "Email is already registered"
			}
		}
	}

	return errors
}

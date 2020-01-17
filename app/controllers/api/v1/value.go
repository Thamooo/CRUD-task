package v1

import (
	"aahframe.work"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
	"test-task/app/models"
	"regexp"
	//"gopkg.in/go-playground/validator.v9"
)

// ValueController is kickstart sample for API implementation.
type ValueController struct {
	*aah.Context
}

var response string

func (c *ValueController) EditClient(val *models.Client) {

	// client := models.Client{}
	// bytes := []byte(value)
	// json.Unmarshal(bytes, client)
	// //fmt.Println(client.Firstname)
	// fmt.Println(client)

	db := dbConnect()

	errors := ValidateUser(val, "edit")
	spew.Dump(errors)

	if len(errors) != 0 {
		c.Reply().JSON(aah.Data{
			"success": "false",
			"errors":  errors,
		})
		return
	}

	//c.Reply().Ok().JSON(response)
	//aah.App().Validate(val)
	fmt.Println(val.Birthday)
	sqlStatement := `
	UPDATE clients
	SET first_name = $2, last_name = $3, birth_date = $4, gender = $5, email = $6, address = $7
	WHERE id = $1`

	_, err := db.Exec(sqlStatement, val.ID, val.Firstname, val.Lastname, val.Birthday, val.Gender, val.Email, val.Address)

	//
	//rowsAffected, err := result.RowsAffected()
	if(err != nil){
		c.Reply().JSON(aah.Data{
			"success": "false",
			"errors":  "{'main' : 'Database connection error'}",
		})
	}
	// if(rowsAffected != 0){
	// 	c.Reply().JSON(aah.Data{
	// 		"success": "false",
	// 		"errors":  "{'main' : 'Database connection error'}",
	// 	})
	// }

	c.Reply().JSON(aah.Data{
		"success": "true",
		"errors":  errors,
	})
}

func (c *ValueController) DeleteClient(val *models.ToDeleteIDs) {

	db := dbConnect()

	for _, element := range val.IDs {
		_, err := db.Exec("DELETE FROM clients WHERE ID= $1",
			element)
		if(err != nil){
			c.Reply().JSON(aah.Data{
				"success": "false",
			})
		}
	}



	c.Reply().JSON(aah.Data{
		"success": "true",
	})
}

func (c *ValueController) AddClient(val *models.Client) {

	// client := models.Client{}
	// bytes := []byte(value)
	// json.Unmarshal(bytes, client)
	// //fmt.Println(client.Firstname)
	// fmt.Println(client)

	db := dbConnect()

	errors := ValidateUser(val, "add")
	spew.Dump(errors)

	if len(errors) != 0 {
		c.Reply().JSON(aah.Data{
			"success": "false",
			"errors":  errors,
		})
		return
	}

	//c.Reply().Ok().JSON(response)
	//aah.App().Validate(val)

	result, err := db.Exec("insert into clients (first_name, last_name, birth_date, gender, email, address) values ($1, $2, $3, $4, $5, $6)",
		val.Firstname, val.Lastname, val.Birthday, val.Gender, val.Email, val.Address)


	rowsAffected, err := result.RowsAffected()
	if(err != nil){

	}
	if(rowsAffected != 0){
		c.Reply().JSON(aah.Data{
			"success": "false",
			"errors":  "{'main' : 'Database connection error'}",
		})
	}

	c.Reply().JSON(aah.Data{
		"success": "true",
		"errors":  errors,
	})
}

func (c *ValueController) GetClients(id int, search string) {

	search_query := " ";
	//fmt.Println(order)
	if(search != ""){
		var re = regexp.MustCompile(` `)
		s := re.ReplaceAllString(search, `%`)

		search_query = `AND first_name || ' ' || last_name LIKE '%%%v%%'`
		search_query = fmt.Sprintf(search_query, s)
		fmt.Println(search_query)
	}
	db := dbConnect()
	//var sqlStatement = `select * from clients FETCH FIRST 10 ROWS ONLY` //STANDART RETURN 10 ROWS
	//var sqlStatement = fmt.Sprintf("select * FROM clients where id >= %v AND first_name LIKE '%Eri%' ORDER BY ID ASC FETCH FIRST 10 ROWS ONLY", id)
	var sqlStatement = fmt.Sprintf(`select * from clients where id >= $1 %v ORDER BY id FETCH FIRST 10 ROWS ONLY`, search_query)

	fmt.Println(sqlStatement)

	var how_many_pages = fmt.Sprintf(`SELECT t.id FROM (SELECT *, row_number() OVER(ORDER BY id) AS row FROM clients WHERE id > 0 %v) t WHERE t.row %% 10 = 0`, search_query)
	fmt.Println(how_many_pages)
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	pages, err := db.Query(how_many_pages)
	if err != nil {
		panic(err)
	}
	fmt.Println(pages)

	defer pages.Close()

	ids := []models.Pagination{}

	for pages.Next() {
		id_one := models.Pagination{}
		err := pages.Scan(&id_one.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ids = append(ids, id_one)
	}

	defer rows.Close()

	clients := []models.Client{}
	for rows.Next() {
		single_client := models.Client{}
		//spew.Dump(single_client.Birthday.String())
		err := rows.Scan(&single_client.ID, &single_client.Firstname, &single_client.Lastname, &single_client.Birthday, &single_client.Gender, &single_client.Email, &single_client.Address)
		if err != nil {
			fmt.Println(err)
			continue
		}
		clients = append(clients, single_client)
	}

	json_clients, err := json.Marshal(clients)
	if err != nil {
	}
	json_pages, err := json.Marshal(ids)
	if err != nil {
	}
	response = fmt.Sprintf(`{"pages" : %v, "data" : %v}`, string(json_pages), string(json_clients))

	//fmt.Sprintf(jsonData);
	c.Reply().Ok().JSON(response)
	return
}

func ValidateUser(val *models.Client, validation_type string) map[string]string {

	check_data := make(map[string]bool)
	errors := make(map[string]string)

	check_data["Email"] = aah.App().ValidateValue(val.Email, "email,required")
	check_data["Firstname"] = aah.App().ValidateValue(val.Firstname, "min=2,max=100,required")
	check_data["Lastname"] = aah.App().ValidateValue(val.Lastname, "min=2,max=100,required")
	check_data["Gender"] = aah.App().ValidateValue(val.Gender, "oneof=male female|oneof=MALE FEMALE")
	check_data["Birthday"] = aah.App().ValidateValue(val.Birthday.String(), "age")
	check_data["Address"] = aah.App().ValidateValue(val.Address, "max=200")
	if(validation_type=="add"){
		//check_data["ID"] = aah.App().ValidateValue(val.Email, "numeric,required")
		check_data["EmailRegistr"] = aah.App().ValidateValue(val.Email, "emailRegistered,required")
	}

	spew.Dump(check_data)
	for k, v := range check_data {
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
func dbConnect() *sql.DB {
	var db *sql.DB
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=test sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	return db
}

// List method returns all the values.
//func (c *ValueController) List() {
//
//}

// Index method returns value for given key.
// If key not found then returns 404 NotFound error.
//func (c *ValueController) Index() {
//
//	c.Reply().NotFound().JSON(aah.Data{
//		"message": "Value not exists",
//	})
//}

// Create method creates new entry in the values map with given payload.
// If key already exists then returns 409 Conflict error.
//func (c *ValueController) Create(val *models.Value) {
//	if _, found := values[val.Key]; found {
//		c.Reply().Conflict().JSON(aah.Data{
//			"message": "Key already exists",
//		})
//		return
//	}
//
//	// Add it to values map
//	values[val.Key] = val
//	newResourceURL := fmt.Sprintf("%s:%s", c.Req.Scheme, c.RouteURL("value_get", val.Key))
//	c.Reply().Created().
//		Header(ahttp.HeaderLocation, newResourceURL).
//		JSON(aah.Data{
//			"key": val.Key,
//		})
//}

// Update method updates value entry on map for given key and Payload.
// If key not exists then returns 400 BadRequest error.
//func (c *ValueController) Update(key string, val *models.Value) {
//	if r, found := values[key]; found {
//		r.Value = val.Value
//		values[key] = r
//		c.Reply().Ok().JSON(aah.Data{
//			"message": "Value updated successfully",
//		})
//		return
//	}
//
//	c.Reply().BadRequest().JSON(aah.Data{
//		"message": "Invalid input",
//	})
//}
//
//// Delete method deletes value for given key.
//// If key not exists then returns 400 BadRequest error.
//func (c *ValueController) Delete(key string) {
//	if _, found := values[key]; found {
//		delete(values, key)
//		c.Reply().Ok().JSON(aah.Data{
//			"message": "Value deleted successfully",
//		})
//		return
//	}
//
//	c.Reply().BadRequest().JSON(aah.Data{
//		"message": "Invalid input",
//	})
//}

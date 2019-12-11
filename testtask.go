package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/go-ozzo/ozzo-validation"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"html/template"
	"net/http"
	"time"
	"strconv"
)

type ViewData struct {
	Title   string
	Message string
}

type client struct {
	ID        int
	Firstname string
	Lastname  string
	Birthday  string
	Gender    string
	Email     string
	Address   string
}

func checkAge(value interface{}) error {
	s, _ := value.(string)
	birthdate, err := time.Parse("01-02-2006", s) //Parse string to date type
	if err != nil {
		// return errors.New("Incorrect Date")
		return errors.New("Incorrect date format")
	}
	now := time.Now()
	years := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		years--
	}
	if years < 18 || years > 60 {
		return errors.New("Age must be from 18 till 60 years")
	}
	return nil
}

func checkMail(value interface{}) error {
	s, _ := value.(string)
	err := checkmail.ValidateHost(s)
	if err != nil {
		fmt.Println(err)
		return errors.New("Email that you have entered is not valid")
	}
	return nil
}

func (c client) Validate() error {

	return validation.ValidateStruct(&c,
		// cannot be empty, and the length must between 2 and 100
		validation.Field(&c.Firstname, validation.Required, validation.Length(2, 100)),
		// cannot be empty, and the length must between 2 and 100
		validation.Field(&c.Lastname, validation.Required, validation.Length(2, 100)),
		// Gender cannot be empty, and the value must be 'MALE' or 'FEMALE'
		validation.Field(&c.Gender, validation.Required, validation.In("FEMALE", "MALE")),
		// checking that age is correct
		validation.Field(&c.Birthday, validation.Required, validation.By(checkAge)),
		//checking email
		validation.Field(&c.Email, validation.Required, validation.By(checkMail)),
		//checking address from 0 to 200 chars
		validation.Field(&c.Address, validation.Length(0, 200)),
	)
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=test sslmode=disable")
	if err != nil {
		return
	}
}

	func checkErr(err error) {
	    if err != nil {
	        panic(err)
	    }
		}

	func checkCount(rows *sql.Rows) (count int) {
	 	for rows.Next() {
	    	err:= rows.Scan(&count)
	    	checkErr(err)
	    }
	    return count
	}

func Get_clients(w http.ResponseWriter, r *http.Request) {
	var sqlStatement = `select * from clients FETCH FIRST 10 ROWS ONLY` //STANDART RETURN 10 ROWS

	//DETECT HOW MUCH PAGES TO RENDER
	rows, err := db.Query("SELECT COUNT(*) as count FROM clients")

	var counted_rows = checkCount(rows)
	var remainder = counted_rows%10;
	var count = 1;
	if remainder != 0 {
		count = (counted_rows/10)+1
	}else{
		count = (checkCount(rows)/10)
	}

	fmt.Printf("pages total : %v \n",count)

	rows, err = db.Query("Select ID FROM clients WHERE ID%10 = 0 ORDER BY ID")
	if err != nil {
	panic(err)
	}
	defer rows.Close()

	for rows.Next() {
	var id int
	err = rows.Scan(&id)
    if err != nil {
      // handle this error
      panic(err)
    }
		fmt.Println(id)
	}

	if r.Method == "POST" { //IF POST IS SET

		if(r.FormValue("id") == ""){
			return
		}
		id := r.FormValue("id")
		i, err := strconv.Atoi(id)
		if err != nil{
			return
		}

		sqlStatement = fmt.Sprintf("select * from clients where id < %v FETCH FIRST 10 ROWS ONLY", i)

		}

		rows, err = db.Query(sqlStatement)
		checkErr(err)

		defer rows.Close()

		clients := []client{}

		for rows.Next() {
			c := client{}
			err := rows.Scan(&c.ID, &c.Firstname, &c.Lastname, &c.Birthday, &c.Gender, &c.Email, &c.Address)
			if err != nil {
				fmt.Println(err)
				continue
			}
			clients = append(clients, c)
		}
		//fmt.Printf("%v", clients)
		jsonData, err := json.Marshal(clients)
		if err != nil {
			fmt.Println(err)
			return
		}

		full_response := fmt.Sprintf(`{ "pages":%v, "data":%v }`, count, string(jsonData))
		fmt.Println(full_response);
		w.Write([]byte(full_response))
}

func Add_client(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte(db_select()))
	if r.Method == "POST" {
		var client_to_add = client{0, r.FormValue("firstname"),
			r.FormValue("lastname"),
			r.FormValue("birthday"),
			r.FormValue("gender"),
			r.FormValue("email"),
			r.FormValue("address")}

		//check structure for errors
		err := client_to_add.Validate()
		if err != nil {
			fmt.Println(err.Error())
			jsonErr, _ := json.Marshal(err)
			full_response := fmt.Sprintf(`{ "success":false, "errors":%v }`, string(jsonErr))
			w.Write([]byte(string(full_response)))
			return
		}

		var email string
		sqlStatement := `SELECT email FROM clients WHERE email=$1;`
		row := db.QueryRow(sqlStatement, client_to_add.Email)
		switch err := row.Scan(&email); err {
		case sql.ErrNoRows:
			fmt.Println("Email is not registered, adding user to the database...")
		case nil:
			w.Write([]byte(string(`{ "success":false, "errors":{ "Email":"email is already registered" } }`)))
			fmt.Printf("User with email %v is already registered", client_to_add.Email)
			return
		default:
			panic(err)
		}

		birthday, _ := time.Parse("01-02-2006", client_to_add.Birthday)
		client_to_add.Birthday=birthday.Format("01/02/2006")
		fmt.Println(client_to_add.Birthday)
		//insert data in table
		fmt.Printf("insert into clients (first_name, last_name, birth_date, gender, email, address) values (%v, %v, %v, %v, %v, %v)", client_to_add.Firstname, client_to_add.Lastname, client_to_add.Birthday, client_to_add.Gender, client_to_add.Email, client_to_add.Address)
		result, err := db.Exec("insert into clients (first_name, last_name, birth_date, gender, email, address) values ($1, $2, $3, $4, $5, $6)",
			client_to_add.Firstname, client_to_add.Lastname, client_to_add.Birthday, client_to_add.Gender, client_to_add.Email, client_to_add.Address)

		rowsAffected, err := result.RowsAffected()
		checkErr(err)

		fmt.Printf("User with email %v is added %v rows affected", client_to_add.Email, rowsAffected)

		w.Write([]byte(string(`{ "success":true }`)))
		// fmt.Println("Receive ajax post data string ", client_to_add)
	}

}

func Home_page(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFiles("view/index.html")
	tmpl.Execute(w, nil)

}

func main() {
	// http.Handler
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", Home_page)
	// mux.HandleFunc("/receive", Send_ajax)
	http.HandleFunc("/", Home_page)
	http.HandleFunc("/clients", Get_clients)
	http.HandleFunc("/client/add", Add_client)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}

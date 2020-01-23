package main


import (
	"aahframe.work"
		"test-task/app/database"
	_ "aahframe.work/minify/html"
    "time"
    "gopkg.in/go-playground/validator.v9"
    "fmt"
	"strings"  
	"database/sql"

)

var app = aah.App()
var validate = aah.App().Validator()


func init() {


    validate.RegisterValidation("emailRegistered", checkMailRegistered)
	validate.RegisterValidation("age", ageValidator)
	
	app.OnStart(SubscribeHTTPEvents)

	app.OnStart(database.Connect)

	app.OnPostShutdown(database.Disconnect)

	app.HTTPEngine().Middlewares(
		aah.RouteMiddleware,
		aah.CORSMiddleware,
		aah.BindMiddleware,
		aah.AntiCSRFMiddleware,
		aah.AuthcAuthzMiddleware,
		aah.ActionMiddleware,
	)

}

func checkMailRegistered(fl validator.FieldLevel) bool{

    var db = database.Instance

    var email string
    emailUser := fl.Field().String()

    sqlStatement := "SELECT email FROM clients WHERE email=$1;"
    row := db.QueryRow(sqlStatement, emailUser)

    fmt.Println(emailUser)

    switch err := row.Scan(&email); err {
    case sql.ErrNoRows:

       return true
    case nil:
        return false
    default:
        return false
    }

}
func ageValidator(fl validator.FieldLevel) bool {

	date_string := fl.Field().String()
	

	spliited := strings.Split(date_string, " ")
	validationResult := checkAge(spliited[0], 18, 60)
	
	return validationResult
	

	
    }

	func checkAge(date string, minAge int, maxAge int) bool {
		birthdate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return false
		}

		now := time.Now()
		years := now.Year() - birthdate.Year()
		if now.YearDay() < birthdate.YearDay() {
			years--
		}
		if years < minAge || years > maxAge {
			return false
		}
		return true
	}

func SubscribeHTTPEvents(_ *aah.Event) {

}

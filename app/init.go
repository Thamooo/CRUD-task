// aah application initialization - configuration, server extensions, middleware's, etc.
// Customize it per application use case.

package main


import (
	"aahframe.work"
	// Registering HTML minifier
		"test-task/app/database"
	_ "aahframe.work/minify/html"
    "time"
    "gopkg.in/go-playground/validator.v9"
    "fmt"
		"strings"

)

var app = aah.App()
var validate = aah.App().Validator()


func init() {


    validate.RegisterValidation("emailRegistered", checkMailRegistered)
    validate.RegisterValidation("age", checkAge)
    //spew.Dump(app.Validator())
	//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
	// Server Extensions
	// Doc: https://docs.aahframework.org/server-extension.html
	//__________________________________________________________________________

	app.OnStart(SubscribeHTTPEvents)

	app.OnStart(database.Connect)

	//app.OnPostShutdown(database.CloseTransaction)
	app.OnPostShutdown(database.Disconnect)


	//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
	// Middleware's
	// Doc: https://docs.aahframework.org/middleware.html
	//
	// Executed in the order they are defined. It is recommended; NOT to change
	// the order of pre-defined aah framework middleware's.
	//__________________________________________________________________________
	app.HTTPEngine().Middlewares(
		aah.RouteMiddleware,
		aah.CORSMiddleware,
		aah.BindMiddleware,
		aah.AntiCSRFMiddleware,
		aah.AuthcAuthzMiddleware,

		//
		// NOTE: Register your Custom middleware's right here
		//

		aah.ActionMiddleware,
	)

	//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
	// Add Application Error Handler
	// Doc: https://docs.aahframework.org/error-handling.html
	//__________________________________________________________________________
	// app.SetErrorHandler(AppErrorHandler)

	//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
	// Add Custom Template Functions
	// Doc: https://docs.aahframework.org/template-funcs.html
	//__________________________________________________________________________
	// app.AddTemplateFunc(template.FuncMap{
	// // ...
	// })

	//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
	// Add Custom Session Store
	// Doc: https://docs.aahframework.org/session.html
	//__________________________________________________________________________
	// app.AddSessionStore("redis", &RedisSessionStore{})

	//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
	// Add Custom Value Parser
	// Doc: https://docs.aahframework.org/request-parameters-auto-bind.html
	//__________________________________________________________________________
	// if err := app.AddValueParser(reflect.TypeOf(CustomType{}), customParser); err != nil {
	//   log.Error(err)
	// }


}

// func dbConnect() *sql.DB {
// 	var db *sql.DB
// 	var err error
// 	db, err = sql.Open("postgres", "user=postgres dbname=test sslmode=disable")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return db
// }

func checkMailRegistered(fl validator.FieldLevel) bool{

    var db = database.Instance

    var email string
    emailUser := fl.Field().String()
    sqlStatement := "SELECT email FROM clients WHERE email=$1;"
    row := db.QueryRow(sqlStatement, emailUser)

    fmt.Println(emailUser)

    switch err := row.Scan(&email); err {
    //case sql.ErrNoRows:
		//	panic(err)
      //  return true
    case nil:
			panic(err)
        return false
    default:
        panic(err)
        return false
    }

}
func checkAge(fl validator.FieldLevel) bool {



	//s, _, _ := fl.GetStructFieldOK()
  date_string := fl.Field().String()
	//fmt.Println(date_string)
	//field := fl.Field().
	//fl.Field().Type()
	spliited := strings.Split(date_string, " ")
	// spew.Dump(date_string)
	fmt.Println(spliited[0])

	birthdate, err := time.Parse("2006-01-02", spliited[0]) //Parse string to date typetime.Parse("Mon Jan 2 15:04:05 2006", "15-12-2000")
	//birthdate, err := time.Parse("Mon Jan 2 15:04:05 2006", date_string) //Parse string to date typetime.Parse("Mon Jan 2 15:04:05 2006", "15-12-2000")
//birthdate.Format("02-01-2006")
fmt.Println(birthdate)
	//spew.Dump(fl.Field())
	//birthdate := reflect.ValueOf(fl.Field())

	if err != nil {
		// return errors.New("Incorrect Date")
		return false
	}
	now := time.Now()
	years := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		years--
	}
	if years < 18 || years > 60 {
		return false
	}
	return true
    }

// SubscribeHTTPEvents method subscribes to HTTP events on app start.
// Doc: https://docs.aahframework.org/server-extension.html
func SubscribeHTTPEvents(_ *aah.Event) {
	// he := aah.App().HTTPEngine()
	// he.OnRequest(myserverext.OnRequest)
	// he.OnPreReply(myserverext.OnPreReply)
	// he.OnHeaderReply(myserverext.OnHeaderReply)
	// he.OnPostReply(myserverext.OnPostReply)
	// he.OnPreAuth(myserverext.OnPreAuth)
	// he.OnPostAuth(myserverext.PostAuthEvent)
}

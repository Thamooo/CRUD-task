package models
import(
	"time"
	//"fmt"
)
//type BirthDate time.Time
// Value model is used by ValueController
type Value struct {
	Key   string      `json:"key" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
}

// type Value struct {
// 	Value []byte `json:"value" validate:"required"`
// }

type Pagination struct {
	ID int
}

type POSTdata struct {
	Value string `json:"value"`
}

type ToDeleteIDs struct{
		IDs []int `bind:"IDs"`
}

type Client struct {
	ID        int 	 `bind:"ID"`
	Firstname string `bind:"Firstname"` //validate:"min=2,max=100,required"`
	Lastname  string `bind:"Lastname"` //validate:"min=2,max=100,required"`
	Birthday  time.Time `bind:"Birthday"` //validate:"age,required"`
	Gender    string `bind:"Gender"` //validate:"oneof=male female|oneof=MALE FEMALE"`
	Email     string `bind:"Email"` //validate:"email,required"`
	Address   string `bind:"Address"` //validate:"max=100"`
}

// func (j *BirthDate) DateParse(b string) error {
//     t, err := time.Parse("02-01-2006", b)
//     if err != nil {
//         return err
//     }
// 		fmt.Println(t)
//     *j = BirthDate(t)
//     return nil
// }
//
// func (j *BirthDate) String() string {
//     fmt.Println(j);
// 		return "lol"
// }

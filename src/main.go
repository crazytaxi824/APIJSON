package main

import (
	"log"
	"src/action"
)

type person struct {
	tableName struct{} `sql:"person"`
	Name      string   `sql:"name" json:"personName"`
	Age       int      `sql:"age" json:"personAge"`
}

func main() {
	log.SetFlags(log.Lshortfile)

	s := `{
"Person":{
},
  "[]": {
	"page":2,
	"count":10,
	"join": "&/User/id@,</Comment/momentId@",
    "Moment":{
        "@column":"id,date,userId",
        "id":12
      },
    "User":{
      "id@":"/Moment/userId",
      "@column":"id,name"
    }
  }
}`
	r, err := action.ParseJsonRequest([]byte(s))
	if err != nil {
		log.Println(err)
		return
	}

	err = action.ParseTableName(r)
	if err != nil {
		log.Println(err)
		return
	}

}

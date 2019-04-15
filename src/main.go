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

	action.CacheStruct(person{})

	log.Println(action.TableJsonStructMap)
}

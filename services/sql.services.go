package services

import (
	"fmt"
	"goGinServer/db"
)

type Result struct {
	ActorID   int
	FirstName string
	LastName  string
}

func GroupByService() []*Result {
	db := db.DbConnection()

	var result []*Result
	db.Raw("SELECT first_name, actor_id, last_name FROM actor").Scan(&result)

	for _, item := range result {
		fmt.Println(item)
	}

	return result
}

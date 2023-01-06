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

	// Query using raw SQL sintax and GORM orm to get all the users in the users table
	// and group them by their age.
	// rows, err := db.Raw("SELECT actor_id, first_name, last_name FROM actor").Rows()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // Print the rows
	// for rows.Next() {
	// 	var actor_id int
	// 	var first_name string
	// 	var last_name string
	// 	rows.Scan(&actor_id, &first_name, &last_name)
	// 	fmt.Println(actor_id, first_name, last_name)
	// }

	var result []*Result
	db.Raw("SELECT first_name, actor_id, last_name FROM actor").Scan(&result)

	for _, item := range result {
		fmt.Println(item)
	}

	return result
}

package main

import (
	"app/stores"
	"app/stores/models"
	"app/stores/mysql"
	"app/stores/postgress"
	"fmt"
)

func main() {
	//Call postgres and mysql package methods using interface variable
	var i stores.DataBase
	u1 := models.User{1, "sandra"}
	u2 := models.User{2, "komal"}
	u3 := models.User{3, "diwakar"}

	i = mysql.NewConn("mysql")
	fmt.Printf("---------------- Printing type of interface %T ---------------------------------\n", i)
	printMapReturn(i.FetchAll())

	printReturn(i.Create(u1))

	printReturn(i.FetchUser(1))

	printReturn(i.Update(1, "new me"))
	printReturn(i.Delete(2))

	printReturn(i.Create(u2))
	printMapReturn(i.FetchAll())

	printReturn(i.Delete(2))

	printMapReturn(i.FetchAll())

	i = postgress.NewConn("postgress")
	fmt.Printf("---------------- Printing type of interface %T ---------------------------------\n", i)
	printMapReturn(i.FetchAll())

	printReturn(i.Create(u1))

	printReturn(i.FetchUser(1))

	printReturn(i.Update(1, "new me"))
	printReturn(i.Delete(2))

	printReturn(i.Create(u2))
	printMapReturn(i.FetchAll())

	printReturn(i.Delete(2))

	printReturn(i.Create(u3))

	printMapReturn(i.FetchAll())

}

func printReturn(u *models.User, ok bool) {
	if !ok {
		fmt.Println("The command was not executed")
	}
	if u == nil {
		fmt.Println("u :", nil, " ok ", ok, "\n\n")
	} else {
		fmt.Println("u :", u, " ok ", ok, "\n\n")
	}

}

func printMapReturn(userDb map[int]*models.User, ok bool) {
	if !ok {
		fmt.Println("The command was not executed")
	}
	for key, value := range userDb {
		fmt.Println("Fetched values from db")
		fmt.Printf("Key: %s\tValue: %v\n", key, value)
	}

}

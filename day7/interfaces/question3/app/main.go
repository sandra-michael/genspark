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
	u := models.User{1, "sandra"}

	i = mysql.NewConn("mysql")
	fmt.Printf("---------------- Printing type of interface %T ---------------------------------\n", i)
	i.Create(u)
	i.Update("new Name")
	i.Delete(2)

	i = postgress.NewConn("postgress")
	fmt.Printf("---------------- Printing type of interface %T ---------------------------------\n", i)
	i.Create(u)
	i.Update("new Name")
	i.Delete(2)

}

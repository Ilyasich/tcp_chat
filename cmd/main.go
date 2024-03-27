package main

import (
	"git/server"

	// "context"
	// "database/sql"
	// "fmt"
	// "log"

	//_ "github.com/lib/pq"
)

func main() {

	// 	connStr := "postgres://postgres:mysecretpassword@localhost/app_db?sslmode=disable"
// 	db, err := sql.Open("postgres", connStr) //универсальная гошная библиотека имеет Open соеденение и принимает два параметра имя драйвера и строка подключения
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	//добавляем в бд строки
// 	_, err = db.ExecContext(context.Background(), `INSERT INTO authors (user_name) VALUES
// ($1)`, "Dima", "Zyz")

// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	err = db.Ping() //создает подключение и дергает базу данных
// 	fmt.Println(err)


	server.StartServer()
}

package main

import (
	"log"
	"os"
)

func main() {

	// var abc string = "ABC"
	// fmt.Print(abc)

	args := Args{
		conn: "postgres://postgres:root@localhost/dana?sslmode=disable",
		port: ":8080",
	}
	if conn := os.Getenv("DB_CONN"); conn != "" {
		args.conn = conn
	}
	if port := os.Getenv("PORT"); port != "" {
		args.port = ":" + port
	}
	// run server
	if err := Run(args); err != nil {
		log.Println(err)
	}
}

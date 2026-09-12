package main

import (
	"flag"
	"fmt"
	"log"

	"agora-backend/internal/config"
	"agora-backend/internal/db"
)

func main() {
	username := flag.String("username", "", "existing username to update")
	action := flag.String("action", "promote", "promote or demote")
	flag.Parse()
	if *username == "" || (*action != "promote" && *action != "demote") {
		log.Fatal("usage: go run ./cmd/admin -username NAME -action promote|demote")
	}
	database, err := db.InitDB(config.LoadConfig().DBDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	role := "admin"
	if *action == "demote" {
		role = "user"
	} else {
		var email string
		if err := database.QueryRow(`SELECT email FROM users WHERE username=$1`, *username).Scan(&email); err != nil {
			log.Fatalf("user %q not found", *username)
		}
		if email == "" {
			log.Fatal("cannot promote an account without an email address")
		}
	}
	result, err := database.Exec(`UPDATE users SET role=$2,updated_at=CURRENT_TIMESTAMP WHERE username=$1`, *username, role)
	if err != nil {
		log.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	if rows != 1 {
		log.Fatalf("user %q not found", *username)
	}
	fmt.Printf("user %q role updated to %s\n", *username, role)
}

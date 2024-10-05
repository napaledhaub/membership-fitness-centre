package main

import (
	"database/sql"
	"fmt"
	"log"
	"membership-fitness-centre/controllers"
	"membership-fitness-centre/middleware"
	"membership-fitness-centre/services"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	memberService := services.NewMemberService(db)
	memberController := controllers.NewMemberController(memberService)

	http.Handle("/members", middleware.Logging(http.HandlerFunc(memberController.GetMembers)))
	http.Handle("/members/create", middleware.Logging(http.HandlerFunc(memberController.CreateMember)))
	http.Handle("/members/update", middleware.Logging(http.HandlerFunc(memberController.UpdateMember)))
	http.Handle("/members/delete", middleware.Logging(http.HandlerFunc(memberController.DeleteMember)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

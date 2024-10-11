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
	"github.com/robfig/cron/v3"
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

	packageService := services.NewPackageService(db)
	packageController := controllers.NewPackageController(packageService)

	subscriptionService := services.NewSubscriptionervice(db)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	c := cron.New()
	c.AddFunc("@every 1m", subscriptionController.CheckSubscriptions)
	c.Start()

	http.Handle("/members/create", middleware.Logging(http.HandlerFunc(memberController.CreateMember)))
	http.Handle("/verify", middleware.Logging(http.HandlerFunc(memberController.VerifyEmailHandler)))
	http.Handle("/members/login", middleware.Logging(http.HandlerFunc(memberController.Login)))
	http.Handle("/member/update_password", middleware.AuthMiddleware(http.HandlerFunc(memberController.UpdatePassword)))

	http.Handle("/package/add_package", middleware.AuthMiddleware(http.HandlerFunc(packageController.AddPackage)))

	log.Fatal(http.ListenAndServe(":8080", nil))

	select {}
}

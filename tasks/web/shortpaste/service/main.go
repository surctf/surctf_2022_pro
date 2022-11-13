package main

import (
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"

	"shortpaste/config"
	"shortpaste/handlers"
	"shortpaste/middleware"
	"shortpaste/models"
)

func AddFlag(db *gorm.DB) {
	u1_pwd, _ := bcrypt.GenerateFromPassword([]byte("ashd!@#kahsdh12kh3hasjdh$@#!jahsjh12jh3jhjsdhajshd!"), bcrypt.DefaultCost)
	u1 := models.User{PasswordHash: u1_pwd}
	db.Create(&u1)
	u1.AddPaste(&models.Paste{Title: "NE FLAG", Content: config.FLAG}, db)
}

func main() {
	os.Remove("data.db")

	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Paste{})
	AddFlag(db)

	auth := middleware.AuthMiddleware{db}
	h := handlers.NewHandlers(db)

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Use(auth.Middleware)

	rGET := r.Methods("GET").Subrouter()
	rGET.HandleFunc("/", http.HandlerFunc(h.IndexHandler))
	rGET.HandleFunc("/paste/{paste_id}", http.HandlerFunc(h.GetPasteHandler))
	rGET.HandleFunc("/logout", http.HandlerFunc(h.LogOutHandler))

	rPOST := r.Methods("POST").Subrouter()
	rPOST.HandleFunc("/sign_up", http.HandlerFunc(h.SignUpHandler))
	rPOST.HandleFunc("/sign_in", http.HandlerFunc(h.SignInHandler))
	rPOST.HandleFunc("/new_paste", http.HandlerFunc(h.CreatePasteHandler))

	s := http.Server{
		Addr:         config.PORT,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      r,
	}

	log.Print("Serving on: ", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

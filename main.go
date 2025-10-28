package main

import (
	"Final_Project/internal/db"
	"Final_Project/internal/server"
	"Final_Project/internal/settings"
	"log"
)

func main() {
	cfg := settings.Load()

	if err := db.Init(); err != nil {
		log.Fatalf("ошибка инициализации БД: %v", err)
	}
	defer db.DB.Close()

	s := server.New(cfg)
	if err := s.Run(); err != nil {
		log.Fatalf("Ошибка запуска сервера")
	}
}

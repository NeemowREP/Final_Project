package main

import (
	"Final_Project/internal/server"
	"Final_Project/internal/settings"
	"log"
)

func main() {
	cfg := settings.Load()

	s := server.New(cfg)
	if err := s.Run(); err != nil{
		log.Fatalf("Ошибка запуска сервера")
	}
}
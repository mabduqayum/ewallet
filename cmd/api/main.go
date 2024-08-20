package main

import (
	"ewallet/internal/server"
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	s := server.New()
	s.RegisterFiberRoutes()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := s.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}

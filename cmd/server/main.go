package main

import _ "github.com/joho/godotenv/autoload"

func main() {
	app := InitializeApp()
	app.Start()
}

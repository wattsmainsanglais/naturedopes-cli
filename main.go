package main

import (
	"github.com/joho/godotenv"
	"github.com/wattsmainsanglais/naturedopes-cli/cmd"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}

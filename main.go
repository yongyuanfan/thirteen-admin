package main

import (
	"errors"
	"log"
	"os"

	"clash-admin/cmd"

	"github.com/joho/godotenv"
)

// main 加载环境变量并执行根命令。
func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("failed to load .env: %v", err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

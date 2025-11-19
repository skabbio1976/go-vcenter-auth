package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	vcauth "github.com/skabbio1976/go-vcenter-auth"
)

func main() {
	// Read credentials from environment variables
	host := os.Getenv("VCENTER_HOST")
	username := os.Getenv("VCENTER_USERNAME")
	password := os.Getenv("VCENTER_PASSWORD")

	if host == "" || username == "" || password == "" {
		log.Fatal("Please set VCENTER_HOST, VCENTER_USERNAME, and VCENTER_PASSWORD environment variables")
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Login to vCenter
	client, err := vcauth.Login(ctx, host, username, password, true)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	fmt.Println("Successfully logged in to vCenter!")

	// Get the vim25 client for further operations
	vim := client.GetVim()
	fmt.Printf("Connected to: %s\n", vim.URL().Host)
	fmt.Printf("API Version: %s\n", vim.ServiceContent.About.ApiVersion)
	fmt.Printf("Product: %s %s\n", vim.ServiceContent.About.FullName, vim.ServiceContent.About.Version)
}

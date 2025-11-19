//go:build windows

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
	// Read vCenter host from environment variable
	host := os.Getenv("VCENTER_HOST")
	if host == "" {
		log.Fatal("Please set VCENTER_HOST environment variable")
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Login using Windows integrated authentication (current user)
	fmt.Println("Attempting SSPI login with current Windows user...")
	client, err := vcauth.LoginSSPI(ctx, host, true)
	if err != nil {
		log.Fatalf("Failed to login via SSPI: %v", err)
	}

	fmt.Println("Successfully logged in to vCenter via SSPI!")

	// Get the vim25 client for further operations
	vim := client.GetVim()
	fmt.Printf("Connected to: %s\n", vim.URL().Host)
	fmt.Printf("API Version: %s\n", vim.ServiceContent.About.ApiVersion)
	fmt.Printf("Product: %s %s\n", vim.ServiceContent.About.FullName, vim.ServiceContent.About.Version)
}

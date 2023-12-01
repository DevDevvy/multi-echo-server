package serverPackage

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// This file contains the code for the TCP server.
// It first creates a WaitGroup to wait for all goroutines to finish before shutting down.
// Then it starts a TCP server on port 8080.
// It uses a goroutine to wait for a shutdown signal.
// It accepts incoming connections and creates a goroutine to handle each connection.
// The handleConnection function reads data from the connection and echoes it back to the client.

// Server starts a TCP server on port 8080
func Server() {
	// Use a WaitGroup to wait for all goroutines to finish before shutting down
	var wg sync.WaitGroup

	// Start a TCP server on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	//os.Interrupt is the signal sent, and syscall.SIGTERM is the signal type
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	//the go routine is used to run the function in the background
	//this allows the program to continue running while the function is running
	go func() {
		// Wait for a shutdown signal
		<-stop

		fmt.Println("\nShutting down the server...")

		// Close the listener to stop accepting new connections
		listener.Close()

		// Wait for all active connections to finish
		wg.Wait()

		// Exit the program
		os.Exit(0)
	}()

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Increment the WaitGroup counter for each new connection
		wg.Add(1)
		go handleConnection(conn, &wg)
	}
}

// Handle incoming connections
func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	// Ensure that the WaitGroup counter is decremented when the function exits
	defer func() {
		// Close the connection
		conn.Close()

		// Decrement the WaitGroup counter when the connection is closed
		wg.Done()
	}()

	// Read data from the connection
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	// Echo the data back to the client
	_, err = conn.Write(buffer)
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		return
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
)

const (
	defaultLocalPort  = 19132
	defaultRemotePort = 19132
)

func log(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	gray := color.New(color.FgHiBlack)
	boldWhite := color.New(color.FgWhite, color.Bold)

	gray.Printf("[%s] ", timestamp)
	boldWhite.Println(message)
}

// don't flush away the error message on failure
func exitWithError(message string) {
	color.New(color.FgRed, color.Bold).Fprintf(os.Stderr, "%s\n", message)
	fmt.Println("\nPress Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}

func promptForInput(prompt string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	blueBold := color.New(color.FgBlue, color.Bold)
	if defaultValue != "" {
		blueBold.Printf("%s", prompt)
		fmt.Printf(" [%s]: ", defaultValue)
	} else {
		blueBold.Printf("%s: ", prompt)
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		exitWithError(fmt.Sprintf("Error reading input: %v", err))
	}

	input = strings.TrimSpace(input)
	if input == "" && defaultValue != "" {
		return defaultValue
	}
	return input
}

func validatePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("port must be a valid number (got '%s')", portStr)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port must be between 1 and 65535 (got %d)", port)
	}

	return port, nil
}

func showInteractiveUI() {
	PrintTitle()
	PrintDiagram()
	PrintDescription()
	PrintWarning()
	PrintConfigPrompt()
}

func main() {
	remoteHost := flag.String("host", "", "Remote host (hostname or IPv4 address)")
	remotePort := flag.Int("port", defaultRemotePort, "Remote port (1-65535)")
	flag.Parse()

	var host string
	var port int

	if *remoteHost == "" {
		// no cli flags given
		showInteractiveUI()

		for {
			host = promptForInput("Enter remote host (hostname or IP)", "")
			if host != "" {
				break
			}
			color.New(color.FgRed, color.Bold).Println("Remote host cannot be empty. Please try again.")
		}

		portStr := promptForInput("Enter remote port", fmt.Sprintf("%d", defaultRemotePort))
		var err error
		port, err = validatePort(portStr)
		if err != nil {
			exitWithError(fmt.Sprintf("Error: %v", err))
		}
	} else {
		// use cli flags
		host = *remoteHost
		port = *remotePort

		// Validate port
		if port < 1 || port > 65535 {
			exitWithError("Error: port must be between 1 and 65535")
		}
	}

	fmt.Println()
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Local:  0.0.0.0:%d (UDP)\n", defaultLocalPort)
	fmt.Printf("  Remote: %s:%d (UDP)\n", host, port)
	fmt.Println()

	// let's do it
	remoteAddr := fmt.Sprintf("%s:%d", host, port)
	localAddr := fmt.Sprintf("0.0.0.0:%d", defaultLocalPort)
	proxy := NewUDPProxy(localAddr, remoteAddr)
	err := proxy.Start()
	if err != nil {
		exitWithError(fmt.Sprintf("Error: %v", err))
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log("Proxy is running. Press Ctrl+C to stop.")

	// cleanup
	<-sigChan
	log("Received interrupt signal, shutting down...")
	proxy.Stop()
	log("Shutdown complete")
}

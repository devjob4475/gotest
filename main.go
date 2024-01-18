package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8888"
	}

	http.HandleFunc("/", HelloHandler)
	http.HandleFunc("/runcmd", CmdHandler)

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func HelloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello from Koyeb")
}

func CmdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // 10 MB max size
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		cmd := r.FormValue("cmd")
		if cmd == "" {
			http.Error(w, "Command is required", http.StatusBadRequest)
			return
		}

		// Execute the command (for demonstration purposes only, this is unsafe)
		output, err := runCommand(cmd)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing command: %s", err), http.StatusInternalServerError)
			return
		}

		// Return the command output as the response
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, output)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func runCommand(cmd string) (string, error) {
	// WARNING: This is a simple example and is NOT secure for production use!
	// You should validate and sanitize user input before executing commands.
	// Using unsafe commands can lead to security vulnerabilities.

	result, err := exec.Command(cmd).Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}

	return string(result), nil
}

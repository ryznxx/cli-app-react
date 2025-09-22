package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func getStoragePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("Failed to get home dir:", err)
		home = "/tmp"
	}
	storagePath := filepath.Join(home, "storage")

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		os.MkdirAll(storagePath, 0755)
	}
	return storagePath
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	storagePath := getStoragePath()
	log.Println("User storage path:", storagePath)

	// Spawn shell per connection
	cmd := exec.Command("bash")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println("Failed to get stdin:", err)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Failed to get stdout:", err)
		return
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		log.Println("Failed to start shell:", err)
		return
	}

	// Baca output shell dan kirim ke client
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				break
			}
			if n > 0 {
				if writeErr := conn.WriteMessage(websocket.TextMessage, buf[:n]); writeErr != nil {
					break
				}
			}
		}
	}()

	// Terima input dari client dan kirim ke shell
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Optional: simpan input ke storage file
		inputFile := filepath.Join(storagePath, "last_input.txt")
		os.WriteFile(inputFile, message, 0644)

		stdin.Write(message)
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

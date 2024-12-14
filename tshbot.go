package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

// Config structure to hold the configuration data
type Config struct {
	BotLogFile  string            `yaml:"bot_log_file"`
	BashCmd     string            `yaml:"bash_cmd"`
	TGBotToken  string            `yaml:"tg_bot_token"`
	TGBotChatID string            `yaml:"tg_bot_chat_id"`
	AllowedCmds map[string]string `yaml:"allowed_cmds"`
	HelpMessage string            `yaml:"help_message"`
}

var config Config
var reservedWords = []string{"help", "commands"}

func init() {

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user's home directory: %v", err)
	}

	// Construct the full path to the configuration file
	configFilePath := filepath.Join(homeDir, ".config/tshbot/tshbot.config")

	// Load the configuration file
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Error loading configuration file: %v", err)
	}

	// Decode config
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("Error decoding configuration file: %v", err)
	}

	// Validate token presence
	if config.TGBotToken == "" || config.TGBotChatID == "" {
		log.Fatal("Configuration must include TGBotToken and TGBotChatID")
	}

	// Validate allowed commands
	if err := validateAllowedCommands(config.AllowedCmds); err != nil {
		log.Fatal(err)
	}
}

func main() {

	// Setup logging (open or create the log file)
	f, err := os.OpenFile(config.BotLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(config.TGBotToken)
	if err != nil {
		log.Panicf("Error initializing bot: %v", err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Send startup message with system info
	sendStartupMessage(bot)

	// Create a new UpdateConfig struct with an offset of 0.
	// Offsets are used to make sure Telegram knows we've handled previous values and we don't need them repeated.
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 40 // Set the timeout to 40 seconds to wait for new updates

	// Start polling Telegram for updates
	updates := bot.GetUpdatesChan(u)

	// Go through each update that we're getting from Telegram
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := strconv.Itoa(int(update.Message.Chat.ID))
		if chatID == config.TGBotChatID {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			handleCommand(update.Message, bot)
		}
	}
}

func sendStartupMessage(bot *tgbotapi.BotAPI) {
	// Get the current time
	currentTime := time.Now().Format(time.RFC1123)

	// Get the current user
	userName, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error retrieving user home directory: %v", err)
		userName = "unknown"
	} else {
		userName = filepath.Base(userName)
	}

	// Get the host name
	hostName, err := os.Hostname()
	if err != nil {
		log.Printf("Error retrieving hostname: %v", err)
		hostName = "unknown"
	}

	// Get the external IP address
	resp, err := http.Get("https://api.ipify.org?format=text")
	var ipAddress string
	if err != nil {
		log.Printf("Error fetching external IP address (using ipify.org): %v", err)
		ipAddress = "unknown"
	} else {
		// Ensure response body is closed after reading
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading external IP address response: %v", err)
			ipAddress = "unknown"
		} else {
			ipAddress = string(body)
		}
	}

	// Construct the startup message
	startupMessage := "tshbot started!\n" +
		"Time: " + currentTime + "\n" +
		"User: " + userName + "\n" +
		"Host: " + hostName + "\n" +
		"External IP: " + ipAddress

	// Send the startup message
	chatID, _ := strconv.ParseInt(config.TGBotChatID, 10, 64)
	sendMessage(chatID, startupMessage, bot)
}

func handleCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	text := message.Text
	chatID := message.Chat.ID
	userID := message.From.ID

	log.Printf("Received command: '%s' from user ID: %d, Username: %s", text, userID, message.From.UserName)

	// Handle specific commands first
	switch text {
	case "/help":
		sendMessage(chatID, config.HelpMessage, bot)
		return
	case "/commands":
		commandsList := ""
		for shortcut, fullCmd := range config.AllowedCmds {
			commandsList += "/" + shortcut + " - " + fullCmd + "\n"
		}
		sendMessage(chatID, "Available commands:\n"+commandsList, bot)
		return
	}

	// Check for user-allowed commands with user-defined shortcuts
	if strings.HasPrefix(text, "/") {
		cmdShortcut := strings.TrimPrefix(text, "/")
		fields := strings.Fields(cmdShortcut)
		if len(fields) > 0 {
			cmdShortcut = fields[0]
		} else {
			sendMessage(chatID, "Invalid command format. Use /commands to see available commands, or type /help.", bot)
			return
		}

		if cmdShortcut == "shell" {
			// Handle the special case for shell shortcut
			cmd := strings.TrimSpace(strings.TrimPrefix(text, "/shell"))
			if cmd == "" {
				sendMessage(chatID, "Please provide a command to execute.", bot)
				return
			}
			output := execShellCommand(cmd, true)
			sendMessage(chatID, output, bot)
		} else if fullCmd, ok := isAllowedCommand(cmdShortcut); ok {
			output := execShellCommand(fullCmd, false)
			sendMessage(chatID, output, bot)
		} else {
			log.Printf("Command not recognized or allowed: %s", cmdShortcut)
			sendMessage(chatID, "Command not recognized or allowed", bot)
		}
	} else {
		sendMessage(chatID, "Commands should start with /. Use /commands to see available commands, or type /help.", bot)
	}
}

// Function to execute commands in bash, with optional highlighting (e.g., for arbitrary `shell` commands)
func execShellCommand(command string, highlight bool) string {
	var logMessage string
	if highlight {
		logMessage = strings.Repeat("*", 10) + " Executing command: " + command + " " + strings.Repeat("*", 10)
	} else {
		logMessage = "Executing command: " + command
	}
	log.Printf("[%s] %s", time.Now().Format(time.RFC3339), logMessage)
	cmd := exec.Command(config.BashCmd, "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while running command: %v", err)
		return fmt.Sprintf("Error executing command: %s\n%s", err, string(out))
	}
	return string(out)
}

func sendMessage(chatID int64, text string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func isAllowedCommand(shortcut string) (string, bool) {
	fullCmd, ok := config.AllowedCmds[shortcut]
	return fullCmd, ok
}

func validateAllowedCommands(commands map[string]string) error {
	seen := make(map[string]bool)
	for shortcut := range commands {
		if contains(reservedWords, shortcut) {
			return errors.New("shortcut '" + shortcut + "' is a reserved word and cannot be used")
		}
		if seen[shortcut] {
			return errors.New("shortcut '" + shortcut + "' is duplicated")
		}
		if shortcut == "shell" {
			log.Println("WARNING: 'shell' shortcut is allowed. This poses a potential security risk as it allows arbitrary command execution.")
		}
		seen[shortcut] = true
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

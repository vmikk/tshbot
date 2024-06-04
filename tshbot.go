package main

import (
	"log"
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
	BotLogFile  string   `yaml:"bot_log_file"`
	BashCmd     string   `yaml:"bash_cmd"`
	TGBotToken  string   `yaml:"tg_bot_token"`
	TGBotChatID string   `yaml:"tg_bot_chat_id"`
	AllowedCmds []string `yaml:"allowed_cmds"`
	HelpMessage string   `yaml:"help_message"`
}

var config Config

func init() {

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting user's home directory:", err)
	}

	// Construct the full path to the configuration file
	configFilePath := filepath.Join(homeDir, ".config/tshbot/tshbot.config")

	// Load the configuration file
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal("Error loading configuration file:", err)
	}

	// Decode config
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal("Error decoding configuration file:", err)
	}

	// Validate token presence
	if config.TGBotToken == "" || config.TGBotChatID == "" {
		log.Fatal("Configuration must include TGBotToken and TGBotChatID")
	}
}

func main() {

	// Setup logging
	f, err := os.OpenFile(config.BotLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(config.TGBotToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Create a new UpdateConfig struct with an offset of 0.
	// Offsets are used to make sure Telegram knows we've handled previous values and we don't need them repeated.
	u := tgbotapi.NewUpdate(0)

	// Wait up to 40 seconds on each request for an update
	u.Timeout = 40

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

func handleCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	text := message.Text
	chatID := message.Chat.ID
	userID := message.From.ID

	log.Printf("Received command: %s from user: %d", text, userID)

	switch {
	case strings.HasPrefix(text, "/sh "):
		cmd := strings.TrimPrefix(text, "/sh ")
		if cmd == "" {
			log.Println("Bot missing parameter")
			sendMessage(chatID, "Missing command parameter", bot)
			return
		}
		if isAllowedCommand(cmd) {
			output := execShellCommand(cmd)
			sendMessage(chatID, output, bot)
		} else {
			log.Println("Command not recognized or allowed")
			sendMessage(chatID, "Command not recognized or allowed", bot)
		}

	case text == "/sysinfo":
		if isAllowedCommand("sysinfo") {
			output := execShellCommand("df -h && free -m")
			sendMessage(chatID, output, bot)
		} else {
			log.Println("Command not allowed")
			sendMessage(chatID, "Command not allowed", bot)
		}

	case text == "/help":
		sendMessage(chatID, config.HelpMessage, bot)

	case text == "/commands":
		commandsList := strings.Join(config.AllowedCmds, "\n")
		sendMessage(chatID, "Available commands:\n"+commandsList, bot)

	default:
		sendMessage(chatID, "Unknown command. Use /help for available commands.", bot)
	}
}

func execShellCommand(command string) string {
	log.Printf("[%s] Executing command: %s %s %s", time.Now().Format(time.RFC3339), config.BashCmd, "-c", command)
	out, err := exec.Command(config.BashCmd, "-c", command).Output()
	if err != nil {
		log.Println("Error while running command:", err)
		return "Error executing command"
	}
	return string(out)
}

func sendMessage(chatID int64, text string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}
}

func isAllowedCommand(command string) bool {
	for _, allowedCmd := range config.AllowedCmds {
		if strings.HasPrefix(command, allowedCmd) {
			return true
		}
	}
	return false
}

# tshbot

`tshbot` is a lightweight tool that enables you to execute shell commands on Linux systems via Telegram.

Initially it was developed with the idea of conveniently running commands on devices like Raspberry Pi remotely, eliminating the need to use SSH. Designed for ease of use, `tshbot` does not require complex setups such as VPNs or NAT traversal, simplifying your remote management experience.

### Features

- **Secure execution**: Send shell commands securely through your Telegram

- **No complex setup**: Forget about VPNs or NAT traversal - `tshbot` simplifies remote access

- **Shortcuts for commands**: Set up shortcuts to run frequently used commands quickly and easily

- **Flexibility**: You have the option to run arbitrary shell commands. This can be enabled or disabled based on your preference for security


# Configuration

The configuration file (`tshbot.config`) should contain the following fields:

- `bot_log_file`: Path to the log file
- `bash_cmd`: Path to the bash executable
- `tg_bot_token`: Telegram bot token
- `tg_bot_chat_id`: Telegram chat ID for the bot
- `allowed_cmds`: A map of command shortcuts to the actual shell commands
- `help_message`: Message displayed when the /help command is issued


# Acknowledgments

This project was inspired by [fnzv/trsh-go](https://github.com/fnzv/trsh-go).

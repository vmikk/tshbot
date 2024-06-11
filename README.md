# tshbot

`tshbot` is a lightweight tool that enables you to execute shell commands on Linux systems via Telegram.

Initially it was developed with the idea of conveniently running commands on devices like Raspberry Pi remotely, eliminating the need to use SSH. Designed for ease of use, `tshbot` does not require complex setups such as VPNs or NAT traversal, simplifying your remote management experience.

### Features

- **Secure execution**: Send shell commands securely through your Telegram

- **No complex setup**: Forget about VPNs or NAT traversal - `tshbot` simplifies remote access

- **Shortcuts for commands**: Set up shortcuts to run frequently used commands quickly and easily

- **Flexibility**: You have the option to run arbitrary shell commands. This can be enabled or disabled based on your preference for security

# Installation

Create a configuration file at `~/.config/tshbot/tshbot.config` with the following content:

``` yaml
bot_log_file: "/path/to/your/logfile.log"
bash_cmd: "/bin/bash"
tg_bot_token: "YOUR_TELEGRAM_BOT_TOKEN"
tg_bot_chat_id: "YOUR_TELEGRAM_CHAT_ID"
allowed_cmds:
  ls: 'ls -la'
  cmd: 'your commands here'
  uptime: 'uptime'
  shell: ''
help_message: "Use /commands to see available commands."
```
Replace `/path/to/your/logfile.log`, `YOUR_TELEGRAM_BOT_TOKEN`, and `YOUR_TELEGRAM_CHAT_ID` with your actual paths and credentials.

### Security considerations

>[!CAUTION]
> Allowing execution of arbitrary shell commands can be potentially dangerous. 
> It may lead to unauthorized access, system compromise, or data loss. 
> To mitigate these risks, if the `shell` shortcut is missing in the `allowed_cmds` configuration, only the white-listed commands specified in `allowed_cmds` are allowed.

# Configuration

The configuration file (`tshbot.config`) should contain the following fields:

- `bot_log_file`: Path to the log file
- `bash_cmd`: Path to the bash executable
- `tg_bot_token`: Telegram bot token
- `tg_bot_chat_id`: Telegram chat ID for the bot
- `allowed_cmds`: A map of command shortcuts to the actual shell commands
- `help_message`: Message displayed when the /help command is issued

# Usage

## Commands

- `/help`: Displays the help message
- `/commands`: Lists all available commands
- `/used_defined_shortcuts`: Executes a predefined command based on the shortcut defined in the configuration (e.g., `/cmd`)
- `/shell` [command]: Executes a shell command directly (if allowed in the configuration)

## Examples

- Sending `/ls` via Telegram will execute `ls -la` on the server.
- Sending `/shell uname -a` will execute `uname -a` on the server.


# Acknowledgments

This project was inspired by [fnzv/trsh-go](https://github.com/fnzv/trsh-go).

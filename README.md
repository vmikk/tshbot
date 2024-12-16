# tshbot

`tshbot` is a lightweight tool that enables you to execute shell commands on Linux systems via Telegram.

Initially it was developed with the idea of conveniently running commands on devices like Raspberry Pi remotely, eliminating the need to use SSH. Designed for ease of use, `tshbot` does not require complex setups such as VPNs or NAT traversal, simplifying your remote management experience.

<p align="middle">
  <img src="assets/tshbot_telegram.webp"  width=28% height=28% title="tshbot - Telegram screenshot"/>
</p>

### Features

- **Secure execution**: Send shell commands securely through your Telegram

- **No complex setup**: Forget about VPNs or NAT traversal - `tshbot` simplifies remote access

- **Shortcuts for commands**: Set up shortcuts to run frequently used commands quickly and easily

- **Flexibility**: You have the option to run arbitrary shell commands. This can be enabled or disabled based on your preference for security

# Installation

## 0. Set up a Telegram bot 

### Create a bot

1. **Start a chat with BotFather:** Open Telegram and search for "@BotFather". Start a chat with this bot, which is the official bot creation tool from Telegram;
2. **Create a new bot:** Send the command `/newbot` to BotFather;
3. **Set a name for your bot:** Follow the prompts and provide a name for your bot. This will be the public name that users see;
4. **Set a username for your bot:** Next, you'll need to choose a unique username for your bot (this username must contain "bot").

After completing these steps, BotFather will provide you with a **token**. This is your bot’s authentication token, which you’ll use to send and receive messages via the Telegram API.

### Obtain chat ID

1. **Start a chat with your bot:** Find your bot on Telegram using the username you set and start a conversation;
2. **Paste the following URL into your browser**, replacing `<YourBotToken>` with the token you received from BotFather:

   ```
   https://api.telegram.org/bot<YourBotToken>/getUpdates
   ```

3. **Find your chat ID:** Look for a number labeled `"chat":{"id":` in the response.


## 1. Download the pre-compiled binary

You can download the appropriate pre-compiled binary for your platform from the [available releases](https://github.com/vmikk/tshbot/releases).

##### E.g., general Linux release (64-bit `amd64`):

```sh
mkdir -p ~/bin
wget \
  -O ~/bin/tshbot \
  https://github.com/vmikk/tshbot/releases/download/0.2/tshbot-linux-amd64
chmod +x ~/bin/tshbot
```

##### Or release for Raspberry Pi Zero 2 W (`arm`):

```sh
mkdir -p ~/bin
wget \
  -O ~/bin/tshbot \
  https://github.com/vmikk/tshbot/releases/download/0.2/tshbot-linux-arm
chmod +x ~/bin/tshbot
```


## 2. Create a configuration file

Create a configuration file at `~/.config/tshbot/tshbot.config` with the following content:

``` yaml
bot_log_file: "/path/to/your/logfile.log"
bash_cmd: "/bin/bash"
tg_bot_token: "YOUR_TELEGRAM_BOT_TOKEN"
tg_bot_chat_id: "YOUR_TELEGRAM_CHAT_ID"
command_timeout: 300
allowed_cmds:
  pingg: 'ping -c 3 8.8.8.8'
  uptime: 'uptime'
  runscript: 'bash ~/bin/myscript.sh'
  shell: ""
help_message: "Use /commands to see available commands."
```
Replace `/path/to/your/logfile.log`, `YOUR_TELEGRAM_BOT_TOKEN`, and `YOUR_TELEGRAM_CHAT_ID` with your actual paths and credentials.  

In the `allowed_cmds` section of the config, you can configure shortcuts for various shell commands that you wish to use frequently. This section is structured as a dictionary where each entry consists of two parts:

- **Shortcut Name:** The key on the left side of the colon (`:`) represents the name of the shortcut. This is a unique identifier you will use to refer to the command;
- **Command:** The value on the right side of the colon is the actual command that will be executed in the shell when the shortcut is used.

For example, the shortcut `pingg` will execute the command `ping -c 3 8.8.8.8`.

Parameter `command_timeout` sets the timeout for old commands (in seconds). If the command is older than the timeout, it will be ignored. For example, setting `command_timeout: 300` will ignore commands older than 5 minutes.

### Security considerations

>[!CAUTION]
> Allowing execution of arbitrary shell commands can be potentially dangerous. 
> It may lead to unauthorized access, system compromise, or data loss. 
> To mitigate these risks, you may remove the `shell` shortcut from the `allowed_cmds` configuration. 
> This will **allow only the white-listed commands** specified in `allowed_cmds`.

## 3. Run the tool

``` sh
~/bin/tshbot
```

Now you can send command via Telegram.

# Configuration

The configuration file (`tshbot.config`) should contain the following fields:

- `bot_log_file`: Path to the log file
- `bash_cmd`: Path to the bash executable
- `tg_bot_token`: Telegram bot token
- `tg_bot_chat_id`: Telegram chat ID for the bot
- `command_timeout`: Timeout for old commands (in seconds)
- `allowed_cmds`: A map of command shortcuts to the actual shell commands
- `help_message`: Message displayed when the /help command is issued

# Usage

## Commands

- `/help`: Displays the help message
- `/commands`: Lists all available commands
- used-defined shortcuts (specified in `allowed_cmds`): Execute a predefined command based on the shortcut defined in the configuration (e.g., `/cmd`)
- `/shell` [command]: Executes a shell command directly (if allowed in the configuration)

## Examples

Using the example configuration file shown above, you can send the following commands via Telegram:

- Sending `/uptime` via Telegram will execute `uptime` on the server.
- Sending `/shell uname -a` will execute `uname -a` on the server.

## Logging

The tool logs its activities and command executions to the specified log file. 
Make sure the bot has write permissions to the log file.

# Auto-starting `tshbot` (on Linux-based systems with `systemd`)

To ensure `tshbot` starts automatically on your system, you can create a `systemd` service file. 
Follow these steps to set it up:

### 1. Create the service file

``` sh
sudo nano /etc/systemd/system/tshbot.service
```

Add the following content to the file:
``` ini
[Unit]
Description=tshbot service
After=network.target

[Service]
Type=simple
User=pi
ExecStart=/home/pi/bin/tshbot
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

- `User=pi`: Change `pi` to the user that should run the service if different
- `ExecStart=/home/pi/bin/tshbot`: If different, change to the actual path of the `tshbot` binary

### 2. Reload `systemd` and enable the service

```sh
sudo systemctl daemon-reload
```

Enable the `tshbot` service to start on boot
``` sh
sudo systemctl enable tshbot.service
```

Start the `tshbot` service
``` sh
sudo systemctl start tshbot.service
```

### 3. Verify the service

To check the status of the `tshbot` service and to ensure it is running correctly, use:
``` sh
sudo systemctl status tshbot.service
```

If everything is set up correctly, you should see an active (running) status.

## Managing the Service

To start the service:
```sh
sudo systemctl start tshbot.service
```

To stop the service:
```sh
sudo systemctl stop tshbot.service
```

To restart the service:
```sh
sudo systemctl restart tshbot.service
```

To check the service logs:
```sh
sudo journalctl -u tshbot.service
```

# Building the binary

To manually build the `tshbot` binary for different architectures, including ARM for Raspberry Pi, follow these instructions:

Prerequisites

- Go programming language installed on your system. You can download it from the official [Go website](https://go.dev/doc/install)
- Source code of `tshbot` cloned from the repository

## Building for local system

To build the binary for your local system, simply run the following command in the project directory:

``` sh
go build tshbot.go
```

This will generate a binary named `tshbot` in the current directory.

## Building for ARM (e.g., Raspberry Pi)

To build a binary for ARM architecture use the following command:

``` sh
env GOOS=linux GOARCH=arm go build -o tshbot tshbot.go
```

This command sets the `GOOS` and `GOARCH` environment variables to target the ARM architecture on Linux systems.

## Building for ARM64

If you need to build for ARM64 architecture, use this command instead:

``` sh
env GOOS=linux GOARCH=arm64 go build -o tshbot-linux-arm64 tshbot.go
```


# Acknowledgments

This project was inspired by [fnzv/trsh-go](https://github.com/fnzv/trsh-go).

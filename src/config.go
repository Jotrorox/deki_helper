package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml"
)

type Config struct {
	DB_PATH      string
	BOT_ID       string
	BOT_TOKEN    string
	CHANNEL      string
	USER_MENTION bool
	CMD_ADD_USER []string
}

func readConfigFile(filename string) *Config {
	data, err := os.ReadFile(filename)
	hpe(err)
	config := &Config{}
	err = toml.Unmarshal(data, config)
	hpe(err)
	return config
}

func createConfigFile(filename string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Database path (commands.db): ")
	db_path, _ := reader.ReadString('\n')
	db_path = strings.Replace(db_path, "\r\n", "", -1)
	db_path = strings.Replace(db_path, "\n", "", -1)

	fmt.Print("Your Bot Name: ")
	bot_id, _ := reader.ReadString('\n')
	bot_id = strings.Replace(bot_id, "\r\n", "", -1)
	bot_id = strings.Replace(bot_id, "\n", "", -1)

	fmt.Print("Your OAuth Token: ")
	bot_token, _ := reader.ReadString('\n')
	bot_token = "oauth:" + bot_token
	bot_token = strings.Replace(bot_token, "\r\n", "", -1)
	bot_token = strings.Replace(bot_token, "\n", "", -1)

	fmt.Print("The Channel: ")
	channel, _ := reader.ReadString('\n')
	channel = strings.Replace(channel, "\r\n", "", -1)
	channel = strings.Replace(channel, "\n", "", -1)

	fmt.Print("A User able to add commands: ")
	cmd_au, _ := reader.ReadString('\n')
	cmd_au = strings.Replace(cmd_au, "\r\n", "", -1)
	cmd_au = strings.Replace(cmd_au, "\n", "", -1)

	fmt.Print("Mention the User after a command (y/n): ")
	mu1, _ := reader.ReadString('\n')
	mu1 = strings.Replace(mu1, "\r\n", "", -1)
	mu1 = strings.Replace(mu1, "\n", "", -1)
	mu2 := false
	if mu1 == "y" {
		mu2 = true
	}

	var cmd_au_arr []string
	cmd_au_arr = append(cmd_au_arr, cmd_au)
	cmd_au_arr = append(cmd_au_arr, channel)

	config := &Config{
		DB_PATH:      db_path,
		BOT_ID:       bot_id,
		BOT_TOKEN:    bot_token,
		CHANNEL:      channel,
		USER_MENTION: mu2,
		CMD_ADD_USER: cmd_au_arr,
	}
	data, err := toml.Marshal(config)
	hpe(err)
	err = os.WriteFile(filename, data, 0644)
	hpe(err)
}

package main

import (
	"bufio"
	"fmt"
	"os"

	twitch "github.com/gempir/go-twitch-irc/v4"
	toml "github.com/pelletier/go-toml"
)

type Config struct {
	DB_PATH   string
	BOT_ID    string
	BOT_TOKEN string
	CHANNEL   string
}

func readConfigFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = toml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func createConfigFile(filename string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Database path (commands.db): ")
	db_path, _ := reader.ReadString('\n')
	fmt.Print("Your Bot ID: ")
	bot_id, _ := reader.ReadString('\n')
	fmt.Print("Your OAuth Token: ")
	bot_token, _ := reader.ReadString('\n')
	fmt.Print("The Channel: ")
	channel, _ := reader.ReadString('\n')
	config := &Config{
		DB_PATH:   db_path,
		BOT_ID:    bot_id,
		BOT_TOKEN: bot_token,
		CHANNEL:   channel,
	}
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	cfg_filename := "config.toml"

	_, err := os.Stat(cfg_filename)
	if os.IsNotExist(err) {
		err = createConfigFile(cfg_filename)
		if err != nil {
			fmt.Println("Failed to create config file:", err)
			return
		}
	}

	cfg, _ := readConfigFile(cfg_filename)

	client := twitch.NewClient(cfg.BOT_ID, cfg.BOT_TOKEN)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.Message == "!help" {
			client.Say(message.Channel, "@"+message.User.Name+" Right now there are only 4 commands: '!dc', 'insta', '!game' and '!lurk'. If you want more please contact deki!")
		}
		if message.Message == "!dc" {
			client.Say(message.Channel, "@"+message.User.Name+" the link to my Discord Server: https://discord.gg/Av9awsZz6K")
		}
		if message.Message == "!insta" {
			client.Say(message.Channel, "@"+message.User.Name+" the link to my Instagram: https://www.instagram.com/dekisenpaitm")
		}
		if message.Message == "!game" {
			client.Say(message.Channel, "@"+message.User.Name+" the link to the game I'm developing: https://dekisenpaitm.itch.io/project-soul")
		}
		if message.Message == "!lurk" {
			client.Say(message.Channel, "@"+message.User.Name+" is now lurking. Enjoy your time and lean back!")
		}
	})

	client.Join(cfg.CHANNEL)

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}

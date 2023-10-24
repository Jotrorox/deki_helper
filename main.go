package main

import (
	"bufio"
	"fmt"
	"os"

	twitch "github.com/gempir/go-twitch-irc/v4"
	"github.com/glebarez/sqlite"
	toml "github.com/pelletier/go-toml"
	"gorm.io/gorm"
)

type Config struct {
	DB_PATH      string
	BOT_ID       string
	BOT_TOKEN    string
	CHANNEL      string
	USER_MENTION bool
}

type Command struct {
	ID       uint
	Trigger  string
	Response string
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
	fmt.Print("Your Bot Name: ")
	bot_id, _ := reader.ReadString('\n')
	fmt.Print("Your OAuth Token: ")
	bot_token, _ := reader.ReadString('\n')
	bot_token = "oauth:" + bot_token
	fmt.Print("The Channel: ")
	channel, _ := reader.ReadString('\n')
	fmt.Print("Mention the User after a command (y/n): ")
	mu1, _ := reader.ReadString('\n')
	mu2 := false
	if mu1 == "y" {
		mu2 = true
	}
	config := &Config{
		DB_PATH:      db_path,
		BOT_ID:       bot_id,
		BOT_TOKEN:    bot_token,
		CHANNEL:      channel,
		USER_MENTION: mu2,
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

func connectToSQLite(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DB_PATH), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createEntryTable(db *gorm.DB) error {
	err := db.AutoMigrate(&Command{})
	if err != nil {
		return err
	}
	return nil
}

func addEntry(db *gorm.DB, cmd Command) error {
	entry := cmd
	result := db.Create(&entry)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func queryEntries(db *gorm.DB) ([]Command, error) {
	var entries []Command
	result := db.Find(&entries)
	if result.Error != nil {
		return nil, result.Error
	}
	return entries, nil
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

	db, err := connectToSQLite(cfg)
	if err != nil {
		panic(err)
	}

	err = createEntryTable(db)
	if err != nil {
		panic(err)
	}

	/*
		test_cmd := Command{
			Trigger:  "!help",
			Response: "Help!",
		}
		err = addEntry(db, test_cmd)
		if err != nil {
			panic(err)
		}
	*/

	client := twitch.NewClient(cfg.BOT_ID, cfg.BOT_TOKEN)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		commands, err := queryEntries(db)
		if err != nil {
			panic(err)
		}

		for _, command := range commands {
			if command.Trigger == message.Message {
				if cfg.USER_MENTION {
					client.Say(message.Channel, "@"+message.User.Name+"! "+command.Response)
				} else {
					client.Say(message.Channel, command.Response)
				}
			}
		}
	})

	client.Join(cfg.CHANNEL)

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}

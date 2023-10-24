package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	CMD_ADD_USER []string
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

func UpdateRow(db *gorm.DB, keyword string, newData Command) error {
	return db.Model(&Command{}).Where("trigger = ?", keyword).Updates(newData).Error
}

func RemoveRow(db *gorm.DB, keyword string) error {
	return db.Where("trigger = ?", keyword).Delete(&Command{}).Error
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

	client := twitch.NewClient(cfg.BOT_ID, cfg.BOT_TOKEN)

	cmd_iter := 0
	tmp_command := Command{}

	update_cmd_iter := 0
	update_tmp_command := Command{}
	update_cmd_key := ""

	delete_cmd_iter := 0
	delete_cmd_key := ""

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		for _, uname := range cfg.CMD_ADD_USER {
			if message.User.Name == uname {
				if cmd_iter == 0 && message.Message == "!add_cmd" {
					client.Say(message.Channel, "@"+message.User.Name+" What should the Trigger for you new command be?")
					cmd_iter = 1
				} else if cmd_iter == 1 {
					tmp_command.Trigger = message.Message
					client.Say(message.Channel, "@"+message.User.Name+" What should the Response for you new command be?")
					cmd_iter = 2
				} else if cmd_iter == 2 {
					tmp_command.Response = message.Message
					client.Say(message.Channel, "@"+message.User.Name+" Review everything, should this command be added? (yes/no)")
					cmd_iter = 3
				} else if cmd_iter == 3 {
					if message.Message == "yes" {
						addEntry(db, tmp_command)
						client.Say(message.Channel, "@"+message.User.Name+" The command has been added")
						cmd_iter = 0
					} else {
						client.Say(message.Channel, "@"+message.User.Name+" The command adding proces has been aborted")
						cmd_iter = 0
					}
				}

				if update_cmd_iter == 0 && message.Message == "!update_cmd" {
					client.Say(message.Channel, "@"+message.User.Name+" Which command do you wanna change?")
					update_cmd_iter = 1
				} else if update_cmd_iter == 1 {
					update_cmd_key = message.Message
					client.Say(message.Channel, "@"+message.User.Name+" What should the new Trigger for your command be?")
					update_cmd_iter = 2
				} else if update_cmd_iter == 2 {
					update_tmp_command.Trigger = message.Message
					client.Say(message.Channel, "@"+message.User.Name+" What should the new Response for your command be?")
					update_cmd_iter = 3
				} else if update_cmd_iter == 3 {
					update_tmp_command.Response = message.Message
					client.Say(message.Channel, "@"+message.User.Name+" Review everything, should this command be updated? (yes/no)")
					update_cmd_iter = 4
				} else if update_cmd_iter == 4 {
					if message.Message == "yes" {
						err = UpdateRow(db, update_cmd_key, update_tmp_command)
						if err != nil {
							client.Say(message.Channel, "@"+message.User.Name+" Something went wrong updating the command")
							update_cmd_iter = 0
						} else {
							client.Say(message.Channel, "@"+message.User.Name+" The command has been updated")
							update_cmd_iter = 0
						}
					} else {
						client.Say(message.Channel, "@"+message.User.Name+" The command updating process has been aborted")
						update_cmd_iter = 0
					}
				}

				if delete_cmd_iter == 0 && message.Message == "!delete_cmd" {
					client.Say(message.Channel, "@"+message.User.Name+" Which command do you wanna delete?")
					delete_cmd_iter = 1
				} else if delete_cmd_iter == 1 {
					delete_cmd_key = message.Message
					client.Say(message.Channel, "@"+message.User.Name+" Review everything, should this command be deleted? (yes/no)")
					delete_cmd_iter = 2
				} else if delete_cmd_iter == 2 {
					if message.Message == "yes" {
						err = RemoveRow(db, delete_cmd_key)
						if err != nil {
							client.Say(message.Channel, "@"+message.User.Name+" Something went wrong while deleting the command")
							delete_cmd_iter = 0
						} else {
							client.Say(message.Channel, "@"+message.User.Name+" The command has been deleted")
							delete_cmd_iter = 0
						}
					} else {
						client.Say(message.Channel, "@"+message.User.Name+" The command deleting process has been aborted")
						delete_cmd_iter = 0
					}
				}
			}
		}

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

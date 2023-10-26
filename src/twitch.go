package main

import (
	twitch "github.com/gempir/go-twitch-irc/v4"
	"gorm.io/gorm"
)

// handleTwitchClient initializes and handles Twitch client operations
func handleTwitchClient(cfg *Config, db *gorm.DB) {
	client := twitch.NewClient(cfg.BOT_ID, cfg.BOT_TOKEN)

	cmd_iter := 0
	tmp_command := Command{}

	update_cmd_iter := 0
	update_tmp_command := Command{}
	update_cmd_key := ""

	delete_cmd_iter := 0
	delete_cmd_key := ""

	// Handle private messages
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
						err := UpdateRow(db, update_cmd_key, update_tmp_command)
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
						err := RemoveRow(db, delete_cmd_key)
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

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

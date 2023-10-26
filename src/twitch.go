package main

import (
	"fmt"

	twitch "github.com/gempir/go-twitch-irc/v4"
	"gorm.io/gorm"
)

func processCommand(client *twitch.Client, db *gorm.DB, message twitch.PrivateMessage, cmd_iter *int, tmp_command *Command, command string, responses []string) {
	if *cmd_iter < 3 {
		tmp_command.Trigger = message.Message
		client.Say(message.Channel, fmt.Sprintf("@%s %s", message.User.Name, responses[*cmd_iter]))
		*cmd_iter++
	} else {
		if message.Message == "yes" {
			addEntry(db, *tmp_command)
		}
		client.Say(message.Channel, fmt.Sprintf("@%s %s", message.User.Name, responses[*cmd_iter]))
		*cmd_iter = 0
	}
}

// handleTwitchClient initializes and handles Twitch client operations
func handleTwitchClient(cfg *Config, db *gorm.DB) {
	client := twitch.NewClient(cfg.BOT_ID, cfg.BOT_TOKEN)

	cmd_iter := 0
	tmp_command := Command{}

	update_cmd_iter := 0
	update_tmp_command := Command{}

	delete_cmd_iter := 0
	delete_tmp_command := Command{}

	// Handle private messages
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		for _, uname := range cfg.CMD_ADD_USER {
			if message.User.Name == uname {
				if cmd_iter == 0 && message.Message == "!add_cmd" {
					responses := []string{"What should the Trigger for your new command be?", "What should the Response for your new command be?", "Review everything, should this command be added? (yes/no)", "The command has been added", "The command adding process has been aborted"}
					processCommand(client, db, message, &cmd_iter, &tmp_command, "!add_cmd", responses)
				} else if update_cmd_iter == 0 && message.Message == "!update_cmd" {
					responses := []string{"Which command do you want to change?", "What should the new Trigger for your command be?", "What should the new Response for your command be?", "Review everything, should this command be updated? (yes/no)", "The command has been updated", "The command updating process has been aborted"}
					processCommand(client, db, message, &update_cmd_iter, &update_tmp_command, "!update_cmd", responses)
				} else if delete_cmd_iter == 0 && message.Message == "!delete_cmd" {
					responses := []string{"Which command do you want to delete?", "Review everything, should this command be deleted? (yes/no)", "The command has been deleted", "The command deleting process has been aborted"}
					processCommand(client, db, message, &delete_cmd_iter, &delete_tmp_command, "!delete_cmd", responses)
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

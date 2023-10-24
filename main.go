package main

import (
	"os"

	twitch "github.com/gempir/go-twitch-irc/v4"
	dotenv "github.com/joho/godotenv"
)

func main() {
	err := dotenv.Load()
	if err != nil {
		panic(err)
	}

	client := twitch.NewClient(os.Getenv("BOT_NAME"), os.Getenv("OAUTH_TOKEN"))

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

	client.Join("deki_senpai_tm")

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}

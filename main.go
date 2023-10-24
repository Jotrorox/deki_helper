package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	twitch "github.com/gempir/go-twitch-irc/v4"
	_ "github.com/mattn/go-sqlite3"
	toml "github.com/pelletier/go-toml"
)

type Config struct {
	DB_PATH   string
	BOT_ID    string
	BOT_TOKEN string
	CHANNEL   string
}

type Command struct {
	ID       int
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

func setupDB(cfg *Config) *sql.DB {
	db, err := sql.Open("sqlite3", cfg.DB_PATH)
	if err != nil {
		fmt.Println("Failed to open database:", err)
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		trigger TEXT,
		response TEXT
	)`)
	if err != nil {
		fmt.Println("Failed to create table:", err)
		panic(err)
	}
	return db
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

	db := setupDB(cfg)

	_, err = db.Exec("INSERT INTO entries (trigger, response) VALUES (?, ?)", "!help", "Help!")
	if err != nil {
		fmt.Println("Failed to add entry:", err)
		return
	}

	client := twitch.NewClient(cfg.BOT_ID, cfg.BOT_TOKEN)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		rows, err := db.Query("SELECT * FROM commands")
		if err != nil {
			fmt.Println("Failed to query entries")
			panic(err)
		}
		defer rows.Close()

		var commands []Command
		for rows.Next() {
			var command Command
			if err := rows.Scan(&command.ID, &command.Trigger, &command.Response); err != nil {
				fmt.Println("Failed to scan entry:", err)
				return
			}
			commands = append(commands, command)
		}

		for _, command := range commands {
			if message.Message == command.Trigger {
				client.Say(message.Channel, command.Response)
			}
		}
	})

	client.Join(cfg.CHANNEL)

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}

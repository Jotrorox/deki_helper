// main.go
package main

import (
	"fmt"
	"os"
)

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

	// Initialize and handle Twitch client operations
	handleTwitchClient(cfg, db)
}

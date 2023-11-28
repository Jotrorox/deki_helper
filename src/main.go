// main.go
package main

import (
	"os"
)

func main() {
	cfg_filename := "config.toml"

	_, err := os.Stat(cfg_filename)
	if os.IsNotExist(err) {
		createConfigFile(cfg_filename)
	}

	cfg := readConfigFile(cfg_filename)

	db, err := connectToSQLite(cfg)
	hpe(err)

	createEntryTable(db)

	// Initialize and handle Twitch client operations
	handleTwitchClient(cfg, db)
}

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Settings *Config

type Config struct {
	Get struct {
		// Settings for daemon
		Address string
		Port    uint
	}

	Post struct {
		// Settings for daemon
		Address string
		Port    uint
	}

	Admin struct {
		// Settings for daemon
		Address string
		Port    uint
	}

	Directories struct {
		// Storage directory for images
		ImageDir     string
		ThumbnailDir string
	}

	// sites for CORS
	CORS struct {
		Sites []string
	}

	Database struct {
		// Database connection settings
		User           string
		Password       string
		Proto          string
		Host           string
		Database       string
		MaxIdle        int
		MaxConnections int
	}

	Redis struct {
		// Redis address and max pool connections
		Protocol       string
		Address        string
		MaxIdle        int
		MaxConnections int
	}

	// HMAC secret for bcrypt
	Session struct {
		Secret string
	}
}

func Print() {

	// Marshal the structs into JSON
	output, err := json.MarshalIndent(Settings, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", output)

}

func init() {
	file, err := os.Open("/etc/pram/pram.conf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Settings = &Config{}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&Settings)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

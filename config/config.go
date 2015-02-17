package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Settings *Config

type Config struct {
	General struct {
		// Settings for daemon
		Address string
		Port    uint

		// Default name if none specified
		DefaultName string

		// Storage directory for images
		ImageDir     string
		ThumbnailDir string
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

	Akismet struct {
		// Akismet settings
		Key  string
		Host string
	}

	Antispam struct {
		// Antispam Key from Prim
		AntispamKey string

		// Antispam cookie
		CookieName  string
		CookieValue string
	}

	StopForumSpam struct {
		// Stop Forum Spam settings
		Confidence float64
	}

	Limits struct {
		// Image settings
		ImageMinWidth  int
		ImageMinHeight int
		ImageMaxWidth  int
		ImageMaxHeight int
		ImageMaxSize   int
		WebmMaxLength  int

		// Max posts in a thread
		PostsMax uint

		// Lengths for posting
		CommentMaxLength int
		CommentMinLength int
		TitleMaxLength   int
		TitleMinLength   int
		NameMaxLength    int
		NameMinLength    int
		TagMaxLength     int
		TagMinLength     int

		// Max thumbnail sizes
		ThumbnailMaxWidth  int
		ThumbnailMaxHeight int

		// Max request parameter input size
		ParamMaxSize uint
	}
}

func Print() {

	// Marshal the structs into JSON
	output, err := json.MarshalIndent(Settings, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", output)

}

func init() {
	file, err := os.Open("/etc/pram/post.conf")
	if err != nil {
		panic(err)
	}

	Settings = &Config{}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&Settings)
	if err != nil {
		panic(err)
	}

}

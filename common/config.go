package common

import (
	"encoding/xml"
	"log"
	"os"
)

var config Config

type Config struct {
	Username        string `xml:"username"`
	Password        string `xml:"password"`
	DatabaseAddress string `xml:"databaseAddress"`
	DatabaseName    string `xml:"databaseName"`
	AssetsPath      string `xml:"assetsPath"`
	ImagesPath      string `xml:"imagesPath"`
}

func CheckError(err error) {
	if err != nil {
		log.Fatalf("Nintendo Channel file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

func GetConfig() Config {
	if config.DatabaseName != "" {
		return config
	}

	data, err := os.ReadFile("config.xml")
	CheckError(err)

	err = xml.Unmarshal(data, &config)
	return config
}

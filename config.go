package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	MailHost      string
	MailPort      string
	MailFrom      string
	MailFromText  string
	MailPass      string
	MailTo        []string
	MailSubject   string
	Url           string
	JobTime       string
	JobTimeHour   int
	JobTimeMinute int
	JobTimeSecond int
}

func readConfig() *Config {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}
	timeParts := strings.Split(config.JobTime, ":")
	config.JobTimeHour, err = strconv.Atoi(timeParts[0])
	if err != nil {
		fmt.Println("error:", err)
	}
	config.JobTimeMinute, err = strconv.Atoi(timeParts[1])
	if err != nil {
		fmt.Println("error:", err)
	}

	if len(timeParts) == 3 {
		config.JobTimeSecond, err = strconv.Atoi(timeParts[2])
		if err != nil {
			fmt.Println("error:", err)
		}
	} else {
		config.JobTimeSecond = 0
	}

	return &config
}

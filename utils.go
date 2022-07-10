package main

import (
	"log"
	"os"
	"os/user"
)

func getCurrentDir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "/"
	}
	return path
}

func getCurrentHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
		return ""
	}
	return hostname
}

func getCurrentUsername() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
		return "nobody"
	}
	return currentUser.Username
}

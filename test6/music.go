package main

import (
	"log"
	"os/exec"
)

func fail() {
	log.Println("playing fail music")
	cmd := exec.Command("sudo", "python", "fail-music.py")
	err := cmd.Run()
	if err != nil {
		log.Println("Errored", err)
	}
}

func succeed() {
	log.Println("playing success music")
	cmd := exec.Command("sudo", "python", "success-music.py")
	err := cmd.Run()
	if err != nil {
		log.Println("Errored", err)
	}
}

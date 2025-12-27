package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/viitorags/orvit/internal/bot"
)

func main() {
	bot.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

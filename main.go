package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-telegram/bot"
)

var gitHash = "development"
var p2hHash = ""

var telegramFromId int64
var telegramDir string = "data"
var telegramSaveJson bool = true
var telegramDelete bool = false

func main() {
	showHelp := flag.Bool("help", false, "show help")
	showVersion := flag.Bool("version", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Println("Version:")
		fmt.Println("     git", gitHash)
		if p2hHash != "" {
			fmt.Println("     p2h", p2hHash)
		}
		os.Exit(0)
	}

	if *showHelp {
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println()
		return
	}

	delete := os.Getenv("TELEGRAM_DELETE")
	if delete == "1" {
		telegramDelete = true
	}

	saveJsonFlag := os.Getenv("TELEGRAM_JSON")
	if saveJsonFlag == "0" {
		telegramSaveJson = false
	}

	dir := os.Getenv("TELEGRAM_DIR")
	if dir != "" {
		telegramDir = dir
	}

	token := os.Getenv("TELEGRAM_APITOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_APITOKEN not found")
	}

	var err error
	telegramFromId, err = strconv.ParseInt(os.Getenv("TELEGRAM_FROMID"), 10, 64)
	if err != nil {
		log.Fatal("TELEGRAM_FROMID not found")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

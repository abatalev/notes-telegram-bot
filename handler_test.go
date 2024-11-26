package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"testing"

	"github.com/go-telegram/bot/models"
)

func TestHandler(t *testing.T) {
	telegramDir = t.TempDir()
	telegramFromId = 111
	telegramSaveJson = false
	variants := []struct {
		jsonName string
	}{
		{jsonName: "handler/empty"},
		{jsonName: "handler/me"},
		{jsonName: "handler/other"},
	}
	for _, variant := range variants {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		message := unmarshalJson(variant.jsonName, t)
		handler(ctx, nil, &models.Update{Message: &message})
	}
}

func TestDownload(t *testing.T) {
	telegramDir = t.TempDir()
	if err := downloadFile(filepath.Join(telegramDir, "index.html"), "https://ya.ru/"); err != nil {
		t.Fatal(err)
	}
}

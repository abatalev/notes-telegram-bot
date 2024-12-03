package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	x := update.Message
	if x == nil || x.From == nil {
		log.Println("skip empty message")
		return
	}

	if x.From.ID != telegramFromId {
		log.Println("skip >>", x.ID, x.From.ID, x.From.FirstName, x.From.LastName)
		return
	}

	// TODO mediagroup
	log.Println("work >>", x.ID, x.From.ID, x.From.FirstName, x.From.LastName)
	x1 := strconv.Itoa(x.ID)
	if telegramSaveJson {
		bb, _ := json.Marshal(x)
		if err := os.WriteFile(filepath.Join(telegramDir, telegramPrefix+x1+".json"), bb, 0666); err != nil {
			panic(err)
		}
	}

	t := InitTemplate(func(id string) string {
		fileUrl, ok := getFileUrlById(ctx, b, id)
		if ok {
			panic("???")
		}
		parsedFileUrl, err := url.Parse(fileUrl)
		if err != nil {
			panic(err)
		}
		shortFileName := telegramPrefix + x1 + "_" + filepath.Base(parsedFileUrl.Path)
		path := filepath.Join(telegramDir, shortFileName)
		log.Println("Download", path, fileUrl)
		if err := downloadFile(path, fileUrl); err != nil {
			panic(err)
		}

		return shortFileName
	})
	if err := os.WriteFile(filepath.Join(telegramDir, telegramPrefix+x1+".md"), Parse(t, *x), 0666); err != nil {
		panic(err)
	}

	if telegramDelete {
		// TODO delete parsed message
		panic("unimplemented")
	}
}

func getFileUrlById(ctx context.Context, b *bot.Bot, fileID string) (string, bool) {
	s, err := b.GetFile(ctx, &bot.GetFileParams{FileID: fileID})
	if err != nil {
		log.Println("ERROR", err)
		return "", true
	}
	return b.FileDownloadLink(s), false
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

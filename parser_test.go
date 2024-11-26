package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/require"
)

func TestSplit80(t *testing.T) {
	variants := []struct {
		nn     int
		input  string
		output string
	}{
		{nn: 5, input: "xxx", output: "xxx"},
		{nn: 5, input: "xxxxx", output: "xxxxx"},
		{nn: 5, input: "xx xx xx", output: "xx xx\nxx"},
		{nn: 5, input: "xx xxxx xx", output: "xx\nxxxx\nxx"},
		{nn: 5, input: "xxxxxxx", output: "xxxxxxx"},
		{nn: 5, input: "xxxxxxx x", output: "xxxxxxx\nx"},

		{nn: 5, input: "яяя", output: "яяя"},
		{nn: 5, input: "яяяяя", output: "яяяяя"},
		{nn: 5, input: "яя яя яя", output: "яя яя\nяя"},
		{nn: 5, input: "яя яяяя яя", output: "яя\nяяяя\nяя"},
		{nn: 5, input: "яяяяяяя", output: "яяяяяяя"},
		{nn: 5, input: "яяяяяяя я", output: "яяяяяяя\nя"},
	}
	for _, variant := range variants {
		require.Equal(t, variant.output, splitNn(variant.nn, variant.input))
	}
}

func TestParseExamples(t *testing.T) {
	telegramSaveJson = false
	telegramDir = t.TempDir()
	nn := []string{
		"photo_n_emoji", "text", "photo_n_caption", "video_n_caption",
		"voice_n_caption", "text_code", "text_pre", "text_pre_lang",
		"text_url", "text_reply",
		"caption_animation_document",
	}
	tmpl := InitTemplate(func(id string) string {
		return "banner.jpg"
	})
	for _, n := range nn {
		require.Equal(t, string(readFile("parser/"+n, "md", t)), string(Parse(tmpl, unmarshalJson("parser/"+n, t))), n)
	}
}

func unmarshalJson(nn string, t *testing.T) models.Message {
	var message models.Message
	if err := json.Unmarshal(readFile(nn, "json", t), &message); err != nil {
		t.Fatal(err)
	}
	return message
}

func readFile(nn, ext string, t *testing.T) []byte {
	b, err := os.ReadFile("examples/" + nn + "." + ext)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

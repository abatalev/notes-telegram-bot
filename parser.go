package main

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/go-telegram/bot/models"
)

type GlobalConfig struct {
	Prefix string
}

func Parse(t *template.Template, message models.Message) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := t.Execute(buff, struct {
		Global  GlobalConfig
		Message models.Message
		Content string
	}{
		Global:  GlobalConfig{Prefix: telegramPrefix},
		Message: message,
		Content: parseContent(message),
	})
	if err != nil {
		return []byte{}, nil
	}
	return buff.Bytes(), nil
}

func getSizedPhoto(photos []models.PhotoSize) string {
	for _, photo := range slices.Backward(photos) {
		if photo.FileSize > 300000 { // TODO TELEGRAM_MAXFILESIZE
			continue
		}
		return photo.FileID
	}
	return photos[0].FileID
}

func InitTemplate(getFileName func(string) string) *template.Template {
	t, err := template.New("note").
		Funcs(template.FuncMap{
			"unixDateTime": func(date int) string {
				return time.Unix(int64(date), 0).UTC().Format(time.RFC3339)
			},
			"getFileName":   getFileName,
			"getSizedPhoto": getSizedPhoto,
			"title":         getTitle,
		}).
		Parse(LoadTemplate())
	if err != nil {
		panic(err)
	}
	return t
}

func getTitle(content string) string {
	s := content
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "```", ".\n")
	s = strings.ReplaceAll(s, "\n", " ")
	linkRegexp := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	s = linkRegexp.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, " @", "")
	idx := strings.Index(s, ". ")
	if idx >= 0 {
		if idx > telegramTitleSize {
			log.Println("!!!!", idx, telegramTitleSize)
			return fineTune(s[0:telegramTitleSize])
		}
		return fineTune(s[0:idx])
	}
	return fineTune(s)
}

func fineTune(s string) string {
	s = strings.TrimSpace(s)
	if strings.Contains(s, ":") {
		s = "\"" + s + "\""
	}
	return s
}

func LoadTemplate() string {
	f, err := os.ReadFile(telegramTemplate)
	if err != nil {
		panic(err)
	}
	return string(f)
}

func parseContent(message models.Message) string {
	if len(message.Caption) > 0 {
		return processing(message.Caption, message.CaptionEntities)
	}
	if len(message.Text) > 0 {
		return processing(message.Text, message.Entities)
	}
	log.Println("!!! Empty")
	return ""
}

func processing(text0 string, entities []models.MessageEntity) string {
	text := []rune(text0)
	points := createMarkers(entities, text)
	o := make([]rune, 0)
	prev := 0
	for _, n := range createIndexes(points) {
		o = append(o, text[prev:n]...)
		o = append(o, []rune(points[n])...)
		prev = n
	}
	o = append(o, text[prev:]...)
	return postProcessing(o)
}

func postProcessing(o []rune) string {
	so := string(o)
	ss := ""
	lines := strings.Split(so, "\n")
	for _, line := range lines {
		if ss == "" {
			ss = splitNn(80, strings.TrimRight(line, " "))
		} else {
			ss = ss + "\n" + splitNn(80, strings.TrimRight(line, " "))
		}
	}
	return ss
}

func splitNn(n int, s string) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}

	if r[n] == rune(' ') {
		return strings.TrimRight(string(r[:n]), " ") + "\n" + splitNn(n, strings.TrimLeft(string(r[n:]), " "))
	}

	idx := strings.LastIndex(string(r[:n]), " ")
	if idx != -1 {
		return strings.TrimRight(s[:idx], " ") + "\n" + splitNn(n, strings.TrimLeft(s[idx:], " "))
	}

	idx = strings.Index(s, " ")
	if idx == -1 {
		return s
	}
	return strings.TrimRight(s[:idx], " ") + "\n" + splitNn(n, strings.TrimLeft(s[idx:], " "))
}

func createIndexes(points map[int]string) []int {
	indexes := make([]int, 0)
	for k := range points {
		indexes = append(indexes, k)
	}
	slices.Sort(indexes)
	return indexes
}

func createOffsets(caption []rune) map[int]int {
	offsets := make(map[int]int, 0)
	idx := 0
	for n, ch := range caption {
		offsets[idx] = n
		if utf8.RuneLen(ch) > 3 {
			idx++
		}
		idx++
	}
	return offsets
}

func createMarkers(entries []models.MessageEntity, text []rune) map[int]string {
	realOffsets := createOffsets(text)
	prefixes := make(map[int]string, 0)
	suffixes := make(map[int]string, 0)
	for _, entry := range entries {
		realOffset := realOffsets[entry.Offset]
		switch entry.Type {
		// TODO “mention” (@username),
		// TODO “hashtag” (#hashtag),
		// TODO “cashtag” ($USD),
		// TODO “bot_command” (/start@jobs_bot),
		// TODO “email” (do-not-reply@telegram.org),
		// TODO “phone_number” (+1-212-555-0123),
		// TODO “text_mention” (for users without usernames)
		// TODO "custom_emoji":
		case models.MessageEntityTypeBold:
			flag := false
			for _, r := range text[realOffset : realOffset+entry.Length] {
				if r != rune(' ') {
					flag = true
				}
			}
			if flag {
				prefixes[realOffset] = prefixes[realOffset] + "**"
				suffixes[realOffset+entry.Length] = "**" + suffixes[realOffset+entry.Length]
			}
		case models.MessageEntityTypeStrikethrough:
			prefixes[realOffset] = prefixes[realOffset] + "~~"
			suffixes[realOffset+entry.Length] = "~~" + suffixes[realOffset+entry.Length]
		case models.MessageEntityTypeUnderline:
			prefixes[realOffset] = prefixes[realOffset] + "<u>"
			suffixes[realOffset+entry.Length] = "</u>" + suffixes[realOffset+entry.Length]
		case models.MessageEntityTypeItalic:
			prefixes[realOffset] = prefixes[realOffset] + "_"
			suffixes[realOffset+entry.Length] = "_" + suffixes[realOffset+entry.Length]
		case models.MessageEntityTypeTextLink:
			prefixes[realOffset] = prefixes[realOffset] + "["
			suffixes[realOffset+entry.Length] = "](" + entry.URL + ")" + suffixes[realOffset+entry.Length]
		case "blockquote":
			prefixes[realOffset] = prefixes[realOffset] + "> "
			for n, r := range text[realOffset : realOffset+entry.Length] {
				if r == rune('\n') {
					prefixes[realOffset+n+1] = prefixes[realOffset+n+1] + "> "
				}
			}
		case models.MessageEntityTypePre:
			prefixes[realOffset] = prefixes[realOffset] + "```" + entry.Language + "\n"
			suffixes[realOffset+entry.Length] = "\n```" + suffixes[realOffset+entry.Length]
		case models.MessageEntityTypeCode:
			prefixes[realOffset] = prefixes[realOffset] + "`"
			suffixes[realOffset+entry.Length] = "`" + suffixes[realOffset+entry.Length]
		case models.MessageEntityTypeURL:
			prefixes[realOffset] = prefixes[realOffset] + "[@]("
			suffixes[realOffset+entry.Length] = ")" + suffixes[realOffset+entry.Length]
		}
	}
	return createPoints(prefixes, suffixes)
}

func createPoints(prefixes map[int]string, suffixes map[int]string) map[int]string {
	points := prefixes
	for k, v := range suffixes {
		points[k] = v + points[k]
	}
	return points
}

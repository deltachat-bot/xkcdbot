package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/deltachat-bot/deltabot-cli-go/botcli"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	"github.com/nishanths/go-xkcd/v2"
	"github.com/spf13/cobra"
)

var cli = botcli.New("xkcdbot")

func reportError(err error, chat *deltachat.Chat) {
	cli.Logger.Error(err)
	chat.SendText(fmt.Sprintf("Error: %v", err))
}

func sendHelp(chat *deltachat.Chat) {
	text := "I am a bot that allows to retrieve comics from https://xkcd.com/\n\n"
	text += "Send me a message with the number of a XKCD comic and I will send you the image. For example, send me: 1254\n\n"
	text += "Available Commands:\n\n/random - get a random comic.\n\n/latest - get latest comic.\n\n/get [number] - get comic corresponding to the given number. Example: /get 1254\n\n/help - send this help message"
	chat.SendText(text)
}

func sendComic(chat *deltachat.Chat, numb int) {
	client := xkcd.NewClient()
	var comic xkcd.Comic
	var err error
	if numb <= 0 {
		comic, err = client.Latest(context.Background())
	} else {
		comic, err = client.Get(context.Background(), numb)
	}
	if err != nil {
		reportError(err, chat)
		return
	}

	data, content_type, err := client.Image(context.Background(), comic.Number)
	if err != nil {
		reportError(err, chat)
		return
	}

	var ext string
	switch content_type {
	case "image/png":
		ext = "png"
	default:
		ext = "jpg"
	}
	file, err := os.CreateTemp("", "image-*."+ext)
	if err != nil {
		return
	}
	filename := file.Name()
	defer os.Remove(filename)

	_, err = file.ReadFrom(data)
	if err != nil {
		return
	}

	chat.SendMsg(deltachat.MsgData{
		Text: fmt.Sprintf("#%v - %v", comic.Number, comic.Title),
		Html: comic.Alt,
		File: filename,
	})
}

func onNewMsg(bot *deltachat.Bot, message *deltachat.Message) {
	msg, err := message.Snapshot()
	if err != nil || msg.IsInfo {
		return
	}
	chat := &deltachat.Chat{bot.Account, msg.ChatId}
	args := strings.Split(msg.Text, " ")
	switch args[0] {
	case "/help":
		sendHelp(chat)
	case "/random":
		comic, err := xkcd.NewClient().Latest(context.Background())
		if err != nil {
			reportError(err, chat)
			return
		}
		numb := rand.Intn(comic.Number) + 1
		sendComic(chat, numb)
	case "/latest":
		sendComic(chat, 0)
	default:
		var text string
		if args[0] == "/get" && len(args) == 2 {
			text = args[1]
		} else {
			chatInfo, err := chat.BasicSnapshot()
			if err != nil || chatInfo.ChatType != deltachat.CHAT_TYPE_SINGLE {
				return
			}
			text = args[0]
		}
		numb, err := strconv.Atoi(text)
		if err != nil {
			sendHelp(chat)
			return
		}
		sendComic(chat, numb)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	cli.OnBotInit(func(bot *deltachat.Bot, cmd *cobra.Command, args []string) {
		name, _ := bot.GetConfig("displayname")
		if name == "" {
			bot.SetConfig("displayname", "XKCD Bot")
			bot.SetConfig("selfstatus", "I am a bot that allows to get XKCD comics, send me /help for more info")
		}
		bot.OnNewMsg(func(message *deltachat.Message) { onNewMsg(bot, message) })
	})
	cli.OnBotStart(func(bot *deltachat.Bot, cmd *cobra.Command, args []string) {
		addr, _ := bot.GetConfig("addr")
		cli.Logger.Infof("Listening at: %v", addr)
	})
	cli.Start()
}

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
	"github.com/deltachat/deltachat-rpc-client-go/deltachat/option"
	"github.com/nishanths/go-xkcd/v2"
	"github.com/spf13/cobra"
)

var cli = botcli.New("xkcdbot")

func reportError(rpc *deltachat.Rpc, accId deltachat.AccountId, chatId deltachat.ChatId, err error) {
	cli.Logger.Error(err)
	rpc.MiscSendTextMessage(accId, chatId, fmt.Sprintf("Error: %v", err))
}

func sendHelp(rpc *deltachat.Rpc, accId deltachat.AccountId, chatId deltachat.ChatId) {
	text := "I am a bot that allows to retrieve comics from https://xkcd.com/\n\n"
	text += "Send me a message with the number of a XKCD comic and I will send you the image. For example, send me: 1254\n\n"
	text += "Available Commands:\n\n/random - get a random comic.\n\n/latest - get latest comic.\n\n/get [number] - get comic corresponding to the given number. Example: /get 1254\n\n/help - send this help message"
	rpc.MiscSendTextMessage(accId, chatId, text)
}

func sendComic(rpc *deltachat.Rpc, accId deltachat.AccountId, chatId deltachat.ChatId, numb int) {
	client := xkcd.NewClient()
	var comic xkcd.Comic
	var err error
	if numb <= 0 {
		comic, err = client.Latest(context.Background())
	} else {
		comic, err = client.Get(context.Background(), numb)
	}
	if err != nil {
		reportError(rpc, accId, chatId, err)
		return
	}

	data, content_type, err := client.Image(context.Background(), comic.Number)
	if err != nil {
		reportError(rpc, accId, chatId, err)
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

	rpc.SendMsg(accId, chatId, deltachat.MsgData{
		Text: fmt.Sprintf("#%v - %v", comic.Number, comic.Title),
		Html: comic.Alt,
		File: filename,
	})
}

func onNewMsg(bot *deltachat.Bot, accId deltachat.AccountId, msgId deltachat.MsgId) {
	msg, err := bot.Rpc.GetMessage(accId, msgId)
	if err != nil || msg.IsInfo || msg.IsBot || msg.FromId <= deltachat.ContactLastSpecial {
		return
	}

	args := strings.Split(msg.Text, " ")
	switch args[0] {
	case "/help":
		sendHelp(bot.Rpc, accId, msg.ChatId)
	case "/random":
		comic, err := xkcd.NewClient().Latest(context.Background())
		if err != nil {
			reportError(bot.Rpc, accId, msg.ChatId, err)
			return
		}
		numb := rand.Intn(comic.Number) + 1
		sendComic(bot.Rpc, accId, msg.ChatId, numb)
	case "/latest":
		sendComic(bot.Rpc, accId, msg.ChatId, 0)
	default:
		var text string
		if args[0] == "/get" && len(args) == 2 {
			text = args[1]
		} else {
			chatInfo, err := bot.Rpc.GetBasicChatInfo(accId, msg.ChatId)
			if err != nil || chatInfo.ChatType != deltachat.ChatSingle {
				return
			}
			text = args[0]
		}
		numb, err := strconv.Atoi(text)
		if err != nil {
			sendHelp(bot.Rpc, accId, msg.ChatId)
			return
		}
		sendComic(bot.Rpc, accId, msg.ChatId, numb)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	cli.OnBotInit(func(cli *botcli.BotCli, bot *deltachat.Bot, cmd *cobra.Command, args []string) {
		accIds, _ := bot.Rpc.GetAllAccountIds()
		for _, accId := range accIds {
			name, _ := bot.Rpc.GetConfig(accId, "displayname")
			if name.UnwrapOr("") == "" {
				bot.Rpc.SetConfig(accId, "displayname", option.Some("XKCD Bot"))
				bot.Rpc.SetConfig(accId, "selfstatus", option.Some("I am a bot that allows to get XKCD comics, send me /help for more info"))
			}
		}
		bot.OnNewMsg(onNewMsg)
	})
	cli.Start()
}

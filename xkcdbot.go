package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/deltachat-bot/deltabot-cli-go/botcli"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	"github.com/nishanths/go-xkcd/v2"
	"github.com/spf13/cobra"
)

var bot *deltachat.Bot
var cli = botcli.New("xkcdbot")

func logEvent(event *deltachat.Event) {
	switch event.Type {
	case deltachat.EVENT_INFO:
		cli.Logger.Info().Msg(event.Msg)
	case deltachat.EVENT_WARNING:
		cli.Logger.Warn().Msg(event.Msg)
	case deltachat.EVENT_ERROR:
		cli.Logger.Error().Msg(event.Msg)
	}
}

func fail(err error, chat *deltachat.Chat) {
	cli.Logger.Error().Err(err)
	chat.SendText(fmt.Sprintf("Error: %v", err))
}

func sendHelp(chat *deltachat.Chat) {
	chat.SendText("Send me a message with the number of a XKCD comic and I will send you the image. For example, send me: 599")
}

func onNewMsg(message *deltachat.Message) {
	msg, err := message.Snapshot()
	if err != nil || msg.IsInfo {
		return
	}
	chat := &deltachat.Chat{bot.Account, msg.ChatId}
	chatInfo, err := chat.BasicSnapshot()
	if err != nil || chatInfo.ChatType != deltachat.CHAT_TYPE_SINGLE {
		return
	}

	client := xkcd.NewClient()
	numb, err := strconv.Atoi(msg.Text)
	if err != nil {
		sendHelp(chat)
		return
	}
	comic, err := client.Get(context.Background(), numb)
	if err != nil {
		fail(err, chat)
		return
	}

	data, content_type, err := client.Image(context.Background(), numb)
	if err != nil {
		fail(err, chat)
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
		fail(err, chat)
		return
	}
	filename := file.Name()
	defer os.Remove(filename)

	_, err = file.ReadFrom(data)
	if err != nil {
		fail(err, chat)
		return
	}

	chat.SendMsg(deltachat.MsgData{
		Text: fmt.Sprintf("#%v - %v", comic.Number, comic.Title),
		Html: comic.Alt,
		File: filename,
	})
}

func main() {
	cli.OnBotInit(func(newBot *deltachat.Bot, cmd *cobra.Command, args []string) {
		bot = newBot
		bot.On(deltachat.EVENT_INFO, logEvent)
		bot.On(deltachat.EVENT_WARNING, logEvent)
		bot.On(deltachat.EVENT_ERROR, logEvent)
		bot.OnNewMsg(onNewMsg)
	})
	cli.OnBotStart(func(bot *deltachat.Bot, cmd *cobra.Command, args []string) {
		addr, _ := bot.GetConfig("addr")
		cli.Logger.Info().Msgf("Listening at: %v", addr)
	})
	cli.Start()
}

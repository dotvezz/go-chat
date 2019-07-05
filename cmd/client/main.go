package main

import (
	"fmt"
	"github.com/dotvezz/gochat/chat"
	"github.com/dotvezz/gochat/chat/connection"
	"github.com/marcusolsson/tui-go"
	"log"
	"net"
	"os"
)

func main() {
	nc, err := net.Dial("tcp", "localhost:1026")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	conn := connection.New(nc)

	myName := new(string)
	input, chatView, ui := initTUI()
	handleSubmit := initSubmitHandler(input, myName, conn)
	input.OnSubmit(handleSubmit)

	listen := initMessageListener(conn, chatView, ui)
	go listen()

	err = ui.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func initSubmitHandler(input *tui.Entry, myName *string, conn chat.Connection) func(e *tui.Entry) {
	return func(e *tui.Entry) {
		m := chat.Message{
			From: *myName,
			To:   "everyone",
			Body: e.Text(),
		}
		conn.Send(m)
		input.SetText("")
	}
}

func initMessageListener(conn chat.Connection, chatView *tui.Box, ui tui.UI) func() {
	return func() {
		for {
			m, err := conn.Receive()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				break
			}

			chatView.Append(tui.NewLabel(fmt.Sprintf(" %s: %s", m.From, m.Body)))
			ui.Repaint() // Repaint to show the new stuff in the chat view
			ui.Repaint() // but twice because it seems to make tui freak out less often
		}
	}
}

func initTUI() (*tui.Entry, *tui.Box, tui.UI) {
	input := tui.NewEntry()
	chatView := tui.NewVBox()
	historyScroll := tui.NewScrollArea(chatView)
	historyScroll.SetAutoscrollToBottom(true)
	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	chatBox := tui.NewVBox(historyBox, inputBox)
	chatBox.SetSizePolicy(tui.Expanding, tui.Expanding)

	ui, err := tui.New(chatBox)
	if err != nil {
		log.Fatal(err)
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	return input, chatView, ui
}

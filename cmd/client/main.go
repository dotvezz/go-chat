package main

import (
	"fmt"
	"github.com/dotvezz/go-chat/chat"
	"github.com/dotvezz/go-chat/chat/connection"
	"github.com/dotvezz/go-chat/chat/domain/client/config"
	"github.com/marcusolsson/tui-go"
	"log"
	"net"
	"time"
)

func main() {
	// Load the configuration
	conf := config.New()
	chat.LoadConfig(&conf)

	// Connect to the server
	nc, err := net.Dial("tcp", fmt.Sprintf("%s%s", conf.Host, conf.Port))
	if err != nil {
		log.Fatal(err)
	}

	// Set the connection
	conn := connection.NewClient(nc)

	// Initialize the TUI
	input, chatView, ui := initTUI()

	// Build the submit handler, injecting the Input, name pointer, and Connection
	handleSubmit := initSubmitHandler(input, conn)
	// Set the submit handler
	input.OnSubmit(handleSubmit)

	// Initialize the listner, injecting the Connection and TUI elements
	listen := initMessageListener(conn, chatView, ui)
	// Listen
	go listen()

	err = ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Initialize the handler function for the input's Submit action
func initSubmitHandler(input *tui.Entry, conn chat.Connection) func(e *tui.Entry) {
	return func(e *tui.Entry) {
		m := chat.Message{
			From: conn.UserName(),
			Body: e.Text(),
		}

		conn.Send(m)
		input.SetText("")
	}
}

// Initialize the listener function
func initMessageListener(conn chat.Connection, chatView *tui.Box, ui tui.UI) func() {
	return func() {
		for {
			m, err := conn.Receive()
			if err != nil {
				log.Fatal("disconnected:", err)
			}

			chatView.Append(tui.NewLabel(
				fmt.Sprintf(
					"%s | %s: %s",
					time.Unix(m.TimeStamp, 0).Format("2006-01-02 3:04:05"),
					m.From,
					m.Body,
				),
			))
			ui.Repaint() // Repaint to show the new stuff in the chat view
			ui.Repaint() // but twice because it seems to make tui freak out less often
		}
	}
}

// Initialize the TUI
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
	ui.SetKeybinding("Ctrl+c", func() { ui.Quit() })

	return input, chatView, ui
}

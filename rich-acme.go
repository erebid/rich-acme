package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"9fans.net/go/acme"
	presence "github.com/hugolgst/rich-go/client"
)

const appID = "721524490957357096"

func main() {
	start := time.Now()
	l, err := acme.Log()
	if err != nil {
		log.Fatal(err)
	}
	presence.Login(appID)
	focusCh := make(chan int)
	go updateFocus(focusCh, l)
	ticker := time.NewTicker(15 * time.Second)
	var focused int
	for {
		select {
		case <-ticker.C:
			err := updatePresence(start, focused)
			if err != nil {
				fmt.Println(err)
			}
		case f := <-focusCh:
			focused = f
		}
	}
}

func updateFocus(c chan int, l *acme.LogReader) {
	for {
		event, err := l.Read()
		if err != nil {
			log.Fatal(err)
		}
		if event.Op == "focus" {
			c <- event.ID
		}
	}
}

func updatePresence(start time.Time, id int) error {
	var filename string
	windows, err := acme.Windows()
	if err != nil {
		return err
	}
	for _, win := range windows {
		if win.ID == id {
			_, filename = filepath.Split(win.Name)
			break
		}
	}
	if filename == "" {
		return nil
	}
	err = presence.SetActivity(presence.Activity{
		Details:    "Editing " + filename,
		State:      strconv.Itoa(len(windows)) + " files open",
		LargeImage: "glenda",
		Timestamps: &presence.Timestamps{
			Start: &start,
		},
	})
	return err
}
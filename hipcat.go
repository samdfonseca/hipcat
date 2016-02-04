package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/tbruyelle/hipchat-go/hipchat"
)

type HipCat struct {
	api      *hipchat.Client
	queue    *StreamQ
	shutdown chan os.Signal
	roomName string
	roomId   string
}

func newHipCat(authToken string, roomId string, roomName string) (*HipCat, error) {
	hc := HipCat{
		api:      hipchat.NewClient(authToken),
		queue:    newStreamQ(),
		shutdown: make(chan os.Signal, 1),
		roomName: roomName,
		roomId:   roomId,
	}
	if hc.roomId == "" {
		err := hc.lookupRoomId()
		if err != nil {
			return nil, err
		}
	}
	signal.Notify(hc.shutdown, os.Interrupt)
	return &hc, nil
}

func (hc *HipCat) trap() {
	sigcount := 0
	for sig := range hc.shutdown {
		if sigcount > 0 {
			exitErr(fmt.Errorf("aborted"))
		}
		output(fmt.Sprintf("got signal: %s", sig.String()))
		output("press ctrl+c again to exit immediately")
		sigcount++
		go hc.exit()
	}
}

func (hc *HipCat) exit() {
	for {
		if hc.queue.isEmpty() {
			os.Exit(0)
		} else {
			output("flushing remaining messages to HipChat...")
			time.Sleep(3 * time.Second)
		}
	}
}

//Lookup id for room by name
func (hc *HipCat) lookupRoomId() error {
	api := hc.api
	room, _, err := api.Room.Get(hc.roomName)
	if err == nil {
		hc.roomId = fmt.Sprint(room.ID)
		return nil
	}
	fmt.Println(err)
	return fmt.Errorf("Unable to find room: %s", hc.roomName)
}

func (hc *HipCat) addToStreamQ(lines chan string) {
	for line := range lines {
		hc.queue.add(line)
	}
	hc.exit()
}

//TODO: handle messages with length exceeding maximum for HipChat
func (hc *HipCat) processStreamQ(noop bool, plain bool) {
	if !(hc.queue.isEmpty()) {
		msglines := hc.queue.flush()
		hc.postMsg(msglines, plain, noop)
	}
	time.Sleep(1 * time.Second)
	hc.processStreamQ(noop, plain)
}

func (hc *HipCat) postMsg(msglines []string, plain bool, noop bool) {
	msgFmtStr := "<code>%s</code>"
	messageFmt := "html"
	if plain {
		msgFmtStr = "%s"
		messageFmt = "text"
	}
	msg := fmt.Sprintf(msgFmtStr, strings.Join(msglines, "<br>"))
	if noop {
		output(fmt.Sprintf("skipped posting of %s message lines to %s", strconv.Itoa(len(msglines)), hc.roomName))
		output(msg)
		return
	}
	notifReq := &hipchat.NotificationRequest{Message: msg, MessageFormat: messageFmt}
	resp, err := hc.api.Room.Notification(hc.roomId, notifReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during room notification %q\n", err)
		fmt.Fprintf(os.Stderr, "Server returns %+v\n", resp)
		if resp.ContentLength >= 0 {
			body := make([]byte, resp.ContentLength)
			_, err := resp.Body.Read(body)
			failOnError(err, "Unable to read response body", true)
			fmt.Fprintf(os.Stderr, "%s\n", body)
		}
		return
	}
	output(fmt.Sprintf("posted %s message lines to %s", strconv.Itoa(len(msglines)), hc.roomName))
}

func (hc *HipCat) postFile(filePath, fileName string, noop bool) {
	//default to timestamp for filename
	if fileName == "" {
		fileName = strconv.FormatInt(time.Now().Unix(), 10)
	}

	if noop {
		output(fmt.Sprintf("skipping upload of file %s to %s", fileName, hc.roomName))
		return
	}

	start := time.Now()
	shareFileReq := &hipchat.ShareFileRequest{Path: filePath, Message: "Shared file from hipcat", Filename: fileName}
	if hc.roomId != "" {
		resp, err := hc.api.Room.ShareFile(hc.roomId, shareFileReq)
		if err != nil {
			fmt.Printf("Error during room file share %q\n", err)
			fmt.Printf("Server returns %+v\n", resp)
			return
		}
	}
	duration := strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)
	output(fmt.Sprintf("file %s uploaded to %s (%ss)", fileName, hc.roomName, duration))
	os.Exit(0)
}

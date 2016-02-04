package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

var version = "dev-build"

func readIn(lines chan string, tee bool) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines <- scanner.Text()
		if tee {
			fmt.Println(scanner.Text())
		}
	}
	close(lines)
}

func writeTemp(lines chan string) string {
	tmp, err := ioutil.TempFile(os.TempDir(), "hipcat-")
	failOnError(err, "unable to create tmpfile", false)

	w := bufio.NewWriter(tmp)
	for line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()

	return tmp.Name()
}

func output(s string) {
	bold := color.New(color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", bold("hipcat"), s)
}

func failOnError(err error, msg string, appendErr bool) {
	if err != nil {
		if appendErr {
			exitErr(fmt.Errorf("%s: %s", msg, err))
		} else {
			exitErr(fmt.Errorf("%s", msg))
		}
	}
}

func exitErr(err error) {
	output(color.RedString(err.Error()))
	os.Exit(1)
}

func main() {
	app := cli.NewApp()
	app.Name = "hipcat"
	app.Usage = "redirect a file to hipchat"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "tee, t",
			Usage: "Print stdin to screen before posting",
		},
		cli.BoolFlag{
			Name:  "stream, s",
			Usage: "Stream messages to HipChat continuously instead of uploading a single snippet",
		},
		cli.BoolFlag{
			Name:  "plain, p",
			Usage: "Write messages as plain texts instead of code blocks",
		},
		cli.BoolFlag{
			Name:  "noop",
			Usage: "Skip posting file to HipChat. Useful for testing",
		},
		cli.StringFlag{
			Name:  "room, r",
			Usage: "HipChat room to post to",
		},
		cli.StringFlag{
			Name: "roomid, i",
			Usage: "HipChat room id to post to",
		},
		cli.StringFlag{
			Name:  "filename, n",
			Usage: "Filename for upload. Defaults to current timestamp",
		},
	}

	app.Action = func(c *cli.Context) {

		config := readConfig()
		fileName := c.String("filename")
		roomName := c.String("room")
		roomId := c.String("roomid")
		if roomName == "" && roomId == "" {
			if config.defaultRoomId == "" && config.defaultRoomName == "" {
				exitErr(fmt.Errorf("'room' flag is required if default_room_id is unset"))
			}
			roomName = config.defaultRoomName
			roomId = config.defaultRoomId
		}


		if !c.Bool("stream") && c.Bool("plain") {
			exitErr(fmt.Errorf("'plain' flag requires 'stream' mode!"))
		}

		hipcat, err := newHipCat(config.authToken, roomId, roomName)
		failOnError(err, "HipChat API Error", true)

		if len(c.Args()) > 0 {
			if c.Bool("stream") {
				output("filepath provided, ignoring stream option")
			}
			filePath := c.Args()[0]
			if fileName == "" {
				fileName = filepath.Base(filePath)
			}
			hipcat.postFile(filePath, fileName, c.Bool("noop"))
			os.Exit(0)
		}

		lines := make(chan string)
		go readIn(lines, c.Bool("tee"))

		if c.Bool("stream") {
			output("starting stream")
			go hipcat.addToStreamQ(lines)
			go hipcat.processStreamQ(c.Bool("noop"), c.Bool("plain"))
			go hipcat.trap()
			select {}
		} else {
			filePath := writeTemp(lines)
			defer os.Remove(filePath)
			hipcat.postFile(filePath, fileName, c.Bool("noop"))
			os.Exit(0)
		}
	}

	app.Run(os.Args)

}

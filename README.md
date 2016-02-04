# hipcat
HipCat is a simple commandline utility to post snippets to HipChat.


  <img width="500px" src="https://raw.githubusercontent.com/samdfonseca/hipcat/master/demo.gif" alt="hipcat"/>


## Quickstart

### Configuration

Generate a personal access token for HipCat: https://hipchat.com/account/api

Create a HipCat config file:
```bash
echo "auth_token = <your-hipchat-token>" >> ~/.hipcat
```

Set a default room to send to:
```bash
echo "default_room_name = Notification Room" >> ~/.hipcat
OR
echo "default_room_id = 1234567" >> ~/.hipcat
```
(Room ID is slightly faster since it saves one request to the HipChat API to lookup the room id)

## Usage
Pipe command output as a message or several messages to the default room:
```bash
$ tail -F -n0 /path/to/log | hipcat --tee --stream
hipcat file hello uploaded to Notification Room
```

Pipe command output as a message or several messages to some other room:
```bash
$ tail -F -n0 /path/to/log | hipcat --stream -r "Notification Room"
hipcat starting stream
hipcat posted 10 message lines to Notification Room
```

Post an existing file:
```bash
$ hipcat --room "Entire Company" /home/user/bot.png
hipcat file bot.png uploaded to general
```

Stream input continously as a formatted message, and print stdin back to stdout:
```bash
$ tail -F -n0 /path/to/log | hipcat --tee --stream
*hipcat* posted 5 message lines to general
*hipcat* posted 2 message lines to general
...
```

## Options

Option | Description
--- | ---
--tee, -t | Print stdin to screen before posting
--stream, -s | Stream messages to HipChat continuously instead of uploading a single snippet
--plain, -p | When streaming, write messages as plain text instead of code blocks
--noop | Skip posting file to HipChat. Useful for testing
--room, -r | HipChat channel, group, or user to post to
--filename, -n | Filename for upload. Defaults to given filename or current timestamp if reading from stdin.

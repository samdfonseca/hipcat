# Configuration Guide
HipChatcat may be configured via a simple or advanced configuration

## Configuration

#### Example ~/.hipcat Config
```bash
auth_token = thisisadefinitelyafaketoken
default_room = Department of Cool New Shiz
```
By default, all messages will be sent to the team1 general channel.

#### Example Usage

Post a file to default room:
```bash
hipcat /path/to/file.txt
```

Post a file to room by name:
```bash
hipcat -r "Notification Black Hole Room" /path/to/file.txt
```

Post a file to room by room id:
(slightly faster since it saves one request to the HipChat api to lookup the room id)
```bash
hipcat -i 1234567 /path/to/file.txt
```

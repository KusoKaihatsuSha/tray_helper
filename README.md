[![godoc](https://godoc.org/github.com/KusoKaihatsuSha/tray_helper?status.svg)](https://godoc.org/github.com/KusoKaihatsuSha/tray_helper) [![Go Report Card](https://goreportcard.com/badge/github.com/KusoKaihatsuSha/tray_helper)](https://goreportcard.com/report/github.com/KusoKaihatsuSha/tray_helper) [![go test](https://github.com/KusoKaihatsuSha/tray_helper/actions/workflows/test.yml/badge.svg)](https://github.com/KusoKaihatsuSha/tray_helper/actions/workflows/test.yml)

# Tray helper*

> App to help with daily routine.

`* The application uses a configuration file, which could be created on pressing the 'Settings' button in the GUI popup, or you can create this file manually, for example by looking inside the folder **__build** after execute **go generate windows_build.go**.`

### **Available actions**

`TARGET📌`
Focusing on the window by title. Click on the middle of screen for protected window

`CLICK TARGET📌`
Run exec file and wait for std output.

`EXEC WAIT OUTPUT TO CLIP🗒️`
Run exec file and wait for std output.

`EXEC🗒️`
Run exec file and do not wait.

`GEN♻️`
Generate random string with length equal to field value. If less than '0' will be generated difficult password.

`TEXT TO CLIPBOARD🖇️`
Copy text to clipboard

`OPEN URL🔖`
Open URL in default browser

`SUPER➕,CTRL➕,SHIFT➕,ALT➕,CTRL➕SHIFT➕,ALT➕SHIFT➕,CTRL➕ALT➕,CTRL➕ALT➕SHIFT➕`
Additional keys

`PASTE🔠`
Paste text at the current whatever place. Select Destination before that.

`SLEEP⌛`
Sleep and wait

`FILE`
Read file and write data to clipboard

`FILE LAST LINE`
Read last line of the file and write data to clipboard

### **Available flags**

> Run with custom settings file

`-config=filename.data` or `-c=filename.data`  

> Address of settings

`-a=localhost:8080`

### **Build**

`go build -ldflags "-s -w -H=windowsgui"`

OR

`go generate windows_build.go`

### Example **settings.data:**

```json
{
    "Generate 16 len PASS": {
        "actions": "GEN@-16",
        "timer": "30s",
        "repeat": 1,
        "silent": true
    },
    "Ping google.com and paste into open notepad": {
        "actions": "EXEC@notepad|EXECSTD@ping google.com|SLEEP@300ms|TARGET@Notepad|SLEEP@300ms|PRESS@CTRL+V|SLEEP@300ms",
        "timer": "",
        "repeat": 1,
        "silent": false
    }
}

```

Screenshots of settings:

<div style="width:50%">
<img src="/files/settings-0.png" ><br>
<img src="/files/settings-list.png" ><br>
<img src="/files/settings.png" ><br>
</div>

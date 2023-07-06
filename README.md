[![godoc](https://godoc.org/github.com/KusoKaihatsuSha/tray_helper?status.svg)](https://godoc.org/github.com/KusoKaihatsuSha/tray_helper) [![Go Report Card](https://goreportcard.com/badge/github.com/KusoKaihatsuSha/tray_helper)](https://goreportcard.com/report/github.com/KusoKaihatsuSha/tray_helper) [![go test](https://github.com/KusoKaihatsuSha/tray_helper/actions/workflows/test.yml/badge.svg)](https://github.com/KusoKaihatsuSha/tray_helper/actions/workflows/test.yml)


# Tray helper

Windows application to help with daily routine. 

When routine activities need to be run frequently, automation saves time.

You can compose simple chains of "actions" and execute them by clicking in the tray menu.

For example:

✔️ Copying some key from a file to the clipboard to avoid showing the location and the file itself to colleagues when presenting. 

✔️ Launching an application and pasting a "login" into it, if it's not saved automatically and getting tired of doing it by hand.

✔️ Open all personal/work sites in one click on the work computer so that nothing is stored in history.

✔️ Automatically switch the focus to the right window after copying the desired file and then automatically paste. This comes in handy for simple automation.


### **Сapability**

1) Simple chains of "actions"
Looks like 
> URL@https://go.dev/play/|URL@https://google.com|EXECSTD@ping -n 10 google.com

2) Set a timer, when it expires all "self-starting" (if app starts without cmd \k as an example) non OS-protected process will be killed.

3) Repeat chains of "actions".

4) Send notification when complete and timer ends


### **Actions**

`TARGET📌` - Focusing on the window by title. Click on the middle of screen for protected window

`CLICK TARGET📌` - Run exec file and wait for std output.

`EXEC WAIT OUTPUT TO CLIP🗒️` - Run exec file and wait for std output.

`EXEC🗒️` - Run exec file and do not wait.

`GEN♻️` - Generate random string with length equal to field value. If less than '0' will be generated difficult password.

`TEXT TO CLIPBOARD🖇️` - Copy text to clipboard

`OPEN URL🔖` - Open URL in default browser

`SUPER➕,CTRL➕,SHIFT➕,ALT➕,CTRL➕SHIFT➕,ALT➕SHIFT➕,CTRL➕ALT➕,CTRL➕ALT➕SHIFT➕` - Additional keys

`PASTE🔠` - Paste text at the current whatever place. Select Destination before that.

`SLEEP⌛` - Sleep and wait

`FILE` - Read file and write data to clipboard

`FILE LAST LINE` - Read last line of the file and write data to clipboard


### **Available flags on start binary**

`-config=filename.data` or `-c=filename.data` - Run with custom settings file

`-a=localhost:8080` - Address of web settings


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

<div style="width:40%">
<label>In tray:</label><br>
<img src="/files/settings-0.png" ><br>
<label>Tasks:</label><br>
<img src="/files/settings-list.png" ><br>
<label>Chains of "actions:"</label><br>
<img src="/files/settings.png" ><br>
</div>


`* The application uses a configuration file, which could be created on pressing the 'Settings' button in the GUI popup, or you can create this file manually (by looking example).`
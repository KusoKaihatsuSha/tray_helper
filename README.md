[![godoc](https://godoc.org/github.com/KusoKaihatsuSha/tray_helper?status.svg)](https://godoc.org/github.com/KusoKaihatsuSha/tray_helper) [![Go Report Card](https://goreportcard.com/badge/github.com/KusoKaihatsuSha/tray_helper)](https://goreportcard.com/report/github.com/KusoKaihatsuSha/tray_helper)

# Tray helper
App to help with daily routine.

`* The application uses a configuration file, which could be created on pressing the 'Settings' button in the GUI popup, or you can create this file manually, for example by looking inside the folder **__build** after execute **go generate windows_build.go**.`

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
<img src="/pictures/settings-0.png" ><br>
<img src="/pictures/settings-list.png" ><br>
<img src="/pictures/settings.png" ><br>
</div>

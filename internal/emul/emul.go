package emul

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/KusoKaihatsuSha/tray_helper/internal/clipwindow"
	"github.com/KusoKaihatsuSha/tray_helper/internal/config"
	"github.com/KusoKaihatsuSha/tray_helper/internal/helpers"
	"github.com/KusoKaihatsuSha/tray_helper/internal/keymap"
	"github.com/KusoKaihatsuSha/tray_helper/internal/notification"

	"github.com/micmonay/keybd_event"
)

var empty = ""

// CtxKey wrapper around context keys
type CtxKey string

// Command consist name and data of one command
type Command struct {
	Name string
	Data string
}

// ParseCommands parse all actions line
func ParseCommands(cmd string) []Command {
	var all []Command
	for _, val := range ParseCommand(cmd, '|') {
		command := ParseCommand(val, '@')
		var c Command
		c.Name = command[0]
		if len(command) > 1 {
			c.Data = command[1]
		}
		all = append(all, c)
	}
	return all
}

// ParseCommand parse one command into struct
func ParseCommand(cmd string, s rune) []string {
	return strings.FieldsFunc(cmd, func(r rune) bool {
		switch r {
		case s:
			return true
		default:
			return false
		}
	})
}

// parsePress - emulating press button
func parsePress(cmd string) (string, bool) {
	if cmd == empty {
		return fmt.Sprintf("warning: %v", "skip because empty"), true
	}

	newSet, err := keybd_event.NewKeyBonding()
	if err != nil {
		return fmt.Sprintf("error %v", err), false
	}

	localCommands := strings.FieldsFunc(cmd, func(r rune) bool {
		switch r {
		case '+':
			return true
		default:
			return false
		}
	})
	newSet.Clear()
	for _, v := range localCommands {
		switch v {
		case "CTRL":
			newSet.HasCTRL(true)
		case "ALT":
			newSet.HasALT(true)
		case "SHIFT":
			newSet.HasSHIFT(true)
		case "SUPER":
			newSet.HasSuper(true)
		default:
			if key, ok := keymap.KeyMap[v]; ok {
				newSet.AddKey(key)
			}
		}
	}
	if err := newSet.Launching(); err != nil {
		return fmt.Sprintf("error %v", err), false
	}
	return "", true
}

// parseClip - fill the clipboard with input data
func parseClip(cmd string) (string, bool) {
	if cmd == empty {
		clipwindow.Clear()
	}
	clipwindow.Set(cmd)
	return "", true // dummy
}

// parsePrint - fill the clipboard with input data and paste it
func parsePrint(cmd string) (string, bool) {
	if cmd == empty {
		clipwindow.Clear()
	}
	clipwindow.Set(cmd)
	return parsePress("CTRL+V")
}

// parseFileLastLine - fill the clipboard with input data (file body last line)
func parseFileLastLine(path string) (string, bool) {
	if path == empty {
		return fmt.Sprintf("warning: %v", "skip because empty"), true
	}

	if helpers.FileExist(path) {
		file, err := os.Open(path)
		if err != nil {
			helpers.ToLog(err)
			return fmt.Sprintf("error %v", err), false
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		line := ""
		for scanner.Scan() {
			line = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			helpers.ToLog(err)
			return fmt.Sprintf("error %v", err), false
		}
		clipwindow.Set(line)
	} else {
		return fmt.Sprintf("file %s not found", path), false
	}
	return "", true
}

// parseFile - fill the clipboard with input data (file body)
func parseFile(path string) (string, bool) {
	if path == empty {
		return fmt.Sprintf("warning: %v", "skip because empty"), true
	}

	if helpers.FileExist(path) {
		file, err := os.ReadFile(path)
		if err != nil {
			helpers.ToLog(err)
			return fmt.Sprintf("error %v", err), false
		}
		clipwindow.Set(string(file))
	} else {
		return fmt.Sprintf("file %s not found", path), false
	}
	return "", true
}

// parseTarget - focusing the window by title name
func parseTarget(ctx context.Context, target string) (string, bool) {
	if target == empty {
		return fmt.Sprintf("warning: %v", "skip because empty"), true
	}
	select {
	case <-clipwindow.FindTitle(target):
		ok := clipwindow.SetTopWindow(target)
		if !ok {
			return fmt.Sprintf("Target window title '%s' not found. Open app and try again.", target), false
		}
	case <-ctx.Done():
		// DO NOTHING
	case <-time.After(10 * time.Second):
		return fmt.Sprintf("Target window title '%s' not found. Open app and try again.", target), false
	}
	return "", true
}

// parseTargetClick - focusing the window by title name and click middle
func parseTargetClick(ctx context.Context, target string) (string, bool) {
	if target == empty {
		return fmt.Sprintf("warning: %v", "skip because empty"), true
	}
	select {
	case <-clipwindow.FindTitle(target):
		ok := clipwindow.SetTopWindowClick(target)
		if !ok {
			return fmt.Sprintf("Target window title '%s' not found. Open app and try again.", target), false
		}
	case <-ctx.Done():
		// DO NOTHING
	case <-time.After(10 * time.Second):
		return fmt.Sprintf("Target window title '%s' not found. Open app and try again.", target), false
	}
	return "", true
}

// parseGenerate - generate string with length
func parseGenerate(cmd string) (string, bool) {
	if v, err := strconv.Atoi(cmd); err == nil {
		if v < 0 {
			// for password
			clipwindow.Set(helpers.RandomStringPass(-v))
		} else {
			// for random string
			clipwindow.Set(helpers.RandomString(v))
		}
		return "", true
	}
	return fmt.Sprintf("error: %v", "check value after GEN@. Must be number."), false
}

// parseExecOutput - run exec with wait output
func parseExecOutput(ctx context.Context, cmd string) (string, bool) {
	if cmd == empty {
		return fmt.Sprintf("warning: %v", "skip because empty"), true
	}

	c := strings.SplitN(cmd, " ", 2)
	cargs := []string{}
	if len(c) > 1 {
		cargs = strings.Split(c[1], " ")
	}
	cmdCtx := exec.CommandContext(ctx, c[0], cargs...)
	if out, err := cmdCtx.CombinedOutput(); err == nil {
		t := time.Now().Add(10 * time.Second)
		for cmdCtx.Process == nil && time.Since(t) <= 0 {
			<-time.After(100 * time.Millisecond)
		}
		select {
		case <-clipwindow.FindPid(cmdCtx.Process.Pid):
			ok := clipwindow.SetTopWindow(cmdCtx.Process.Pid)
			if !ok {
				return fmt.Sprintf("error exec %v", err), false
			}
		case <-ctx.Done():
			// DO NOTHING
		case <-time.After(1 * time.Second):
			// WAIT
		}
		clipwindow.Set(string(out))
	} else {
		return fmt.Sprintf("error exec %v", err), false
	}
	return "", true
}

// parseExec - run exec
func parseExec(ctx context.Context, cmd string) (string, bool) {
	if cmd == empty {
		return fmt.Sprintf("warning: %v", "skip because empty"), true
	}

	c := strings.SplitN(cmd, " ", 2)
	cargs := []string{}
	if len(c) > 1 {
		cargs = strings.Split(c[1], " ")
	}
	if cmd := exec.CommandContext(ctx, c[0], cargs...); cmd.Start() == nil {
		t := time.Now().Add(10 * time.Second)
		for cmd.Process == nil && time.Since(t) <= 0 {
			<-time.After(100 * time.Millisecond)
		}
		select {
		case <-clipwindow.FindPid(cmd.Process.Pid):
			ok := clipwindow.SetTopWindow(cmd.Process.Pid)
			if !ok {
				return fmt.Sprintf("error exec %v", cmd), false
			}
		case <-ctx.Done():
			// DO NOTHING
		case <-time.After(1 * time.Second):
			// WAIT
		}
	} else {
		return fmt.Sprintf("error exec %v", cmd), false
	}
	return "", true
}

// Handle - handler the command
func Handle(ctx context.Context, cancel context.CancelFunc, texts string, repeat int) int {
	counter := 0
ex:
	for i := 0; i < repeat; i++ {
		for _, c := range ParseCommands(texts) {
			okResult := true
			errPoint := ""
			select {
			case <-ctx.Done():
				break ex
			default:
			}

			switch c.Name {
			case "CRYPT":
				// TODO:
			case "BASE64":
				// TODO:
			case "FIND":
				// TODO:
			case "URL":
				if c.Data != empty {
					helpers.OpenUrl(c.Data)
				}
			case "CLIP":
				errPoint, okResult = parseClip(c.Data)
			case "GEN":
				errPoint, okResult = parseGenerate(c.Data)
			case "EXECSTD":
				errPoint, okResult = parseExecOutput(ctx, c.Data)
			case "EXEC":
				errPoint, okResult = parseExec(ctx, c.Data)
			case "TARGET":
				errPoint, okResult = parseTarget(ctx, c.Data)
			case "CLICKTARGET":
				errPoint, okResult = parseTargetClick(ctx, c.Data)
			case "PRINT":
				errPoint, okResult = parsePrint(c.Data)
			case "PRESS":
				errPoint, okResult = parsePress(c.Data)
			case "SLEEP":
				if v, err := time.ParseDuration(c.Data); err == nil {
					time.Sleep(v)
				}
			case "FILE_READ_LAST_LINE_TO_CLIPBOARD":
				errPoint, okResult = parseFileLastLine(c.Data)
			case "FILE_READ_TO_CLIPBOARD":
				errPoint, okResult = parseFile(c.Data)
			case "":
				// DO NOTHING
			default:
				// DO NOTHING
			}
			if !okResult {
				notification.Toast(
					ctx,
					fmt.Sprintf("Error %s", config.Title),
					errPoint,
					false,
				)
				cancel()
			}
		}
		counter++
	}
	return counter
}

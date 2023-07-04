package config

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KusoKaihatsuSha/tray_helper/internal/helpers"
)

// General annotation.
const (
	DefTimePostfix = "s"
	HTTPString     = `http://`
	HTTPSString    = `https://`

	RealIPName  = "X-Real-IP"
	IcoNotif    = `icons/notif.ico`
	IcoApp      = `icons/icon.ico`
	IcoSettings = `icons/settings.ico`
	IcoClose    = `icons/close.ico`
	Title       = "Easy tray helper"
	ToolTip     = "<" + Title + "> make you the daily work some easy"

	Repeat         = "repeat"
	Silent         = "silent"
	WaitUntilClose = "timer"
	EmulateAction  = "actions"
)

type CtxKey string

var (
	this      Configuration
	onceFlags sync.Once = sync.Once{}
)

//go:embed icons/*
var EmbedFiles embed.FS

// Configuration consist settings, which filling on the start.
// valid list:
//
//	       "tmp_file"  - Crop prefix path and add tmp path for windows. For linux same string
//			  "file"      - Check exist and return correct /\
//			  "url"       - Check correct URL
//			  "bool"      - Check correct bool
//			  "timer"     - Check correct timer. If not included "s"... suffix
//	       "-", "none" - Nothing to do. Same string
type Configuration struct {
	PrefixURL string
	Tmp
	ConfigurationMain
}

// Tmp struct consist tmp information when reading/valid configs
type Tmp struct {
	tmp map[string]TagInfo
}

// ConfigurationMain - block consist the main settings
type ConfigurationMain struct {
	Config  string `json:"config" default:"settings.data" flag:"f,c,config" env:"CONFIG" text:"Config file" valid:"file"`
	Ico     string `json:"ico_file" default:"app.ico" flag:"ico" env:"ICO_FILE" text:"Filepath for ico notif" valid:"tmp_file"`
	Address string `json:"address" default:"ip:0" flag:"a" env:"ADDRESS" text:"Address and port" valid:"url"`
}

// TagInfo store tags Config struct
type TagInfo struct {
	ConfigFile func(map[string]any) TagInfo
	Valid      func() any
	DummyFlags func() TagInfo
	Env        func() TagInfo
	Flag       func() TagInfo
	flagSet    *flag.FlagSet
	flags      []*flag.Flag
	name       string
	Tag
}

// Tag store general values
type Tag struct {
	valid string
	json  string
	desc  string
	env   string
	def   string
	flag  []string
}

// tags fill the 'tag'
func tags[T Configuration](field string) TagInfo {
	var structAny T
	ret := TagInfo{}
	ret.name = field
	ret.flagSet = flag.NewFlagSet("", flag.ContinueOnError)
	tmp, exist := reflect.TypeOf(structAny).FieldByName(ret.name)
	if !exist {
		return TagInfo{}
	}
	if v, ok := tmp.Tag.Lookup("default"); ok {
		ret.def = v
	}
	if v, ok := tmp.Tag.Lookup("flag"); ok {
		ret.flag = strings.Split(v, ",")
		var all string
		for _, flagTag := range ret.flag {
			ret.flagSet.StringVar(&all, flagTag, ret.def, ret.desc)
			ret.flags = append(ret.flags, pointerFlag(flagTag, ret.flagSet))
		}
	}
	if v, ok := tmp.Tag.Lookup("env"); ok {
		ret.env = v
	}
	if v, ok := tmp.Tag.Lookup("text"); ok {
		ret.desc = v
	}
	if v, ok := tmp.Tag.Lookup("valid"); ok {
		ret.valid = v
	}
	if v, ok := tmp.Tag.Lookup("json"); ok {
		ret.json = v
	}
	ret.ConfigFile = func(m map[string]any) TagInfo {
		for k, v := range m {
			if ret.json == k {
				for _, f := range ret.flags {
					switch val := v.(type) {
					case string:
						f.DefValue = val
					default:
						f.DefValue = fmt.Sprintf("%v", val)
					}
					helpers.ToLog(
						f.Value.Set(f.DefValue),
					)
				}
			}
		}
		return ret
	}
	ret.DummyFlags = func() TagInfo {
		for _, flagTag := range ret.flag {
			if flag.Lookup(flagTag) == nil {
				flag.StringVar(new(string), flagTag, ret.def, ret.desc)
			}
		}
		return ret
	}
	ret.Flag = func() TagInfo {
		def := ret.flagSet.Output()
		ret.flagSet.SetOutput(io.Discard)
		for _, arg := range os.Args {
			err := ret.flagSet.Parse([]string{arg})
			helpers.ToLogWithType(err, helpers.LogNull)
		}
		ret.flagSet.SetOutput(def)
		return ret
	}
	ret.Env = func() TagInfo {
		env, ok := os.LookupEnv(ret.env)
		if ok {
			for _, f := range ret.flags {
				err := f.Value.Set(env)
				helpers.ToLog(err)
				break
			}
		}
		return ret
	}
	ret.Valid = func() any {
		value := ""
		for _, f := range ret.flags {
			value = f.Value.String()
			break
		}
		switch ret.valid {
		case "tmp_file":
			return validTempFile(value)
		case "file":
			return validFile(value)
		case "url":
			return validURL(value)
		case "bool":
			return validBool(value)
		case "timer":
			return validTimer(value)
		default:
			return value
		}
	}
	return ret
}

// pointerFlag return the flag pointer
func pointerFlag(name string, fs *flag.FlagSet) *flag.Flag {
	var current *flag.Flag
	fs.VisitAll(func(f *flag.Flag) {
		if f.Name == name {
			current = f
		}
	})
	return current
}

// Init find the info inside the tags
func (c *Configuration) Init() {
	c.tmp = make(map[string]TagInfo)
	elements := reflect.ValueOf(&c.ConfigurationMain).Elem()
	c.tmp = make(map[string]TagInfo, elements.NumField())
	for i := 0; i < elements.NumField(); i++ {
		name := elements.Type().Field(i).Name
		c.tmp[name] = tags[Configuration](name)
	}
}

// Dummy needed for the default printing of Flags info
func (c *Configuration) Dummy() {
	for _, v := range c.tmp {
		v.DummyFlags()
	}
	flag.Parse()
}

// Flag get info from the Flags
func (c *Configuration) Flag() {
	for _, v := range c.tmp {
		v.Flag()
	}
}

// Env get info from the Environment
func (c *Configuration) Env() {
	for _, v := range c.tmp {
		v.Env()
	}
}

// Conf get info from the configuration file
func (c *Configuration) Conf() {
	confFile := ""
	for _, f := range c.tmp["Config"].flags {
		if f.Value.String() != "" {
			confFile = f.Value.String()
		}
	}
	if confFile != "" {
		reflect.ValueOf(&c.ConfigurationMain).Elem().FieldByName("Config").Set(reflect.ValueOf(c.tmp["Config"].Valid()))
		tmpConfig := SettingsFile(c.Config)
		for _, v := range c.tmp {
			v.ConfigFile(tmpConfig)
		}
	}
}

// Valid check info and make some correcting
func (c *Configuration) Valid() {
	for k, v := range c.tmp {
		reflect.ValueOf(&c.ConfigurationMain).Elem().FieldByName(k).Set(reflect.ValueOf(v.Valid()))
	}
}

// Get create new conf.
func Get() Configuration {
	onceFlags.Do(
		func() {
			// ContinueOnError by default
			// flag.CommandLine.Init("", flag.ContinueOnError)
			// Section first init additional block 'maps'
			this.Init()

			// Dummy call for pretty print Usage()
			this.Dummy()

			// Section preload. For reading correct variant path to settings file.
			this.Flag() // 1 lowest
			this.Env()  // 2

			// Section settings weight by order
			this.Conf() // 1 lowest
			this.Flag() // 2
			this.Env()  // 3

			// Section valid last values
			this.Valid()

			defAddress := strings.Split(this.Address, ":")
			if len(defAddress) > 1 {
				newAddress1 := defAddress[0]
				newAddress2 := defAddress[1]
				if defAddress[0] == "ip" {
					newAddress1 = helpers.GetLocalIP()
				}
				if defAddress[1] == "0" {
					newAddress2 = helpers.OpenPort()
				}
				this.Address = fmt.Sprintf("%s:%s", newAddress1, newAddress2)
			}

			if ico, err := EmbedFiles.ReadFile(IcoNotif); err == nil {
				this.Ico = helpers.CreateTmp()
				os.WriteFile(this.Ico, ico, 0644)
			}
		})
	return this
}

// Set change conf.
func Set(c Configuration) {
	this = c
}

// Init if not need use init()
func Init() Configuration {
	return Get()
}

// validBool prepare bool string for Storage.
func validBool(v string) bool {
	if tmp, err := strconv.ParseBool(v); err == nil {
		return tmp
	}
	return false
}

// validTempFile prepare tmp filepath for Storage.
func validTempFile(v string) string {
	if strings.TrimSpace(v) == "" {
		return ""
	}
	var file string
	if runtime.GOOS == "windows" {
		check := strings.FieldsFunc(v, func(ss rune) bool {
			return strings.ContainsAny(string(ss), `\/`)
		})
		tmpString := filepath.Join(os.TempDir(), check[len(check)-1])
		file = tmpString
	} else {
		file = "/" + v
	}
	return file
}

// validFile prepare filepath for Storage.
func validFile(v string) string {
	check := strings.FieldsFunc(v, func(ss rune) bool {
		return strings.ContainsAny(string(ss), `\/`)
	})
	if strings.TrimSpace(v) == "" {
		return ""
	}
	if !helpers.FileExist(v) {
		helpers.ToLog(fmt.Sprintf("file '%s' not found", v))
	}
	return strings.Join(check, string(os.PathSeparator))
}

// validTimer prepare time string for Storage.
func validTimer(v string) time.Duration {
	if v[0] == '-' || v[0] == '+' {
		v = v[1:]
	}
	l := len(v) - 1
	if '0' <= v[l] && v[l] <= '9' {
		v += DefTimePostfix
	}
	tmp, err := time.ParseDuration(v)
	helpers.ToLog(err)
	return tmp
}

// validURL prepare url.
func validURL(v string) string {
	// trim prefix
	re := regexp.MustCompile(`^.*(://|^)[^/]+`)
	trimPrefix := re.FindString(v)
	re = regexp.MustCompile(`^.*(://|^)`)
	fullAddress := re.ReplaceAllString(trimPrefix, "")
	// trim port
	re = regexp.MustCompile(`^[^/:$]+`)
	address := re.FindString(fullAddress)
	// fill address
	if strings.TrimSpace(address) == "" {
		return ""
	}
	// check ip
	isIP := false
	re = regexp.MustCompile(`\d+`)
	isIPTest := re.ReplaceAllString(address, "")
	isIPTest = strings.ReplaceAll(isIPTest, ".", "")
	if strings.TrimSpace(isIPTest) == "" {
		isIP = true
	}
	// correct IP
	if isIP {
		re = regexp.MustCompile(`\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}`)
		addressIP := re.FindString(address)
		if strings.TrimSpace(addressIP) == "" {
			return ""
		}
	}
	// check and correct port
	re = regexp.MustCompile(`:.*`)
	correctPort := re.FindString(fullAddress)
	correctPort = strings.Replace(correctPort, ":", "", 1)
	re = regexp.MustCompile(`\D`)
	correctPort = re.ReplaceAllString(correctPort, "")
	correctPort = strings.Replace(correctPort, ":", "", 1)
	if strings.TrimSpace(correctPort) == "" {
		return address + ":80"
	}
	return address + ":" + correctPort
}

// SettingsFile return map[string]any from the setting file
func SettingsFile(filename string) (compare map[string]any) {
	// generate conf file with default if file not exist?
	// Examples:
	//
	// {
	// 	"test1": {
	// 		"actions": "EXEC@pathtoexe|SLEEP@1s|TARGET@window_name|SLEEP@1s|PRESS@1|SLEEP@1s|PRINT@testtesttest|SLEEP@1s|GEN@12",
	// 		"repeat": 3,
	// 		"timer": "30s",
	// 		"silent": false
	// 	},
	// 	"test2": {
	// 		"actions": "TARGET@Notepad|GEN@-16|SLEEP@300ms|PRESS@CTRL+V|SLEEP@300ms|PRESS@ENTER|PRINT@||SLEEP@300ms|EXECSTD@ping google.com|SLEEP@300ms|PRESS@CTRL+V|SLEEP@300ms",
	// 		"repeat": 2,
	// 		"timer": "",
	// 		"silent": true
	// 	}
	// }
	//
	//
	// {
	// 	"config": "settings.data",
	// 	"ico_file": "app.ico",
	// 	"address": "localhost:8080",
	// }
	//
	if runtime.GOOS != "windows" {
		filename = "/" + filename
	}

	f, err := os.ReadFile(filename)
	helpers.ToLog(err)
	err = json.Unmarshal(f, &compare)
	helpers.ToLog(err)
	return
}

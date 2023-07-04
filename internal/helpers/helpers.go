// Package helpers working with logging
// and other non-main/other help-function.
package helpers

import (
	"bytes"
	"compress/gzip"
	cryptoRand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	mathRand "math/rand"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Text output types
const (
	dnsServerTest = "8.8.8.8:80"
)

const (
	LogOutNullTS = iota
	LogOutFastTS
	LogOutHumanTS
)

// Error output types
const (
	LogErrNullTS = iota + 100
	LogErrFastTS
	LogErrHumanTS
)

// LogNull Error null
const (
	LogNull = 1000
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

// OpenUrl open URL.
func OpenUrl(url string) {
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
	ToLog(err)
}

// ToLog splitting notifications to Err and Debug.
func ToLog(err any) {
	if err == nil || err == "" {
		return
	}
	switch val := err.(type) {
	case error:
		if val.Error() != "" {
			log.Error().Msgf("%v", val)
		}
	default:
		log.Debug().Msgf("%v", val)
	}
}

// ToLogWithType splitting notifications to Err and Debug. With the type of timestamps.
func ToLogWithType(err any, typ int) {
	if err == nil || err == "" {
		return
	}
	if typ == LogNull {
		std := logger(typ)
		std.Debug().Msgf("%v", err)
		return
	}
	switch val := err.(type) {
	case error:
		if val.Error() != "" {
			std := logger(typ + LogErrNullTS)
			std.Error().Msgf("%v", val)
		}
	default:
		std := logger(typ)
		std.Debug().Msgf("%v", val)
	}
}

// logger return logger with TS.
func logger(typ int) zerolog.Logger {
	switch typ {
	case LogOutNullTS:
		return zerolog.New(os.Stdout).With().Logger()
	case LogErrNullTS:
		return zerolog.New(os.Stderr).With().Logger()
	case LogOutFastTS:
		return zerolog.New(os.Stdout).With().Timestamp().Logger()
	case LogErrFastTS:
		return zerolog.New(os.Stderr).With().Timestamp().Logger()
	case LogOutHumanTS:
		return zerolog.New(os.Stdout).With().Str("time", time.Now().Format("200601021504")).Logger()
	case LogErrHumanTS:
		return zerolog.New(os.Stderr).With().Str("time", time.Now().Format("200601021504")).Logger()
	case LogNull:
		return zerolog.New(io.Discard).With().Logger()
	default:
		return zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
}

// Compress data.
func Compress(data []byte, levelCompression int) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, levelCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}

// UnCompress data.
func UnCompress(data []byte) ([]byte, error) {
	w, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	b, _ := io.ReadAll(w)
	return b, nil
}

// Random return random value using channel.
func Random() <-chan float64 {
	var rand = mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	min, max := float64(100), float64(999)
	r := min + rand.Float64()*(max-min)
	c := make(chan float64, 1)
	go func() {
		c <- r
	}()
	return c
}

// Generate alias for RandomString.
func Generate(i int) string {
	return RandomString(i)
}

// RandomString - random string with 'i' length.
func RandomString(l int) string {
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		switch RandInt(1, 7) {
		// case 1:
		// 	bytes[i] = byte(RandInt(48, 57))
		case 2:
			b[i] = byte(RandInt(65, 90))
		case 3:
			b[i] = byte(RandInt(97, 122))
		default:
			b[i] = byte(RandInt(97, 122))
		}
	}
	return string(b)
}

// RandomStringPass - random difficult string with 'i' length.
func RandomStringPass(l int) string {
	b := make([]byte, l)
	spc := []int{33, 35, 36, 37, 38, 42, 43, 40, 41, 60, 62, 63, 64, 123, 125, 91, 93}
	checkABC := int(math.Round(float64(l) * 0.25))
	checkAbc := 0
	check123 := 0
	checkSpc := 0
	if checkABC < l {
		checkAbc = int(math.Round(float64(l) * 0.5))
	}
	if checkABC+checkAbc < l {
		check123 = int(math.Round(float64(l) * 0.125))
	}
	if checkABC+checkAbc+check123 < l {
		checkSpc = int(math.Round(float64(l) * 0.125))
	}
	for i := 0; i < checkABC; i++ {
		rval := RandInt(0, int64(l-1))
		if b[rval] == 0 {
			b[rval] = byte(RandInt(65, 90))
		} else {
			i--
		}
	}
	for i := 0; i < checkAbc; i++ {
		rval := RandInt(0, int64(l-1))
		if b[rval] == 0 {
			b[rval] = byte(RandInt(97, 122))
		} else {
			i--
		}
	}
	for i := 0; i < check123; i++ {
		rval := RandInt(0, int64(l-1))
		if b[rval] == 0 {
			b[rval] = byte(RandInt(48, 57))
		} else {
			i--
		}
	}
	for i := 0; i < checkSpc; i++ {
		rval := RandInt(0, int64(l-1))
		if b[rval] == 0 {
			r := RandInt(1, 17)
			b[rval] = byte(spc[r-1])
		} else {
			i--
		}
	}
	return string(b)
}

// RandInt generate int value between min and max.
func RandInt(min, max int64) int {
	bval, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(max-min+1))
	ToLog(err)
	val := bval.Int64()
	if val < min {
		return int(val + min)
	}
	return int(val)
}

// Round float64 to value 'count' digit after dot.
func Round(v float64, count int) float64 {
	base := math.Pow(10, float64(count))
	return math.Round(v*base) / base
}

// SplitPrefix split names by prefix as "ST0001" --> "ST" "0001"
func SplitPrefix(val string, pre ...string) (string, string) {
	if pre == nil {
		pre = append(pre, "")
	}
	if len(val) > 0 {
		r := rune(val[0])
		if r < 48 || r > 57 { // non digit
			return SplitPrefix(val[1:], pre[0]+string(val[0]))
		}
	}
	return string(pre[0]), val
}

// SplitPrefixLine split strings first line and other lines
func SplitPrefixLine(val string, pre ...string) (string, string) {
	if pre == nil {
		pre = append(pre, "")
	}
	if len(val) > 0 {
		r := rune(val[0])
		if r != 10 { // '\n'
			return SplitPrefixLine(val[1:], pre[0]+string(val[0]))
		}
	}
	return string(pre[0]), val
}

// Concatenate join slices
func Concatenate[T any](values ...[]T) []T {
	count := 0
	for _, value := range values {
		count += len(value)
	}
	all := make([]T, 0, count)
	for _, value := range values {
		all = append(all, value...)
	}
	return all
}

// MapToSlice convert the map to the slice
func MapToSlice[K comparable, V any](m map[K]V) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

// FindError find errors
func FindError[T error](err error) error {
	var check T
	if errors.As(err, &check) {
		return check
	}
	return err
}

// CreateTmp create file in Temp folder
func CreateTmp() string {
	fileEnv, err := os.CreateTemp("", "tmp_golang_")
	ToLog(err)
	defer func(path string) {
		ToLog(fileEnv.Close())
	}(fileEnv.Name())
	if runtime.GOOS != "windows" {
		return string(os.PathSeparator) + fileEnv.Name()
	}
	return fileEnv.Name()
}

// DeleteTmp delete file in Temp folder. Actually imply delete any file by path
func DeleteTmp(path string) {
	ToLog(os.RemoveAll(path))
}

// FileExist check exist
func FileExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

// OpenPort return open port
func OpenPort() string {
	listener, _ := net.Listen("tcp", ":0")
	defer listener.Close()
	return strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
}

// GetLocalIP return local ip
func GetLocalIP() string {
	conn, _ := net.Dial("udp", dnsServerTest)
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

// GetLocalIPs return local ip
func GetLocalIPs() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		ToLog(err)
		return ""
	}
	IPs := make([]string, 0)
	for _, inter := range interfaces {
		if inter.Flags&net.FlagUp == 0 {
			continue
		}
		if inter.Flags&net.FlagLoopback != 0 {
			continue
		}
		addresses, _ := inter.Addrs()
		for _, address := range addresses {
			if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil && ipNet.IP.IsPrivate() {
					IPs = append(IPs, ipNet.IP.To4().String())
				}
			}
		}
	}
	return strings.Join(IPs, ",")
}

// CheckIP check ip into CIDR
func CheckIP(data string, dataCIDR string) bool {
	IPs := strings.Split(data, ",")
	_, addr, err := net.ParseCIDR(dataCIDR)
	if err != nil {
		return false
	}
	for _, p := range IPs {
		if addr.Contains(net.ParseIP(p)) {
			return true
		}
	}
	return false
}

// Print - clear spaces
func Print(v any) {
	fmt.Println(
		strings.ReplaceAll(
			fmt.Sprint(v),
			"  ",
			" ",
		),
	)
}

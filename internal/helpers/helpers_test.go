package helpers_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/rs/zerolog"

	"github.com/KusoKaihatsuSha/tray_helper/internal/helpers"
)

func Example_logs() {

	var nilerr error
	nilrr := struct {
		name string
		eman int
	}{name: "foo", eman: 1}

	var tests []any
	tests = append(
		tests,
		"test",
		"", // skip
		17,
		0,
		true,
		false,
		errors.New("test error"), // prints double
		errors.New(""),           // skip
		fmt.Errorf(""),           // skip
		nilerr,                   // skip
		fmt.Errorf("1:%w", fmt.Errorf("2:%w", fmt.Errorf("3:%w", errors.New("4")))), // prints double
		nilrr,
	)

	tmpOut, tmpErr := os.Stdout, os.Stderr
	r, w, _ := os.Pipe()
	re, we, _ := os.Pipe()
	os.Stdout, os.Stderr = w, we
	// printings
	bakLog := zerolog.GlobalLevel()
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	for _, tt := range tests {
		helpers.ToLogWithType(tt, helpers.LogOutNullTS)
	}
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	for _, tt := range tests {
		helpers.ToLogWithType(tt, helpers.LogOutNullTS)
	}
	zerolog.SetGlobalLevel(bakLog)
	// end printings
	c, ce := make(chan string), make(chan string)

	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			panic("check test logs files")
		}
		c <- buf.String()
	}()

	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, re)
		if err != nil {
			panic("check test logs files")
		}
		ce <- buf.String()
	}()

	err := w.Close()
	if err != nil {
		panic("check test logs files")
	}

	err = we.Close()
	if err != nil {
		panic("check test logs files")
	}

	os.Stdout, os.Stderr = tmpOut, tmpErr
	out, outerr := <-c, <-ce

	fmt.Println("DEBUG:")
	fmt.Print(out)

	fmt.Println("ERROR:")
	fmt.Print(outerr)

	// Output:
	// DEBUG:
	// {"level":"debug","message":"test"}
	// {"level":"debug","message":"17"}
	// {"level":"debug","message":"0"}
	// {"level":"debug","message":"true"}
	// {"level":"debug","message":"false"}
	// {"level":"debug","message":"{foo 1}"}
	// ERROR:
	// {"level":"error","message":"test error"}
	// {"level":"error","message":"1:2:3:4"}
	// {"level":"error","message":"test error"}
	// {"level":"error","message":"1:2:3:4"}
}

func Example_compress() {
	text := "testtesttesttesttest"
	fmt.Println(text)
	data := []byte(text)
	compress, err := helpers.Compress(data, 9)
	if err == nil {
		fmt.Println(compress)
	}
	unCompress, err := helpers.UnCompress(compress)
	if err == nil {
		fmt.Println(string(unCompress))
	}
	if text == string(unCompress) {
		fmt.Println(true)
	}

	// Output:
	// testtesttesttesttest
	// [31 139 8 0 0 0 0 0 2 255 42 73 45 46 65 199 128 0 0 0 255 255 109 80 182 70 20 0 0 0]
	// testtesttesttesttest
	// true
}

func Example_rand() {
	test := helpers.RandInt(105, 1205)
	if test >= 105 {
		fmt.Println(true)
	}
	if test <= 1205 {
		fmt.Println(true)
	}

	// Output:
	// true
	// true
}

func Example_round() {
	test := 123456.664321
	fmt.Println(helpers.Round(test, 4))
	fmt.Println(helpers.Round(test, 1))

	// Output:
	// 123456.6643
	// 123456.7
}

func Example_tolog1() {
	// another stdout
	helpers.ToLog("test")
	helpers.ToLog(errors.New("test"))

	// Output:

}

func Example_tolog2() {
	// another stdout
	helpers.ToLogWithType("test", helpers.LogErrFastTS)
	helpers.ToLogWithType("test", helpers.LogErrHumanTS)
	helpers.ToLogWithType("test", 999)

	// Output:

}

func Example_randomString() {
	// another stdout
	r1 := helpers.RandomString(7)
	r2 := helpers.RandomString(7)

	if r1 == r2 {
		fmt.Println(false)
	}

	// Output:

}

func Example_splitPrefix() {
	text001 := `first
second
------`
	text002 := "first\nsecond------"
	text003 := "SA00001"
	text004 := "SA00002"
	text005 := "00003"
	text006 := "00004S"
	tf, _ := helpers.SplitPrefixLine(text001)
	fmt.Println(tf)
	tf, _ = helpers.SplitPrefixLine(text002)
	fmt.Println(tf)
	tf, _ = helpers.SplitPrefix(text003)
	fmt.Println(tf)
	tf, _ = helpers.SplitPrefix(text004)
	fmt.Println(tf)
	tf, _ = helpers.SplitPrefix(text005)
	fmt.Println(tf)
	tf, _ = helpers.SplitPrefix(text006)
	fmt.Println(tf)

	// Output:
	// first
	// first
	// SA
	// SA

}

func Example_mapToSlice() {
	m := make(map[int]string, 5)
	m[1] = "1"
	m[2] = "2"
	m[3] = "3"
	m[4] = "4"
	m[5] = "5"
	s := helpers.MapToSlice(m)
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
	fmt.Printf("%#v", s)
	// Output:
	// []string{"1", "2", "3", "4", "5"}

}

type TestErr struct {
	Text  string
	Value int
}

func (e TestErr) Error() string {
	return fmt.Sprintf("%#v", e)
}

func Example_errorsUnwrap() {
	e1 := TestErr{
		Text:  "-1-",
		Value: 1234,
	}
	e2 := fmt.Errorf("new 1: %w", e1)
	e3 := fmt.Errorf("new 2: %w", e2)
	e4 := fmt.Errorf("new 3: %w", e3)
	e5 := fmt.Errorf("new 4: %w", e4)
	exact := helpers.FindError[TestErr](e5)
	fmt.Println(e5)
	fmt.Printf("%#v", exact)
	// Output:
	// new 4: new 3: new 2: new 1: helpers_test.TestErr{Text:"-1-", Value:1234}
	// helpers_test.TestErr{Text:"-1-", Value:1234}

}

func Example_tmp() {
	n := helpers.CreateTmp()
	err := os.WriteFile(n, []byte("test"), 0755)
	helpers.ToLog(err)
	text, err := os.ReadFile(n)
	fmt.Println("body file:", string(text))
	fmt.Println("exist file before delete:", helpers.FileExist(n))
	helpers.DeleteTmp(n)
	fmt.Println("exist file after delete:", helpers.FileExist(n))
	// Output:
	// body file: test
	// exist file before delete: true
	// exist file after delete: false

}

func Example_checkIP() {
	fmt.Println(helpers.CheckIP("192.168.1.111", "192.168.1.1/30"))
	fmt.Println(helpers.CheckIP("192.168.1.111", "192.168.1.1/24"))
	fmt.Println(helpers.CheckIP("192.168.1.111", "192.168.1.1/16"))
	fmt.Println(helpers.CheckIP("192.168.1.111", "192.168.1.1/8"))
	fmt.Println(helpers.CheckIP("127.10.1.111", "192.168.1.1/1"))
	fmt.Println(helpers.CheckIP("127.10.1.111", "192.168.1.1/0"))

	// Output:
	// false
	// true
	// true
	// true
	// false
	// true

}

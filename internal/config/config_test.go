package config_test

import (
	"fmt"
	"os"

	"github.com/KusoKaihatsuSha/tray_helper/internal/config"
	"github.com/KusoKaihatsuSha/tray_helper/internal/helpers"
)

func Example_configOrder() {
	testDataFlag := `{
                 "address": "http://file-flag"
         }`
	testDataEnv := `{
                 "address": "http://file-env"
         }`
	fileFlag := helpers.CreateTmp()
	err := os.WriteFile(fileFlag, []byte(testDataFlag), 0755)
	if err != nil {
		fmt.Println(err)
	}
	fileEnv := helpers.CreateTmp()
	err = os.WriteFile(fileEnv, []byte(testDataEnv), 0755)

	if err != nil {
		fmt.Println(err)
	}

	err = os.Setenv("ADDRESS", "http://env")
	if err != nil {
		fmt.Println(err)
	}
	err = os.Setenv("CONFIG", fileEnv)
	if err != nil {
		fmt.Println(err)
	}

	os.Args = append(
		os.Args,
		"-config="+fileFlag,
		"-a=http://flag",
	)
	testConfiguration := config.Init()
	testConfiguration.Init()

	testConfiguration.Flag() // preload
	testConfiguration.Env()  // preload

	// test section
	testConfiguration.Conf() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Flag() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Env() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Conf()
	testConfiguration.Flag() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Conf()
	testConfiguration.Env() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Flag()
	testConfiguration.Env() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Conf()
	testConfiguration.Flag()
	testConfiguration.Env() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Flag()
	testConfiguration.Conf()
	testConfiguration.Env() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Conf()
	testConfiguration.Env()
	testConfiguration.Flag() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Env()
	testConfiguration.Conf()
	testConfiguration.Flag() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Flag()
	testConfiguration.Env()  // last filepath config here
	testConfiguration.Conf() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	testConfiguration.Env()
	testConfiguration.Flag() // last filepath config here
	testConfiguration.Conf() // <--
	testConfiguration.Valid()
	fmt.Println(testConfiguration.Address)

	helpers.DeleteTmp(fileEnv)
	helpers.DeleteTmp(fileFlag)

	// Output:
	// file-env:80
	// flag:80
	// env:80
	// flag:80
	// env:80
	// env:80
	// env:80
	// env:80
	// flag:80
	// flag:80
	// file-env:80
	// file-flag:80

}

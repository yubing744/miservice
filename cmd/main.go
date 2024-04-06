package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/longbai/miservice"
)

func usage() {
	fmt.Printf("MiService - XiaoMi Cloud Service\n")
	fmt.Printf("Usage: The following variables must be set:\n")
	fmt.Printf("           export MI_USER=<Username>\n")
	fmt.Printf("           export MI_PASS=<Password>\n")
	fmt.Printf("           export MI_DID=<Device ID|Name>\n\n")
	fmt.Printf(miservice.IOCommandHelp("", os.Args[0]+" "))
}

func main() {
	args := os.Args
	argCount := len(args)
	argIndex := 1

	if argCount > argIndex {
		token := fmt.Sprintf("%s/.mi.token", os.Getenv("HOME"))
		account := miservice.NewAccount(os.Getenv("MI_USER"),
			os.Getenv("MI_PASS"),
			miservice.NewTokenStore(token),
		)

		var result interface{}
		var err error
		cmd := strings.Join(args[argIndex:], " ")

		service := miservice.NewIOService(account, nil)
		result, err = miservice.IOCommand(service, os.Getenv("MI_DID"), cmd, os.Args[0]+" ")

		if err != nil {
			fmt.Println(err)
		} else {
			if resStr, ok := result.(string); ok {
				fmt.Println(resStr)
			} else {
				resBytes, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(resBytes))
			}
		}
	} else {
		usage()
	}
}

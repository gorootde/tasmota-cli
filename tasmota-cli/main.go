package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var verbose bool

func verbosePrintln(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a)
		fmt.Printf("\n")
	}
}

func printUsage() {
	fmt.Printf("Usage: %s [OPTIONS] <command> <ip> (<ip> <ip> ...)\n", os.Args[0])
	flag.PrintDefaults()
}

func sendCommand(ip string, command string, user string, password string) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("http://%s/cm", ip)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	if len(password) > 0 && len(user) > 0 {
		req.SetBasicAuth(user, password)
	}

	q := req.URL.Query()
	q.Add("cmnd", command)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		verbosePrintln("Error: %s", err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		verbosePrintln("Response: %s", string(body))
		var parsed map[string]map[string]interface{}
		json.Unmarshal([]byte(body), &parsed)

		return parsed["StatusFWR"]["Version"].(string), nil
	}
	return "", fmt.Errorf("HTTP %d. %s", resp.StatusCode, body)

}

func main() {
	flag.Usage = printUsage
	user := flag.String("u", "", "The username used to authenticate to tasmota")
	password := flag.String("p", "", "The password used to authenticate to tasmota")
	flag.BoolVar(&verbose, "v", false, "Enable verbose mode")
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := flag.Args()[0]
	for i := 1; i < len(flag.Args()); i++ {
		ip := flag.Args()[i]
		fmt.Printf("%s\t\t", ip)
		var result string
		var error error
		switch command {
		case "version":
			result, error = sendCommand(ip, "Status 2", *user, *password)
		}

		if error == nil {
			fmt.Printf("%+v", result)
		} else {
			fmt.Printf("Error: %+v", error)
		}
		fmt.Printf("\n")
	}
}

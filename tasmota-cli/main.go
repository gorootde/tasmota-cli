package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var verbose bool
var legacyAuth bool

func verbosePrintln(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a...)
		fmt.Printf("\n")
	}
}

func printUsage() {
	fmt.Printf("Usage: %s [OPTIONS] <command> <ip> (<ip> <ip> ...)\n\n", os.Args[0])
	fmt.Println("  <command> 	Any tasmota command. See https://tasmota.github.io/docs/Commands/")
	fmt.Println("  <ip> 		List of the IPs of the tasmota devices to execute the command on")
	fmt.Println("\nCLI Arguments always take precedence over environment eariables")
	fmt.Println("  TASMOTACLI_USERNAME  Username")
	fmt.Println("  TASMOTACLI_PASSWORD  Password")
	fmt.Println("\nFlags:")
	flag.PrintDefaults()
}

func downloadBackup(ip string, user string, password string) (string, error) {

	url := fmt.Sprintf("http://%s/dl?", ip)

	req, _ := http.NewRequest("GET", url, nil)

	addAuthentication(req, user, password)

	filename := fmt.Sprintf("tasmota_backup_%s.bin", ip)
	fp, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	resp, err := performRequest(req)
	if err != nil {
		verbosePrintln("Error: %s", err)
		return "", err
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		verbosePrintln("Response: %d %s", resp.StatusCode, string(body))
		return string(body), nil
	}
	defer resp.Body.Close()
	_, err = io.Copy(fp, resp.Body)

	return "OK", nil
}

func performRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: 1000 * time.Millisecond,
	}
	return client.Do(req)

}

func addAuthentication(req *http.Request, user string, password string) {
	//Authentication
	if len(password) > 0 || len(user) > 0 {
		if legacyAuth {
			q := req.URL.Query()
			q.Add("user", user)         //required for Tasmota <= v9.2.0
			q.Add("password", password) //required for Tasmota <= v9.2.0
			req.URL.RawQuery = q.Encode()
		} else {
			req.SetBasicAuth(user, password)
		}
	}
}

func sendCommand(ip string, command string, user string, password string) (string, error) {
	url := fmt.Sprintf("http://%s/cm", ip)

	client := &http.Client{
		Timeout: 1000 * time.Millisecond,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")

	addAuthentication(req, user, password)

	//Create Querystring
	q := req.URL.Query()
	q.Add("cmnd", command)
	req.URL.RawQuery = q.Encode()

	//Perform request
	resp, err := client.Do(req)
	if err != nil {
		verbosePrintln("Error: %s", err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		verbosePrintln("Response: %s", string(body))
		return string(body), nil
	}
	return "", fmt.Errorf("HTTP %d. %s", resp.StatusCode, body)
}

func parseFwVersion(response string) string {
	var parsed map[string]map[string]interface{}
	json.Unmarshal([]byte(response), &parsed)
	return parsed["StatusFWR"]["Version"].(string)
}

func main() {
	flag.Usage = printUsage
	var user string
	var password string

	flag.StringVar(&user, "u", os.Getenv("TASMOTACLI_USERNAME"), "The username used to authenticate to tasmota")
	flag.StringVar(&password, "p", os.Getenv("TASMOTACLI_PASSWORD"), "The password used to authenticate to tasmota")
	flag.BoolVar(&legacyAuth, "la", false, "Enable legacy authentication mode (for tasmota versions <= 9.2.0)")
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
			result, error = sendCommand(ip, "Status 2", user, password)
			if error == nil {
				result = parseFwVersion(result)
			}
		case "backup":
			result, error = downloadBackup(ip, user, password)
		default:
			result, error = sendCommand(ip, command, user, password)
		}

		if error == nil {
			fmt.Printf("%+v", result)
		} else {
			fmt.Printf("Error: %+v", error)
		}
		fmt.Printf("\n")
	}
}

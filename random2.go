package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func checkURL(url string, client *http.Client) bool {
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// Check if the response body is empty
	return len(body) == 0
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [domain]")
		os.Exit(1)
	}

	domain := os.Args[1]
	urls := []string{
		"http://" + domain + "/wp-content/plugins/forminator/forminator.php",
		"https://" + domain + "/wp-content/plugins/forminator/forminator.php",
	}

	// Custom HTTP client with SSL verification disabled
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	for _, url := range urls {
		if checkURL(url, client) {
			fmt.Printf("Success: %s\n", url)
		} else {
			fmt.Printf("Fail: %s\n", url)
		}
	}
}

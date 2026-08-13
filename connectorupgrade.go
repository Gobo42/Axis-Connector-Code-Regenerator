package main

import (
    "fmt"
    "bufio"
    "encoding/json"
    "io"
    "os"
    "crypto/tls"
    "net/http"
    "strings"
    "strconv"
)

type entry struct {
    ID string `json:"id"`
    Name string `json:"name"`
}

func parseEntries(msg []byte) ([]entry, error) {
    var response struct {
        Data []entry `json:"data"`
    }
    if err := json.Unmarshal(msg, &response); err != nil {
        return nil, err
    }
    return response.Data, nil
}

func main() {
    fmt.Println("Axis Connector Install String Regenerator")
    fmt.Println("v1.0.1 by matt.hum@hpe.com")

    apikey := ""
    dat, err := os.ReadFile("apikey")
    if err != nil {
        fmt.Println("Missing API Key")
        fmt.Print("Enter Key here: ")
        reader := bufio.NewReader(os.Stdin)
        text,_ := reader.ReadString('\n')
        text = strings.Replace(text, "\n","",-1)
        
        f, err := os.Create("apikey")
        if err!=nil {
            fmt.Println("Couldn't open file for opening")
        }
        defer f.Close()

        w:=bufio.NewWriter(f)
        _, err = w.WriteString(text)
        if err!=nil {
            fmt.Println("Couldn't write API key to file")
        }
        apikey = text
        w.Flush()
    } else {
        apikey = string(dat)
    }
    bearer:= "Bearer " + strings.TrimSpace(apikey)

    tr := &http.Transport {
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    url := "https://admin-api.axissecurity.com/api/v1/connectors?pageSize=100&pageNumber=1"
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        fmt.Println("Error formatting request")
        panic(1)
    }
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Authorization", bearer)

    res, err := client.Do(req)
    if err != nil {
        fmt.Println("status code: ", res.StatusCode)
        fmt.Println("Error sending req: ", err)
        panic(1)
    }
    defer res.Body.Close()
    if res.StatusCode != 200 {
        fmt.Println("Error, got status code: ", res.StatusCode)
        fmt.Println(res)
        panic(1)
    }
    msg, _ := io.ReadAll(res.Body)
    
    var arr []entry
    arr, err = parseEntries(msg)
    if err != nil {
        fmt.Println("Error parsing connector response:", err)
        panic(1)
    }
    count:=len(arr)
    fmt.Println("I found", count, "connectors")
    for i:=0; i<count; i++ {
        fmt.Printf("%v: %v\n",i,arr[i].Name)
    }
    text:=""
    fmt.Print("Enter number of connector to regen a command: ")
    reader := bufio.NewReader(os.Stdin)
    text,_ = reader.ReadString('\n')
    text = strings.Replace(text, "\n","",-1)
    num, err :=strconv.Atoi(text)
    if err != nil {
        panic("Couldn't read number")
    }
    fmt.Printf("Regenerating command for %v\n",arr[num].Name)
    url="https://admin-api.axissecurity.com/api/v1/connectors/"+arr[num].ID+"/regenerate"

    req, err = http.NewRequest(http.MethodPost, url, nil)
    if err != nil {
        fmt.Println("Error formatting request")
        panic(1)
    }
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Authorization", bearer)

    res, err = client.Do(req)
    if err != nil {
        fmt.Println("status code: ", res.StatusCode)
        fmt.Println("Error sending req: ", err)
        panic(1)
    }
    defer res.Body.Close()
    msg, _ = io.ReadAll(res.Body)

    b:=strings.Split(string(msg),",")
    for i:=0; i<len(b); i++ {
        if strings.Contains(b[i],"command") {
            c:=strings.Split(b[i],"\"")
            fmt.Println(c[3])
        }
    }
}

package main

import (
    "fmt"
    "bufio"
    "io"
    "os"
    "crypto/tls"
    "net/http"
    "strings"
    "strconv"
)


func main() {
    fmt.Println("Axis Connector Install String Regenerator")
    fmt.Println("v1.0.0 by matt.hum@hpe.com")

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
    
    data := strings.SplitAfter(string(msg),"[")[1]
    data2 := strings.Split(data,"]")[0]
    connectors:= strings.Split(data2,"},")
    count:=len(connectors)
    fmt.Println("I found", count, "connectors")
    type entry struct {
        id string
        name string
    }
    var arr []entry
    for i:=0; i<count; i++ {
        b:= strings.Split(connectors[i],",")
        var t entry
        for j:=0; j<4; j++ {
            if strings.Contains(b[j],"id") {
                c:=strings.Split(b[j],"\"")
                t.id = strings.TrimSpace(c[3])
            }
            if strings.Contains(b[j],"\"name\"") {
                c:=strings.Split(b[j],"\"")
                t.name = strings.TrimSpace(c[3])
                fmt.Printf("%v: %v\n",i,c[3])
            }
        }
        arr = append(arr, t)
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
    fmt.Printf("Regenerating command for %v\n",arr[num].name)
    url="https://admin-api.axissecurity.com/api/v1/connectors/"+arr[num].id+"/regenerate"

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

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
//    var msg = []byte(`{"pageNumber":1,"pageSize":100,"firstPage":"http://admin-api.axissecurity.com/api/v1/connectors?pageNumber=1&pageSize=100","lastPage":"http://admin-api.axissecurity.com/api/v1/connectors?pageNumber=1&pageSize=100","totalPages":1,"totalRecords":6},"data":[{"enabled":true,"connectorZoneId":"4ba50724-69a9-45c6-ac36-05e080bd1635","name":"AWS-SELab","id":"9360b5f6-6882-461c-a5c6-ea4517826fa0"},{"enabled":true,"connectorZoneId":"46c58153-487e-44d2-abdc-3d6f1d0bd295","name":"HomeLab-2","id":"149e5798-00f3-48b5-9e92-463c1d71211e"},{"enabled":true,"connectorZoneId":"168d7790-eadf-499f-9546-9dd6b7393707","name":"LabEnv-1","id":"35bb71bc-1344-4108-87a6-fef2018c1cec"},{"enabled":true,"connectorZoneId":"24e65e14-3f54-4d1a-9edf-60cb6fbee88d","name":"Log Streaming Connector_001","id":"43b00820-6f32-46c1-9082-cd051aa9b279"},{"enabled":true,"connectorZoneId":"46c58153-487e-44d2-abdc-3d6f1d0bd295","name":"test","id":"77ddb467-2ef1-4ebc-a7d3-367cd5a0d058"},{"enabled":true,"connectorZoneId":"46c58153-487e-44d2-abdc-3d6f1d0bd295","name":"test2","id":"2674be7e-6dba-46a6-9491-a6921685abbe"}]}`)
    
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
//                d:=strings.ReplaceAll(c[1],"\"","")
//                e:=strings.ReplaceAll(d,"}","")
                t.id = strings.TrimSpace(c[3])
            }
            if strings.Contains(b[j],"\"name\"") {
                c:=strings.Split(b[j],"\"")
//                d:=strings.ReplaceAll(c[1],"\"","")
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

    //msg = []byte(`{"id":"2674be7e-6dba-46a6-9491-a6921685abbe","command":"sudo bash < <(curl -fsSL https://ops.axissecurity.com/1tAGfvDS/install)"}`)
    b:=strings.Split(string(msg),",")
    for i:=0; i<len(b); i++ {
        if strings.Contains(b[i],"command") {
            c:=strings.Split(b[i],"\"")
            fmt.Println(c[3])
        }
    }
}

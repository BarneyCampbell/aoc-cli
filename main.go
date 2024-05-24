package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Method int
const (
    Input Method = iota
    Submit // todo
)

type Settings struct {
    Token string            `json:"token"`
    DownloadLocation string `json:"location"`
}

func main() {
    yearFlag := flag.Int("year", -1, "The year of advent of code to query.")
    flag.IntVar(yearFlag, "y", -1, "Shorthand of -year")

    dayFlag := flag.Int("day", -1, "The day of the advent of code to query.")
    flag.IntVar(dayFlag, "d", -1, "Shorthand of -day")

    filepathFlag := flag.String("filepath", "", "The location to download the input file to")
    flag.StringVar(filepathFlag, "f", "", "Shorthand of -filepath")

    filenameFlag := flag.String("filename", "", "The filename to download to")
    flag.StringVar(filenameFlag, "n", "", "Shorthand of -filename")

    toStdout := flag.Bool("stdout", false, "Whether to print the output or save to file")
    flag.Parse()

    if *yearFlag != -1 && *yearFlag < 2015 || *yearFlag > 2023 {
        log.Println("Please specify a valid year: 2015 <= year <= 2023")
        *yearFlag = -1
    }
    if *dayFlag != -1 && *dayFlag <= 1 || *dayFlag >= 25 {
        log.Println("Please specfy a valid day: 1 <= day <= 25")
        *yearFlag = -1
    }

    var year int
    var day int
    if *yearFlag == -1 {
        year = getNumFromCmd("Year")
        if year < 2015 || year > 2023 {
            log.Fatal("Invalid year: 2015 <= year <= 2023")
        }
    } else {
        year = *yearFlag
    }

    if *dayFlag == -1 {
        day = getNumFromCmd("Day")
        if day <= 1 || day >= 25 {
            log.Fatal("Invalid day: 1 <= day <= 25")
        }
    } else {
        day = *dayFlag
    }

    var settings Settings
    err := getConfig(&settings)

    if err != nil {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("No config found. Please enter your access token: ")

        text, err := reader.ReadString('\n')
        text = strings.Split(text, "\n")[0]
        checkErr(err)

        settings.Token = text
        settings.DownloadLocation = "."
    }
    if settings.DownloadLocation == "" {
        settings.DownloadLocation = "."
    }

    url := getEndpointFormatString(Input)
    data := getResponse(
        fmt.Sprintf(url, year, day),
        fmt.Sprintf("session=%s", settings.Token))

    if data != nil {
        if *toStdout {
            fmt.Println(data)
        } else {
            var filename string
            if *filenameFlag == "" {
                filename = fmt.Sprintf("%d-day%d.txt", year, day)
            } else {
                filename = *filenameFlag
            }

            if *filepathFlag != "" {
                downloadFile(data, *filepathFlag, filename)
            } else {
                downloadFile(data, settings.DownloadLocation, filename)
            }
        }
    }
}

func downloadFile(content []byte, location string, filename string) error {
    filepath := fmt.Sprintf("%s/%s", location, filename)

    err := os.WriteFile(filepath, content, 0644)
    if err != nil {
        log.Fatal(err)
    }

    return nil
}

func getEndpointFormatString(method Method) string {
    switch method {
    case Input:
        return "https://adventofcode.com/%d/day/%d/input"
    case Submit:
        log.Fatal("Submit not implemented")
    }
    return ""
}

func getResponse(url string, token string) []byte {
    client := &http.Client {}
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("Cookie", token)

    resp, err := client.Do(req)
    checkErr(err)
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        bytes, err := io.ReadAll(resp.Body)
        checkErr(err)

        return bytes
    }

    return nil
}

func getConfig(settings *Settings) error {
    fileLocations := [2]string{"./config.json"}

    configLocation, err := os.UserConfigDir()
    if err == nil {
        fileLocations[1] = fmt.Sprintf("%s/aoc/config.json", configLocation)
    }

    for location := range fileLocations {
        file, err := os.ReadFile(fileLocations[location])
        if err != nil {
            continue
        }
        s, err := decodePossibleConfig(file)
        if err != nil {
            continue
        }
        *settings = s

        return nil
    }
    return errors.New("No config found")
}

func decodePossibleConfig(content []byte) (Settings, error) {
    var settings Settings
    err := json.Unmarshal(content, &settings)

    if err != nil {
        log.Println(err)
    }
    return settings, err
}

func getNumFromCmd(name string) int {
    reader := bufio.NewReader(os.Stdin)

    fmt.Printf("%s: ", name)
    text, err := reader.ReadString('\n')
    text = strings.Split(text, "\n")[0]
    checkErr(err)

    num, err := strconv.Atoi(text)
    checkErr(err)

    return num
}

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

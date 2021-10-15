package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	curl "github.com/andelf/go-curl"
)

func main() {
	requests := flag.Int("requests", 1, "requests number")
	flag.Parse()
	fmt.Println("requests: ", *requests)

	file, err := os.Open("address")
	if err != nil {
		fmt.Println("Open file failed: ", err.Error())
		return
	}
	defer file.Close()
	lines_byte, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read file failed: ", err.Error())
		return
	}
	lines_str := string(lines_byte)
	lines := strings.Split(lines_str, "\n")
	url0 := lines[0]
	url1 := lines[1]
	url2 := lines[2]
	url3 := lines[3]
	url4 := lines[4]

	coordinator, err := os.OpenFile("coordinator.txt", os.O_WRONLY, 0666)
	defer coordinator.Close()

	start := time.Now()
	for i := 0; i < *requests; i++ {
		mh := curl.MultiInit()
		ch0 := curl.EasyInit()
		ch1 := curl.EasyInit()
		ch2 := curl.EasyInit()
		ch3 := curl.EasyInit()
		ch4 := curl.EasyInit()

		ch0.Setopt(curl.OPT_URL, url0)
		ch0.Setopt(curl.OPT_POSTFIELDS, "a")
		mh.AddHandle(ch0)

		ch1.Setopt(curl.OPT_URL, url1)
		ch1.Setopt(curl.OPT_POSTFIELDS, "a")
		mh.AddHandle(ch1)

		ch2.Setopt(curl.OPT_URL, url2)
		ch2.Setopt(curl.OPT_POSTFIELDS, "a")
		mh.AddHandle(ch2)

		ch3.Setopt(curl.OPT_URL, url3)
		ch3.Setopt(curl.OPT_POSTFIELDS, "a")
		mh.AddHandle(ch3)

		ch4.Setopt(curl.OPT_URL, url4)
		ch4.Setopt(curl.OPT_POSTFIELDS, "a")
		mh.AddHandle(ch4)

		s := time.Now().UnixNano() / 1000000
		running := 1
		for running != 0 {
			running, _ = mh.Perform()
		}
		e := time.Now().UnixNano() / 1000000
		str := strconv.Itoa(i) + " start: " + strconv.FormatInt(s, 10) + " end: " + strconv.FormatInt(e, 10)
		_, err := io.WriteString(coordinator, str)
		if err != nil {
			fmt.Println("write to coordinator failed: ", err.Error())
			return
		}
		fmt.Println(str)
	}
	end := time.Now()

	fmt.Println("start: ", start)
	fmt.Println("end: ", end)
}

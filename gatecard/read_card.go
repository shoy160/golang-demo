package main

import (
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
	ser "go.bug.st/serial"
	"golang.org/x/net/websocket"
)

// 获取串口列表
func getComms() []string {
	ports, err := ser.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("没有找到串口")
	}
	return ports
}

// 串口监听
func comListener(name string, baud int, read func(int)) {
	config := &serial.Config{
		Name:        name,
		Baud:        baud,
		ReadTimeout: time.Second * 5,
	}
	s, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer s.Close()
	log.Printf("connect %s succ\n", config.Name)
	for true {
		buf := make([]byte, 11)
		buffer := make([]byte, 11)
		index := 0
		num, _ := s.Read(buf)
		for num > 0 {
			buffer = append(buffer[:index], buf[:num]...)
			index += num
			if index >= 8 {
				break
			}
			num, _ = s.Read(buf)
		}
		if index >= 8 {
			log.Println(strings.ToUpper(hex.EncodeToString(buffer)))
			//解析卡号
			code := hex.EncodeToString(buffer[1:4])
			n, _ := strconv.ParseUint(code, 16, 32)
			read(int(n))
		}
		index = 0
		buf = buf[0:0]
		buffer = buffer[0:0]
	}
}

var conn *websocket.Conn

// 监听回调
func readCode(code int) {
	log.Println(code)
	if conn == nil || !conn.IsServerConn() {
		return
	}
	var err error
	if err = websocket.Message.Send(conn, strconv.Itoa(code)); err != nil {
		log.Fatal(err)
	}
}

// websocket
func upper(ws *websocket.Conn) {
	conn = ws
	var err error
	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			log.Fatal(err)
			continue
		}

		if err = websocket.Message.Send(ws, strings.ToUpper(reply)); err != nil {
			log.Fatal(err)
			continue
		}
	}
}

// 启动服务
func start() {
	ports := getComms()
	for _, port := range ports {
		go comListener(port, 9600, readCode)
	}
	http.Handle("/gate_card", websocket.Handler(upper))

	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

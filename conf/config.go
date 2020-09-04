package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var properties = make(map[string]string)

func getConfig(key string, def string) string {
	if _, ok := properties[key]; ok {
		return properties[key]
	}
	return def
}

func getConfigAsInt(key string, def int) int {
	if _, ok := properties[key]; ok {
		val, _ := strconv.Atoi(properties[key])
		return val
	}
	return def
}

func init() {
	f, err := ioutil.ReadFile("app.conf")
	if err != nil {
		log.Fatal(err)
		return
	}
	arr := strings.Split(string(f), "\n")
	for _, str := range arr {
		key := strings.Replace(strings.Split(str, "=")[0], " ", "", -1)
		value := strings.Replace(strings.Split(str, "=")[1], " ", "", -1)
		fmt.Printf("%s,%s\n", key, value)
		properties[key] = value
	}
}

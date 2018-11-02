package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
)

func main() {
	var serverBuf bytes.Buffer
	serverWriter := bufio.NewWriter(&serverBuf)

	var watchBuf bytes.Buffer
	watchWriter := bufio.NewWriter(&watchBuf)

	var loadBuf bytes.Buffer
	loadWriter := bufio.NewWriter(&loadBuf)

	if err := getService(serverWriter); err != nil {
		log.Fatal(err)
	}

	if err := getServer(serverWriter); err != nil {
		log.Fatal(err)
	}

	if err := serverWriter.Flush(); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("server.yaml", serverBuf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	if err := getWatch(watchWriter); err != nil {
		log.Fatal(err)
	}

	if err := watchWriter.Flush(); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("watch.yaml", watchBuf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	if err := getLoad(loadWriter); err != nil {
		log.Fatal(err)
	}

	if err := loadWriter.Flush(); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("load.yaml", loadBuf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}

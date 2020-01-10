package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

func main() {
	remoteAddress := flag.String("remote", "", "Remote Address in host:port")
	proxyAddress := flag.String("proxy", "", "HTTPS Proxy Address in host:port")

	flag.Parse()

	conn, err := tls.Dial("tcp", *proxyAddress, nil)
	if err != nil {
		log.Fatalln(err)
	}
	err = conn.Handshake()

	if err != nil {
		log.Fatalln(err)
	}

	requestHead := fmt.Sprintf("CONNECT %s HTTP/1.1\r\n"+
		"UserAgent: SSHOVERHTTPS\r\n"+
		"Proxy-Connection: Keep-Alive\r\n"+
		"Host: %s\r\n\r\n", *remoteAddress, *remoteAddress)

	_, err = conn.Write([]byte(requestHead))
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1500)
	n, err := conn.Read(buffer)
	response := string(buffer[:n])
	response = strings.TrimSpace(response)
	reg := regexp.MustCompile("^HTTP/1\\.1\\ 200")
	if !reg.Match([]byte(response)) {
		log.Fatal("Unconnect err:", response)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		io.Copy(os.Stdout, conn)
		wg.Done()
	}()
	go func() {
		io.Copy(conn, os.Stdin)
		wg.Done()
	}()
	wg.Wait()
}

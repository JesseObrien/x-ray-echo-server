package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-xray-sdk-go/xray"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go startTCPServer(&wg)
	go startUDPServer(&wg)

	wg.Wait()
}

func startTCPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Info("Starting TCP server...")

	var port = ":2000"

	l, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Info("⚡ TCP server lstening on port 2000")

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}
		go handleTCPConnection(c)
	}
}

func handleTCPConnection(c net.Conn) {
	netData, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("TCP Read -> ", string(netData))
}

func startUDPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Info("Starting UDP server...")
	s, err := net.ResolveUDPAddr("udp4", "0.0.0.0:2000")
	if err != nil {
		log.Fatal(err)
		return
	}

	c, err := net.ListenUDP("udp4", s)
	if err != nil {
		log.Fatal(err)
		return
	}
	c.SetReadBuffer(1048576)

	log.Info("⚡ UDP server lstening on port 2000")

	go handleUDPConnection(c)
}

func handleUDPConnection(c *net.UDPConn) {
	defer c.Close()
	for {
		// Largest packet should be 64kB from: https://docs.aws.amazon.com/xray/latest/devguide/xray-api-segmentdocuments.html
		buffer := make([]byte, 65535)
		bytesRead, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Ignore the first chunk of the udp
		str := string(buffer[0:bytesRead])
		chunks := strings.Split(str, "\n")
		if len(chunks) != 2 {
			continue
		}

		payload := chunks[1]

		// Unmarshal the segment so we can figure out if we need to log it or not
		seg := &xray.Segment{}
		if err := json.Unmarshal([]byte(payload), seg); err != nil {
			log.Warn(err)
			continue
		}

		switch seg.Name {
		case "dns":
			fallthrough
		case "dial":
			fallthrough
		case "connect":
			fallthrough
		case "response":
			continue
		default:
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(payload), "", "  ")
			if err != nil {
				log.Warn(err)
				continue
			}
			log.Info("Xray Segment Received")
			log.Infof("%v", string(prettyJSON.Bytes()))
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"net"
)

// 글로벌 변수 : 모든 피어의 연결 정보 저장하는 리스트
var connectedPeers []net.Conn

func main() {
	// 명령줄
	mode := flag.String("mode", "fullNode", "Start in 'Bootstrap Node' or 'FullNode' ")
	port := flag.Int("port", 30303, "The port on which the erver listen (TCP & UDP)")
	// 명령줄 인자 파싱 (flag.Parse() 필수)
	flag.Parse()

	tcpAddress := make(chan string)
	udpAddress := make(chan string)
	bootstrapAddress := "localhost:8282"

	if *mode == "bootstrap" {
		startBootstrapServer()

	} else if *mode == "fullNode" {
		// 1. TCP 서버 실행
		go startTCPServer(tcpAddress, *port)

		// 2. UDP 서버 실행
		go startUDPServer(udpAddress, tcpAddress, *port)

		// 3. UDP 서버 주소 받아옴
		udpServerAddress := <-udpAddress

		// 3. 부트스트랩 노드에 연결하고, 내 서버 정보 전달
		nodeAddress := connectBootstrapNode(bootstrapAddress, udpServerAddress)
		fmt.Println("부트스트랩 노드로부터 받은 노드들 주소 :", nodeAddress)

		// 4. 받은 노드들과 연결 시도
		startClient(nodeAddress)
	} else {
		fmt.Println("Invalid mode. Use -mode=bootstrap or -mode=fullNode")
	}

}

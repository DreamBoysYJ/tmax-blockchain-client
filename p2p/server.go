package p2p

import (
	"bufio"
	"fmt"
	"io"
	"net"
	pc "simple_p2p_client/protocol_constants"
	"simple_p2p_client/utils"
	"strconv"
	"strings"
)

// TCP 서버 실행
func StartTCPServer(tcpAddress chan<- string, port int) {

	// TCP 서버 시작
	address := "127.0.0.1:" + strconv.Itoa(port)
	// min := 6666
	// max := 9999
	// randomPort := rand.Intn(max-min+1) + min

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		utils.PrintError(fmt.Sprintf("Error starting server: %v", err))
		return
	}

	fmt.Println("TCP Server is listening on :", address)
	tcpAddress <- address

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.PrintError(fmt.Sprintf("Error accepting connection: %v", err))
			continue
		}
		utils.PrintMessage(fmt.Sprintf("Peer connected from: %v", conn.RemoteAddr().String()))

		// 연결된 피어를 글로벌 변수에 저장
		ConnectedPeers = append(ConnectedPeers, conn)

		// 피어와 통신 처리
		go HandleIncomingMessages(conn)
	}
}

// UDP 서버 실행
func StartUDPServer(udpAddress chan<- string, tcpAddress <-chan string, port int) {

	// 1. Create UDP address
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("127.0.0.1"),
	}

	// 2. Waiting for UDP connection
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting UDP server :", err)
		return
	}

	fmt.Println("UDP Server is listening on : ", fmt.Sprintf("%s:%d", addr.IP.String(), addr.Port))

	udpAddress <- addr.String()

	defer conn.Close()

	// TCP 주소 받기 전까지 대기
	tcpListeningAddr := <-tcpAddress
	fmt.Println("Recevied TCP address for UDP handling : ", tcpListeningAddr)

	// create Buffer
	buf := make([]byte, 1024)

	// Listening UDP Server
	// Node Discovery Protocol
	for {
		// n : 바이트 수 (일단 생략)
		_, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receving UDP data : ", err)
			continue
		}

		messageType := buf[0]

		switch messageType {

		// Node가 payload에 Ping(0x01)을 보냈을 시
		case pc.NodeDiscoveryPing:
			fmt.Printf("Received Ping message from %s\n", remoteAddr.String())
			// Pong으로 응답
			_, err = conn.WriteToUDP([]byte{pc.NodeDiscoveryPong}, remoteAddr)
			if err != nil {
				fmt.Println("Error sending UDP response : ", err)
			}
			fmt.Printf("Sent Pong to %s\n", remoteAddr.String())
		// Node가 payload에 ENRRequest(0x05)를 보냈을 시
		case pc.NodeDiscoveryENRRequest:
			fmt.Printf("Received ENRRequest from %s\n", remoteAddr.String())
			// ENRResponse(0x06)로 응답
			message := append([]byte{pc.NodeDiscoveryENRResponse}, []byte(tcpListeningAddr)...)
			_, err := conn.WriteToUDP(message, remoteAddr)
			if err != nil {
				fmt.Println("Error sending UDP message : ", err)
			}
			fmt.Printf("Send ENRResponse to %s\n", remoteAddr.String())

		}

	}
}

// 받은 메시지 처리 함수
func HandleIncomingMessages(conn net.Conn) {
	defer conn.Close()

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				utils.PrintError(fmt.Sprintf("Peer disconnected: %s", conn.RemoteAddr().String()))
				ConnectedPeers = utils.RemoveConn(ConnectedPeers, conn)

			} else {
				utils.PrintError(fmt.Sprintf("Error reading from peer: %v", err))
			}
			return
		}

		message = strings.TrimSpace(message)
		utils.PrintMessage(fmt.Sprintf("Message received from peer : %s from : %s", message, conn.RemoteAddr().String()))
	}
}

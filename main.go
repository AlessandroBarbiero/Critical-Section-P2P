package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	token "github.com/AlessandroBarbiero/Critical-Section-P2P/grpc"
	"google.golang.org/grpc"
)

// Peer data structure, the mutex is used to assure that the token is passed only after the
// next peer has been found and the connection established
type peer struct {
	token.UnimplementedTokenServer
	id           int32
	nextPeer     token.TokenClient
	nextPeerPort int32
	request      bool
	mutex        sync.RWMutex
	ctx          context.Context
}

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	arg2, _ := strconv.ParseInt(os.Args[2], 10, 32)
	ownPort := int32(arg1)
	// node that should start sending token
	firstNodePort := int32(arg2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:      ownPort,
		request: false,
		ctx:     ctx,
		mutex:   sync.RWMutex{},
	}

	//set log file fo
	f, err := os.OpenFile("network.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	//set logger prefix
	log.SetPrefix(fmt.Sprintf("Node %v: ", p.id))

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	token.RegisterTokenServer(grpcServer, p)

	// Waiting for other peers to connect to me on another thread
	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("Failed server function at P %v: %v", p.id, err)
		}
		log.Printf("Server %v has started\n", p.id)
	}()

	// Find my next peer
	p.readConfigFile()

	//Create the connection with the next peer
	conn := p.dialNextPeer()
	defer conn.Close()

	p.mutex.Lock()
	p.nextPeer = token.NewTokenClient(conn)
	p.mutex.Unlock()

	go func() {
		if p.id == firstNodePort {
			request := &token.Request{}
			p.nextPeer.Token(ctx, request)
		}
	}()

	fmt.Printf("Hi I am node %v\n", ownPort)

	// Take input and wait for the token to actually write in the restricted area
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if p.request {
			log.Println("Wait, we are processing the previous request")
		} else {
			log.Println("Request accepted, waiting for the token to write in restricted area")
			p.request = true
		}
	}
}

func (p *peer) Token(ctx context.Context, req *token.Request) (*token.Reply, error) {
	// if there is a request to process
	// access the critical area
	if p.request {
		log.Println("Got token, and I need it")
		p.criticalArea()
		p.request = false
	}
	// send token to next node
	p.giveTokenToNextPeer()
	rep := &token.Reply{}
	return rep, nil
}

// If the nextPeer doesn't exist yet wait for his instantiation
func (p *peer) giveTokenToNextPeer() {
	request := &token.Request{}

	// check if connection is ready
	p.mutex.RLock()
	np := p.nextPeer
	p.mutex.RUnlock()
	for np == nil {
		time.Sleep(time.Second * 1)
		p.mutex.RLock()
		np = p.nextPeer
		p.mutex.RUnlock()
	}

	p.nextPeer.Token(p.ctx, request)
}

func (p *peer) dialNextPeer() *grpc.ClientConn {
	var conn *grpc.ClientConn
	log.Printf("Trying to dial: %v\n", p.nextPeerPort)
	conn, err := grpc.Dial(fmt.Sprintf(":%v", p.nextPeerPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect to %v: %s", p.nextPeerPort, err)
	}
	log.Printf("Connected to %v", p.nextPeerPort)
	return conn
}

// Read config file to save port of next peer
func (p *peer) readConfigFile() {
	name := "clients.info"
	file, err := os.Open(name)
	if err != nil {
		log.Fatalln("Couldn't read file with ports")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	foundMyPort := false

	// read config file
	// wait for finding my port
	// when my port is found read next line, store it as a nextPeerPort and break the loop
	for scanner.Scan() {
		id, e := strconv.Atoi(scanner.Text())
		if e != nil {
			log.Fatalln("Invalid value in config file")
		}
		if foundMyPort {
			p.nextPeerPort = int32(id)
			return
		}
		if int32(id) == p.id {
			foundMyPort = true
		}
	}

	// when my port is last between all ports seek go back to the beginning of file
	// and read and store first port in file as nextPeerPort
	if foundMyPort {
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatalln("Couldn't seek start of file, while reading config file")
		}
		scanner = bufio.NewScanner(file)
		scanner.Scan()
		id, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatalln("Invalid value in config file")
		}
		p.nextPeerPort = int32(id)
	}
}

// Critical area that should be accessible by only one node at a time, to show the correctness of the algorithm
// a delay is added to simulate a heavy utilization of the restricted area, in this time only one node can access it
func (p *peer) criticalArea() {
	log.Printf("Critical area was reached by node %v\n", p.id)

	//Add a delay for demonstration purpose before writing in the file
	time.Sleep(10 * time.Second)

	// open criticalArea file for append
	name := "criticalArea.log"
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Couldn't read file with ports")
	}
	defer file.Close()

	_, errw := file.WriteString(fmt.Sprintf("Critical area was reached by node %v\n", p.id))
	if errw != nil {
		log.Fatalf("Couldn't write to file, error: %v\n", errw)
	}
}

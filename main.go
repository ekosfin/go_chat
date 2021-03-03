package main

import (
	"log"
	"net"
)

func main() {

	//Creates the servet that processes the requests
	s := newServer()
	go s.run()

	//Starts to listen to the tcp port
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}
	//This will exicute as the main function ends
	defer listener.Close()
	log.Printf("server started on :8888")

	//infite while loop where listener.Accept() is a blocking line
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}
		//Starts a thread for each client
		go s.newClient(conn)
	}
}

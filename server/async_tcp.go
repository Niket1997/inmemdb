package server

import (
	"fmt"
	"github.com/Niket1997/inmemdb/config"
	"github.com/Niket1997/inmemdb/core"
	"io"
	"log"
	"net"
	"strings"
	"syscall"
)

var connClients = 0

func readCommand(c io.ReadWriter) (*core.RedisCmd, error) {
	// TODO: Max read in one shot is 512 bytes
	// To allow input > 512 bytes, then repeated read until
	// we get EOF or designated delimiter
	var buf []byte = make([]byte, 512)
	// Read is blocking command, i.e. the code will wait until there
	// is something to read from the socket connection
	n, err := c.Read(buf[:])
	if err != nil {
		return nil, err
	}

	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}
	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, c io.ReadWriter) {
	_, err2 := c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
	if err2 != nil {
		log.Println(err2)
	}
}

func respond(cmd *core.RedisCmd, c io.ReadWriter) {
	err := core.EvalAndRespond(cmd, c)
	if err != nil {
		respondError(err, c)
	}
}

func RunAsyncTCPServer() error {
	log.Println("starting a asynchronous TCP server on", config.Host, config.Port)

	maxClients := 20000

	// create EPOLL event objects to hold events
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, maxClients)

	// Create a socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD)

	// Set the Socket operate in a non-blocking mode
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	// Bind the ip & the port
	ip4 := net.ParseIP(config.Host)
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	// Start listening
	if err = syscall.Listen(serverFD, maxClients); err != nil {
		return err
	}

	// AsyncIO starts here!!

	// create EPOLL instance
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}
	defer syscall.Close(epollFD)

	// specify the events we want to get callbacks about
	var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}

	// listen to read events on server itself
	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerEvent); err != nil {
		return err
	}

	for {
		// see if any FD is ready for IO
		nevents, err := syscall.EpollWait(epollFD, events[:], -1)
		if err != nil {
			continue
		}

		for i := 0; i < nevents; i++ {
			// if the socket server itself is ready for an IO
			// i.e. new client wants to connect to the server
			if int(events[i].Fd) == serverFD {
				// accept the incoming connection from client
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("err", err)
					continue
				}

				// increase the number of concurrent clients
				connClients++
				syscall.SetNonblock(serverFD, true)

				// add this TCP connection to be monitored
				var socketClientEvent syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(fd),
				}

				if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvent); err != nil {
					log.Fatal(err)
				}
			} else {
				comm := core.FDComm{Fd: int(events[i].Fd)}
				cmd, err := readCommand(comm)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					connClients--
					continue
				}
				respond(cmd, comm)
			}
		}
	}
}

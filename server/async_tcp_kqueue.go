package server

//func RunAsyncTCPServerUnix() error {
//	log.Println("starting an asynchronous TCP server on", config.Host, config.Port)
//
//	maxClients := 20000
//	connClients := 0
//
//	// Create kqueue event objects to hold events
//	events := make([]unix.Kevent_t, maxClients)
//
//	// Create a socket
//	serverFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_TCP)
//	if err != nil {
//		return err
//	}
//	defer unix.Close(serverFD)
//
//	// Set the socket to operate in non-blocking mode
//	if err = unix.SetNonblock(serverFD, true); err != nil {
//		return err
//	}
//
//	// Allow address reuse
//	if err = unix.SetsockoptInt(serverFD, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
//		return err
//	}
//
//	// Bind the IP & the port
//	ip4 := net.ParseIP(config.Host).To4()
//	if ip4 == nil {
//		return fmt.Errorf("invalid IPv4 address")
//	}
//	sa := &unix.SockaddrInet4{
//		Port: config.Port,
//	}
//	copy(sa.Addr[:], ip4)
//	if err = unix.Bind(serverFD, sa); err != nil {
//		return err
//	}
//
//	// Start listening
//	if err = unix.Listen(serverFD, maxClients); err != nil {
//		return err
//	}
//
//	// Async I/O starts here!!
//
//	// Create kqueue instance
//	kq, err := unix.Kqueue()
//	if err != nil {
//		return err
//	}
//	defer unix.Close(kq)
//
//	// Specify the events we want to get callbacks about
//	kev := unix.Kevent_t{
//		Ident:  uint64(serverFD),
//		Filter: unix.EVFILT_READ,
//		Flags:  unix.EV_ADD,
//	}
//
//	// Register the serverFD with kqueue
//	_, err = unix.Kevent(kq, []unix.Kevent_t{kev}, nil, nil)
//	if err != nil {
//		return err
//	}
//
//	for {
//		// See if any FD is ready for I/O
//		nevents, err := unix.Kevent(kq, nil, events, nil)
//		if err != nil {
//			continue
//		}
//
//		for i := 0; i < nevents; i++ {
//			ev := events[i]
//			fd := int(ev.Ident)
//
//			// If the socket server itself is ready for an I/O
//			// i.e., new client wants to connect to the server
//			if fd == serverFD {
//				// Accept the incoming connection from client
//				nfd, _, err := unix.Accept(serverFD)
//				if err != nil {
//					log.Println("err", err)
//					continue
//				}
//
//				// Increase the number of concurrent clients
//				connClients++
//				if err := unix.SetNonblock(nfd, true); err != nil {
//					log.Println("err", err)
//					unix.Close(nfd)
//					continue
//				}
//
//				// Add this TCP connection to be monitored
//				kev := unix.Kevent_t{
//					Ident:  uint64(nfd),
//					Filter: unix.EVFILT_READ,
//					Flags:  unix.EV_ADD,
//				}
//
//				_, err = unix.Kevent(kq, []unix.Kevent_t{kev}, nil, nil)
//				if err != nil {
//					log.Fatal(err)
//				}
//			} else {
//				// Handle client I/O
//				comm := core.FDComm{Fd: fd}
//				cmd, err := readCommand(comm)
//				if err != nil {
//					// Remove the socket from kqueue
//					kev := unix.Kevent_t{
//						Ident:  uint64(fd),
//						Filter: unix.EVFILT_READ,
//						Flags:  unix.EV_DELETE,
//					}
//					unix.Kevent(kq, []unix.Kevent_t{kev}, nil, nil)
//					unix.Close(fd)
//					connClients--
//					continue
//				}
//				respond(cmd, comm)
//			}
//		}
//	}
//}

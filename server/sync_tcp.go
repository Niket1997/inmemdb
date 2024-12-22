package server

//func RunSyncTCPServer() {
//	log.Println("starting a synchronous TCP server on", config.Host, config.Port)
//
//	var conClients int = 0
//
//	// listening to the configured hots:port
//	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
//	if err != nil {
//		panic(err)
//	}
//
//	for {
//		// blocking call: waiting for the new client to connect
//		c, err := lsnr.Accept()
//		if err != nil {
//			panic(err)
//		}
//
//		// increment number of concurrent clients
//		conClients++
//		log.Println("client connected with address:", c.RemoteAddr(), "concurrent clients", conClients)
//
//		for {
//			// over the socket, continuously read the command & print it
//			cmd, err := readCommand(c)
//			if err != nil {
//				c.Close()
//				conClients--
//				log.Println("client disconnected", c.RemoteAddr(), "concurrent clients", conClients)
//				if err != io.EOF {
//					log.Println("err", err)
//				}
//				break
//			}
//			respond(cmd, c)
//		}
//	}
//}

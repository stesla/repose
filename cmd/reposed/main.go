package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var listener net.Listener
var connections = map[string]net.Conn{}

var graceful = flag.String("graceful", "", "get listener from parent")

func main() {
	var err error

	log.Println("pid =", os.Getpid())

	flag.Parse()

	if *graceful != "" {
		fds := strings.Split(*graceful, ",")
		lfd, err := strconv.Atoi(fds[0])
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("old listener, fd =", lfd)
		listener, err = net.FileListener(os.NewFile(uintptr(lfd), "listener"))
		if err != nil {
			log.Fatalln(err)
		}

		for _, str := range fds[1:] {
			cfd, err := strconv.Atoi(str)
			if err != nil {
				log.Fatalln(err)
			}
			conn, err := net.FileConn(os.NewFile(uintptr(cfd), "connection"))
			fmt.Fprintf(conn, "...Restarted\n")
			go clientSession(conn)
		}
	} else {
		listener, err = net.Listen("tcp", ":2000")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("new listener")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go clientSession(conn)
	}

	listener.Close()
}

func clientSession(c net.Conn) {
	connections[c.RemoteAddr().String()] = c
	defer func() {
		c.Close()
		delete(connections, c.RemoteAddr().String())
		log.Printf("connection to %v closed", c.RemoteAddr())
	}()

	s := bufio.NewScanner(c)
	for s.Scan() {
		switch t := s.Text(); t {
		case "close", "quit":
			return
		case "who":
			for _, conn := range connections {
				fmt.Fprintf(c, "%v\n", conn.RemoteAddr())
			}
		case "restart":
			doRestart()
		case "shutdown":
			os.Exit(0)
		default:
			fmt.Fprintln(c, "YOU SAID:", t)
		}
	}
}

func doRestart() {
	file, _ := listener.(*net.TCPListener).File()
	files := []*os.File{file}
	for _, conn := range connections {
		fmt.Fprintf(conn, "Restarting...\n")
		file, _ := conn.(*net.TCPConn).File()
		files = append(files, file)
	}
	fds := make([]string, len(files))
	for i, file := range files {
		syscall.Syscall(syscall.SYS_FCNTL, file.Fd(), syscall.F_SETFD, 0)
		fds[i] = fmt.Sprintf("%v", file.Fd())
	}
	path := os.Args[0]
	args := []string{path, "-graceful", strings.Join(fds, ",")}
	syscall.Exec(path, args, os.Environ())
}

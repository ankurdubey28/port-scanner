package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type config struct {
	hosts []string
	port  int
}

var ErrHostNotSpecified = errors.New("hosts not specified")
var ErrInvalidPort = errors.New("port out of range [1,65535]")
var ErrUnexpectedPosArgs = errors.New("received positional args while expected none")

func main() {
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	err = runCmd(c)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
}

func parseArgs(w io.Writer, args []string) (*config, error) {
	c := config{}
	fs := flag.NewFlagSet("port scanner", flag.ContinueOnError)
	fs.SetOutput(w)
	var hostStr string
	fs.StringVar(&hostStr, "hosts", "", "hosts name / names")
	fs.IntVar(&c.port, "port", -1, "give target port")

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}
	parsedArgs := fs.Args()
	if len(parsedArgs) > 0 {
		return nil, ErrUnexpectedPosArgs
	}
	if hostStr != "" {
		c.hosts = strings.Split(hostStr, ",")
	}
	if c.port != -1 && (c.port < 1 || c.port > 65535) {
		return nil, ErrInvalidPort
	}
	return &c, nil
}

func runCmd(c *config) error {
	// given single host and port , find if that port is open or not
	if len(c.hosts) == 1 && c.port != -1 {
		err := connect(c.hosts[0], c.port)
		if err != nil {
			return err
		}
		// run scan on single/multiple host to find any open port
	} else if len(c.hosts) >= 1 {
		for _, host := range c.hosts {
			err := scan(host)
			if err != nil {
				return err
			}
		}
	} else {
		return ErrHostNotSpecified
	}
	return nil
}

func connect(host string, port int) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, 1000*time.Millisecond)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func scan(host string) error {
	var wg sync.WaitGroup
	activePorts := make(chan int, 20)

	for i := 1; i <= 65535; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			err := connect(host, port)
			if err == nil {
				activePorts <- port
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(activePorts)
	}()

	for port := range activePorts {
		fmt.Printf("port %d is active\n", port)
	}
	return nil
}

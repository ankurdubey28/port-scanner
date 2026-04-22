package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type config struct {
	host string
	port int
}

var ErrHostNotSpecified = errors.New("host not specified")
var ErrInvalidPort = errors.New("port out of range [1,65535]")

func main() {
	c, err := parseArgs(os.Stdout, os.Args[1:])
	if err != nil {
		fmt.Fprint(os.Stdout, err)
		return
	}
	err = runCmd(c)
	if err != nil {
		fmt.Fprint(os.Stdout, err)
		return
	}
}

func parseArgs(w io.Writer, args []string) (*config, error) {
	c := config{}
	fs := flag.NewFlagSet("port scanner", flag.ContinueOnError)
	fs.SetOutput(w)

	fs.StringVar(&c.host, "host", "", "give host name")
	fs.IntVar(&c.port, "port", -1, "give target port")

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}
	parsedArgs := fs.Args()
	if len(parsedArgs) > 0 {
		return nil, err
	}
	if c.port != -1 && (c.port < 1 || c.port > 65535) {
		return nil, ErrInvalidPort
	}
	return &c, nil
}

func runCmd(c *config) error {
	if c.host != "" && c.port != -1 {
		err := connect(c.host, c.port)
		if err != nil {
			return err
		}
	} else if c.host != "" {
		err := scan(c.host)
		if err != nil {
			return err
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
	activePorts := make(chan int, 2)

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

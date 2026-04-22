package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type config struct {
	host string
	port int
}

var ErrHostNotSpecified = errors.New("host not specified")

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
	fmt.Println("connection established")
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
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func scan(host string) error {
	for i := 1; i <= 65535; i++ {
		err := connect(host, i)
		if err != nil {
			continue
		}
		fmt.Printf("port %d is open\n", i)
	}
	return nil
}

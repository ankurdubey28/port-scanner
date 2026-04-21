package main

import (
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
	fs.IntVar(&c.port, "port", 00, "give target port")

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
	addr := net.JoinHostPort(c.host, strconv.Itoa(c.port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

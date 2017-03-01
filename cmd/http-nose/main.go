package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"github.com/buglloc/http-nose/httpclient"
	"github.com/buglloc/http-nose/httpfeature"
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] target-host:port\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	hostFlag := flag.String("host", "localhost", "request host")
	pathFlag := flag.String("path", "/", "request path")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	baseRequest := httpclient.Request{
		Method:     "GET",
		RequestURI: *pathFlag,
		Proto:      "HTTP/1.0",
		Headers: []httpclient.Header{
			{Name: "Host", Value: *hostFlag},
		},
	}

	client := httpclient.Client{Target: args[0]}
	baseResponse, err := client.MakeRequest(&baseRequest)
	if err != nil {
		log.Fatal("Failed to make base request: ", err)
	}

	httpfeatures := httpfeature.NewFeatures(client, baseRequest, *baseResponse)
	for _, f := range httpfeatures.Collect() {
		fmt.Printf("%s: %s\n", f.Name(), f.ToString())
	}
}

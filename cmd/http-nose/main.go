package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/buglloc/http-nose/httpclient"
	"github.com/buglloc/http-nose/httpfeature"
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] target-host:port\n", os.Args[0])
	flag.PrintDefaults()
}

func PrepareFeaturesForJson(features *httpfeature.Features) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, f := range features.Collect() {
		result = append(result,
			map[string]interface{}{
				"Name":  f.Name(),
				"Value": f.Export(),
			},
		)
	}
	return result
}

func main() {
	flag.Usage = usage
	hostFlag := flag.String("host", "localhost", "request host")
	pathFlag := flag.String("path", "/", "request path")
	argsFlag := flag.String("args", "", "args")
	formatFlag := flag.String("format", "text", "output format")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	baseRequest := httpclient.Request{
		Method: "GET",
		Path:   *pathFlag,
		Args:   *argsFlag,
		Proto:  "HTTP/1.0",
		Host:   *hostFlag,
	}

	client := httpclient.Client{Target: args[0]}
	baseResponse, err := client.MakeRequest(&baseRequest)
	if err != nil {
		log.Fatal("Failed to make base request: ", err)
	}

	httpfeatures := httpfeature.NewFeatures(client, baseRequest, *baseResponse)
	switch *formatFlag {
	case "text":
		for _, f := range httpfeatures.Collect() {
			fmt.Printf("%s: %s\n", f.Name(), f.String())
		}
	case "json":
		preparedFeatures := PrepareFeaturesForJson(httpfeatures)
		formatted, err := json.Marshal(preparedFeatures)
		if err != nil {
			log.Fatal("Failed to encode: ", err)
		}
		fmt.Println(string(formatted))
	}
}

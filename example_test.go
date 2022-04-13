package boa_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/feliixx/boa"
)

func Example() {

	config := `
	{
		"http_server": {
			"enabled": true,
			"host": "127.0.0.1"
		}
	}`

	boa.SetDefault("http_server.port", 80)

	err := boa.ParseConfig(strings.NewReader(config))
	if err != nil {
		log.Fatal(err)
	}

	srvHost := boa.GetString("http_server.host")
	srvPort := boa.GetInt("http_server.port")

	fmt.Printf("%s:%d", srvHost, srvPort)
	// Output: 127.0.0.1:80
}

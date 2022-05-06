[![Go Report Card](https://goreportcard.com/badge/github.com/feliixx/boa)](https://goreportcard.com/report/github.com/feliixx/boa)
[![codecov](https://codecov.io/gh/feliixx/boa/branch/master/graph/badge.svg)](https://codecov.io/gh/feliixx/boa)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/feliixx/boa)](https://pkg.go.dev/github.com/feliixx/boa)

## boa

A small configuration library for go application with a viper-like API, but with limited scope and no external dependency 

It supports: 

  * reading a config in JSON or JSONC ( JSON with comments) 
  * setting default 



## example


```go
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
```


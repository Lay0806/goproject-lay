package main

import (
	"myproject/httpproxy"
)

func main() {
	go httpproxy.HttpProxy()
}

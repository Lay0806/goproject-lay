package httpproxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

func HttpProxy() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	_, err := client.Read(b[:])
	if err != nil {
		log.Panic(err)
		return
	}
	var method, host, address string
	fmt.Scanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Panic(err)
		return
	}
	if hostPortURL.Opaque == "443" { // https 访问
		address = hostPortURL.Host + ":443"
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 {
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}

	//获取了请求的hos和port，就开始拨号了
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Panic(err)
		return
	}
	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connecction established\r\n\r\n")
	} else {
		server.Write(b[:])
	}

	//进行转发
	go io.Copy(server, client)
	io.Copy(client, server)
}

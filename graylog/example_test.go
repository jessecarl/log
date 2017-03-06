package graylog_test

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net"
	"os"

	"github.com/PermissionData/log"
	"github.com/PermissionData/log/graylog"
)

func ExampleNew() {
	server, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 9090})
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	defer close(done)
	go func() {
		defer server.Close()
		for {
			select {
			case <-done:
				return
			default:
			}
			b := make([]byte, 256)
			_, err := server.Read(b)
			if err != nil && err != io.EOF {
				panic(err)
			}
			os.Stdout.Write(b)
		}
	}()
	// END fake server implementation

	conn, err := net.ListenPacket("udp4", "0.0.0.0:9091")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	gw, err := graylog.New(graylog.Config{
		ClientPacketConn: conn,
		ServerAddr:       server.LocalAddr(),
		CompressionLevel: gzip.DefaultCompression,
	})
	if err != nil {
		panic(err)
	}
	logger := log.New(log.Config{
		Threshold: log.ErrorLevel,
		Encoder:   json.NewEncoder(gw),
		Filters: []log.Filter{
			log.DefaultFilter,
			func(lvl, threshold log.Level, data log.Data) log.Data {
				data["hey"] = &struct{ Ho bool }{true}
				return data
			},
			func(lvl, threshold log.Level, data log.Data) log.Data {
				if data == nil {
					return nil
				}
				data["@timestamp"] = nil
				return data
			},
		},
	})

	logger.Log(log.InfoLevel, log.Data{})
}

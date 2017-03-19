package graylog_test

import (
	"compress/gzip"
	"io"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/PermissionData/log/graylog"
)

func BenchmarkWrite(b *testing.B) {
	// For the sake of benchmarking, we'll use a NoCompression level of compression
	compressionLvl := gzip.NoCompression

	b.Run("Size=10", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			source := &randomWriterTo{
				size: 10,
				r:    rand.New(rand.NewSource(10)),
			}
			for pb.Next() {
				source.WriteTo(c)
			}
		})
	})

	b.Run("Size=100", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			source := &randomWriterTo{
				size: 100,
				r:    rand.New(rand.NewSource(10)),
			}
			for pb.Next() {
				source.WriteTo(c)
			}
		})
	})

	b.Run("Size=1000", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			source := &randomWriterTo{
				size: 1000,
				r:    rand.New(rand.NewSource(10)),
			}
			for pb.Next() {
				source.WriteTo(c)
			}
		})
	})

	b.Run("Size=10000", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			source := &randomWriterTo{
				size: 10000,
				r:    rand.New(rand.NewSource(10)),
			}
			for pb.Next() {
				source.WriteTo(c)
			}
		})
	})

	b.Run("Size=100000", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			source := &randomWriterTo{
				size: 100000,
				r:    rand.New(rand.NewSource(10)),
			}
			for pb.Next() {
				source.WriteTo(c)
			}
		})
	})
}

type randomWriterTo struct {
	size int

	mu sync.Mutex
	r  *rand.Rand
}

func (wt *randomWriterTo) WriteTo(w io.Writer) (int64, error) {
	b := make([]byte, wt.size+1)
	wt.mu.Lock()
	wt.r.Read(b[:wt.size])
	wt.mu.Unlock()
	b[wt.size] = '\n'
	n, err := w.Write(b)
	return int64(n), err
}

type nopPacketConn struct{}

func (nopPacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	return 0, nil, nil
}

func (nopPacketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	return 0, nil
}

func (nopPacketConn) Close() error {
	return nil
}

func (nopPacketConn) LocalAddr() net.Addr {
	return nil
}

func (nopPacketConn) SetDeadline(t time.Time) error {
	return nil
}

func (nopPacketConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (nopPacketConn) SetWriteDeadline(t time.Time) error {
	return nil
}

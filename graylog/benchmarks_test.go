package graylog_test

import (
	"compress/gzip"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/PermissionData/log/graylog"
)

func BenchmarkWrite(b *testing.B) {
	// For the sake of benchmarking, we'll use a NoCompression level of compression
	compressionLvl := gzip.NoCompression

	// returns a newline terminated string of length 20+size*3 (not including the newline)
	stringToTest := func(size int) string {
		return fmt.Sprintf("{\"long string\" : \"%s\"}\n", strings.Repeat("pop ", size))
	}

	b.Run("Size=1", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(1)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=2", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(2)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=3", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(3)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=5", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(5)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=8", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(8)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=13", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(13)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=21", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(21)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=34", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(34)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=55", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(55)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=89", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(89)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=144", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(144)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=233", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(233)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=377", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(377)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=610", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(610)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=987", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(987)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=1597", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(1597)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=2584", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(2584)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=4181", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(4181)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=6765", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(6765)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=10946", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(10946)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=17711", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(17711)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
	b.Run("Size=28657", func(b *testing.B) {
		c, _ := graylog.New(graylog.Config{
			CompressionLevel: compressionLvl,
			ClientPacketConn: nopPacketConn{},
		})
		b.RunParallel(func(pb *testing.PB) {
			s := stringToTest(28657)
			for pb.Next() {
				c.Write([]byte(s))
			}
		})
	})
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

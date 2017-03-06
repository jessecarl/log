package graylog_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	mock_net "git.permissiondata.com/libs/log/graylog/mock"

	"github.com/PermissionData/log/graylog"
)

func TestWriteImplementsWriter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAddr := mock_net.NewMockAddr(mockCtrl)

	mockPacketConn := mock_net.NewMockPacketConn(mockCtrl)
	mockPacketConn.EXPECT().WriteTo(gomock.Any(), gomock.Any()).Times(0)

	var (
		w   io.Writer
		err error
	)
	w, err = graylog.New(graylog.Config{
		ClientPacketConn: mockPacketConn,
		ServerAddr:       mockAddr,
	})
	if err != nil {
		t.Fatalf("Unexpected error for zero-value configuration: %+v", err)
	}
	p := []byte("some bytes to write")
	p2 := []byte("some bytes to write")

	n, err := w.Write(p)
	if n > len(p) {
		t.Fatalf("Write must not write more bytes than the length of the input")
	}
	if n < len(p) && err == nil {
		t.Fatalf("Write must return a non-nil error if it returns n < len(p)")
	}
	if !reflect.DeepEqual(p, p2) {
		t.Fatalf("Write cannot modify the input slice")
	}
}

func TestWriteRequiresNewline(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAddr := mock_net.NewMockAddr(mockCtrl)
	mockPacketConn := mock_net.NewMockPacketConn(mockCtrl)

	testCases := []struct {
		name          string
		p             []byte
		expectIsError bool
	}{
		{"just a newline",
			[]byte("\n"),
			false,
		},
		{"valid JSON without newline",
			[]byte("{\"something\":3.14}"),
			true,
		},
		{"valid JSON wit newline",
			[]byte("{\"something\":3.14}\n"),
			false,
		},
		{"gibberish without newline",
			[]byte("random gibberish and stuff..~!@#$%^&*("),
			true,
		},
		{"gibberish with newline",
			[]byte("random gibberish and stuff..~!@#$%^&*(\n"),
			false,
		},
	}

	w, _ := graylog.New(graylog.Config{
		ClientPacketConn: mockPacketConn,
		ServerAddr:       mockAddr,
	})
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expectIsError {
				mockPacketConn.EXPECT().WriteTo(gomock.Any(), gomock.Eq(mockAddr)).Times(1)
			}
			_, err := w.Write(tc.p)
			if (err != nil) != tc.expectIsError {
				t.Fatalf("w.Write(%q) = _, %+v, expected error? %v", tc.p, err, tc.expectIsError)
			}
		})
	}
}

func TestWriteIgnoresEmptyBytes(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAddr := mock_net.NewMockAddr(mockCtrl)
	mockPacketConn := mock_net.NewMockPacketConn(mockCtrl)
	mockPacketConn.EXPECT().WriteTo(gomock.Any(), gomock.Any()).Times(0)

	w, _ := graylog.New(graylog.Config{
		ClientPacketConn: mockPacketConn,
		ServerAddr:       mockAddr,
	})
	if n, err := w.Write(nil); n != 0 || err != nil {
		t.Fatalf("w.Write(<nil>) = %d, %+v, expected 0, <nil>", n, err)
	}
	if n, err := w.Write([]byte{}); n != 0 || err != nil {
		t.Fatalf("w.Write(<nil>) = %d, %+v, expected 0, <nil>", n, err)
	}
	if n, err := w.Write([]byte("")); n != 0 || err != nil {
		t.Fatalf("w.Write(<nil>) = %d, %+v, expected 0, <nil>", n, err)
	}
}

func TestWriteAddsHeaders(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
	}{
		{"single message",
			[]byte("{\"short\":{\"pi\":3.14,\"phi\":1.618}}\n"),
		},
		{"multiple messages",
			[]byte(fmt.Sprintf("{\"long string\":\"%s\"}\n", strings.Repeat("pop ", 10000))),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockPacketConn := mock_net.NewMockPacketConn(mockCtrl)

			w, err := graylog.New(graylog.Config{
				ClientPacketConn: mockPacketConn,
			})
			if err != nil {
				t.Fatalf("error constructing New Client: %+v", err)
			}

			var msgID []byte
			var (
				i uint8
				j uint8
			)
			mockPacketConn.EXPECT().WriteTo(gomock.Any(), gomock.Any()).Do(func(p []byte, addr net.Addr) (int, error) {
				if len(p) <= 12 {
					t.Errorf("got total packet length of %d, expected > 12 bytes", len(p))
					return len(p), nil
				}
				if p[0] != 0x1e || p[1] != 0x0f {
					t.Errorf("for magic bytes, got %v, expected %v", p[0:2], []byte{0x1e, 0x0f})
					return len(p), nil
				}
				if len(msgID) == 0 {
					msgID = p[2:10]
				}
				if !bytes.Equal(p[2:10], msgID) {
					t.Errorf("for message ID, got %v, expected %v", p[2:10], msgID)
					return len(p), nil
				}
				if p[10] != i {
					t.Errorf("for sequence index, got %v, expected %v", p[10], byte(i))
					return len(p), nil
				}
				i++
				if p[11] == 0 {
					t.Errorf("for chunk count, got %v, expected non-zero", p[11])
					return len(p), nil
				}
				if j == 0 {
					j = p[11]
				}
				if p[11] != j {
					t.Errorf("for chunk count, got %v, expected %v", p[11], byte(j))
					return len(p), nil
				}
				return len(p), nil
			}).AnyTimes()
			w.Write(tc.input)
		})
	}
}

func TestWriteZips(t *testing.T) {
	testCases := []struct {
		name          string
		input         []byte
		expectIsError bool
	}{
		{"small message",
			[]byte("{\"short\":{\"pi\":3.14,\"phi\":1.618}}\n"),
			false,
		},
		{"large message",
			[]byte(fmt.Sprintf("{\"long string\":\"%s\"}\n", strings.Repeat("pop ", 10000))),
			false,
		},
		{"too large message",
			[]byte(fmt.Sprintf("{\"long string\":\"%s\"}\n", strings.Repeat("pop ", 128*1400))),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockAddr := mock_net.NewMockAddr(mockCtrl)
			mockPacketConn := mock_net.NewMockPacketConn(mockCtrl)

			w, err := graylog.New(graylog.Config{
				ClientPacketConn: mockPacketConn,
				ServerAddr:       mockAddr,
			})
			if err != nil {
				t.Fatalf("error constructing New Client: %+v", err)
			}

			var buf bytes.Buffer
			mockPacketConn.EXPECT().WriteTo(gomock.Any(), gomock.Eq(mockAddr)).Do(func(p []byte, addr net.Addr) (int, error) {
				if !reflect.DeepEqual(addr, mockAddr) {
					t.Errorf("got server address of %+v, expected %+v", addr, mockAddr)
					return len(p), nil
				}
				if len(p) <= 12 {
					t.Errorf("for %s: got total packet length of %d, expected > 12 bytes", tc.name, len(p))
					return len(p), nil
				}
				if _, err := (&buf).Write(p[12:]); err != nil {
					t.Errorf("for %s: unexpected error writing to test buffer: %+v", tc.name, err)
				}
				return len(p), nil
			}).AnyTimes()
			_, err = w.Write(tc.input)
			if (err != nil) != tc.expectIsError {
				t.Errorf("for %s: Write() = _, %+v, expected err? %v", tc.name, err, tc.expectIsError)
				return
			}
			if err != nil {
				return
			}

			zbuf, err := gzip.NewReader(&buf)
			if err != nil && err != io.EOF {
				t.Errorf("for %s: failed to create gzip reader for results: %+v", tc.name, err)
				return
			}
			if zbuf == nil {
				t.Errorf("zip buffer is nil, expected non-nil")
			}
			var output bytes.Buffer
			_, err = (&output).ReadFrom(zbuf)
			if err != nil {
				t.Errorf("for %s: reading from Zip, got unexpected error: %+v", tc.name, err)
			}
			if !bytes.Equal(output.Bytes(), tc.input[:len(tc.input)-1]) {
				t.Errorf("once unzipped, got %d bytes, expected different %d bytes", output.Len(), len(tc.input)-1)
			}

		})
	}
}

func TestNew_Conn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPacketConn := mock_net.NewMockPacketConn(mockCtrl)
	mockPacketConn.EXPECT().Close().Times(1)

	udpConn, err := net.ListenPacket("udp", ":12201")
	if err != nil {
		t.Fatalf("dialing udp connection for test: %+v", err)
	}
	defer udpConn.Close()

	testCases := []struct {
		name          string
		conn          net.PacketConn
		wantNilWriter bool
		wantNilError  bool
	}{
		{"nil Conn",
			nil,
			true,
			false,
		},
		{"mock Conn",
			mockPacketConn,
			false,
			true,
		},
		{"udp Conn",
			udpConn,
			false,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.conn != nil {
				defer tc.conn.Close()
			}
			w, err := graylog.New(graylog.Config{
				ClientPacketConn: tc.conn,
			})
			if (w == nil) != tc.wantNilWriter {
				t.Errorf("New(Config{ClientPacketConn:%v}) = %+v, <error>, expected nil writer? %v", tc.conn, w, tc.wantNilWriter)
				return
			}
			if (err == nil) != tc.wantNilError {
				t.Errorf("New(Config{ClientPacketConn:%v}) = <Writer>, %+v, expected nil error? %v", tc.conn, err, tc.wantNilError)
			}
		})
	}
}

func TestNew_CompressionLevel(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPacketConn := mock_net.NewMockPacketConn(mockCtrl)

	testCases := []struct {
		level         int
		wantNilWriter bool
		wantNilError  bool
	}{
		{
			gzip.NoCompression, // same as 0
			false,
			true,
		},
		{
			gzip.BestSpeed,
			false,
			true,
		},
		{
			gzip.BestCompression,
			false,
			true,
		},
		{
			gzip.DefaultCompression,
			false,
			true,
		},
		{
			(gzip.BestCompression + gzip.BestSpeed) / 2,
			false,
			true,
		},
		{
			gzip.BestCompression + 1,
			true,
			false,
		},
		{
			-2,
			true,
			false,
		},
	}
	for _, tc := range testCases {
		w, err := graylog.New(graylog.Config{
			CompressionLevel: tc.level,
			ClientPacketConn: mockPacketConn,
		})
		if (w == nil) != tc.wantNilWriter {
			t.Errorf("New(Config{CompressionLevel:%d}) = %+v, <error>, expected nil writer? %v", tc.level, w, tc.wantNilWriter)
			continue
		}
		if (err == nil) != tc.wantNilError {
			t.Errorf("New(Config{CompressionLevel:%d}) = <Writer>, %+v, expected nil error? %v", tc.level, err, tc.wantNilError)
		}
	}
}

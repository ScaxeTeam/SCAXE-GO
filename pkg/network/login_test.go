package network

import (
	"net"
	"testing"
	"time"

	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

func TestLoginSequence(t *testing.T) {

	done := make(chan bool)
	logger.SetDebug(true)

	server := NewServer("0.0.0.0:19132", "TestServer", 10)

	clientConn, serverConn := createPipe()
	defer clientConn.Close()
	defer serverConn.Close()

	session := NewSession(serverConn, server)

	go func() {

		server.handleSession(session)
	}()

	loginPk := protocol.NewLoginPacket()
	loginPk.Username = "TestUser"
	loginPk.Protocol = 70
	loginPk.ClientUUID = "0123456789abcdef"
	loginPk.ServerAddress = "test"
	loginPk.SkinID = "Standard_Custom"

	stream := protocol.NewBinaryStream()
	loginPk.Encode(stream)

	packetData := stream.Bytes()
	batchStream := protocol.NewBinaryStream()

	batchStream.WriteInt(int32(len(packetData)))
	batchStream.WriteBytes(packetData)

	batch := protocol.NewBatchPacket()
	batch.Payload = batchStream.Bytes()

	if err := batch.Compress(batchStream.Bytes()); err != nil {
		t.Errorf("Compress failed: %v", err)
	}

	if _, err := batch.Decompress(); err != nil {
		t.Fatalf("Self-verification failed! Batch decompression error: %v. Payload len: %d", err, len(batch.Payload))
	} else {
		t.Logf("Self-verification passed. Compressed size: %d", len(batch.Payload))
	}

	finalStream := protocol.NewBinaryStream()
	finalStream.WriteByte(protocol.IDBatch)
	finalStream.WriteInt(int32(len(batch.Payload)))
	finalStream.WriteBytes(batch.Payload)

	go func() {
		t.Logf("Writing %d bytes to clientConn", finalStream.Len())
		n, err := clientConn.Write(finalStream.Bytes())
		if err != nil {
			t.Errorf("Write failed: %v", err)
		}
		t.Logf("Wrote %d bytes", n)
	}()

	readPacket := func() (protocol.DataPacket, error) {
		buf := make([]byte, 4096)
		n, err := clientConn.Read(buf)
		if err != nil {
			return nil, err
		}
		data := buf[:n]

		if data[0] == protocol.IDBatch {
			b := protocol.NewBatchPacket()
			stream := protocol.NewBinaryStreamFromBytes(data[1:])
			if err := b.Decode(stream); err != nil {
				return nil, err
			}

			decomp, err := b.Decompress()
			if err != nil {
				return nil, err
			}
			pkts, err := protocol.DecodePackets(decomp)
			if err != nil {
				return nil, err
			}
			return pkts[0], nil
		}
		return nil, nil
	}

	go func() {

		p1, err := readPacket()
		if err != nil {
			t.Error(err)
			return
		}
		if ps, ok := p1.(*protocol.PlayStatusPacket); !ok || ps.Status != protocol.PlayStatusLoginSuccess {
			t.Errorf("Expected PlayStatus(LoginSuccess 0), got %v", p1)
		}

		p2, err := readPacket()
		if err != nil {
			t.Error(err)
			return
		}
		if _, ok := p2.(*protocol.StartGamePacket); !ok {
			t.Errorf("Expected StartGamePacket, got %v", p2)
		}

		p3, err := readPacket()
		if err != nil {
			t.Error(err)
			return
		}
		if ps, ok := p3.(*protocol.PlayStatusPacket); !ok || ps.Status != protocol.PlayStatusPlayerSpawn {
			t.Errorf("Expected PlayStatus(PlayerSpawn 3), got %v", p3)
		}

		done <- true
	}()

	select {
	case <-done:

	case <-time.After(2 * time.Second):
		t.Error("Test timed out")
	}
}

type pipeAddr struct{}

func (pipeAddr) Network() string { return "pipe" }
func (pipeAddr) String() string  { return "pipe" }

func createPipe() (net.Conn, net.Conn) {
	return net.Pipe()
}

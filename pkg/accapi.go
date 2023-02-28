package accapi

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "net"
    "time"
)

// ACCUDPClient represents a UDP client that receives live data from ACC.
type ACCUDPClient struct {
    addr          *net.UDPAddr
    conn          *net.UDPConn
    lastPacket    []byte
    lastPacketSeq uint32
    connected     bool
}

// NewACCUDPClient creates a new instance of ACCUDPClient and sets the remote address.
func NewACCUDPClient(remoteAddr string) (*ACCUDPClient, error) {
    udpAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
    if err != nil {
        return nil, err
    }

    conn, err := net.ListenUDP("udp", nil)
    if err != nil {
        return nil, err
    }

    return &ACCUDPClient{
        addr:      udpAddr,
        conn:      conn,
        connected: true,
    }, nil
}

// ReadPacket reads a packet from the UDP connection.
func (c *ACCUDPClient) ReadPacket() ([]byte, uint32, error) {
    packet := make([]byte, 4096)
    c.conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
    n, _, err := c.conn.ReadFromUDP(packet)
    if err != nil {
        if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
            return nil, 0, fmt.Errorf("read timeout")
        }
        return nil, 0, err
    }
    packet = packet[:n]

    var seq uint32
    r := bytes.NewReader(packet[0:4])
    err = binary.Read(r, binary.LittleEndian, &seq)
    if err != nil {
        return nil, 0, err
    }

    if seq < c.lastPacketSeq {
        // If we've received an older packet, discard it.
        return nil, 0, nil
    }

    if seq > c.lastPacketSeq {
        // If there's a gap between packets, fill it with the last received packet.
        c.lastPacketSeq = seq
        c.lastPacket = packet
        return c.lastPacket, c.lastPacketSeq, nil
    }

    // If the sequence number is the same as the last received packet, use the last packet.
    return c.lastPacket, c.lastPacketSeq, nil
}

// Close closes the UDP connection.
func (c *ACCUDPClient) Close() error {
    c.connected = false
    return c.conn.Close()
}

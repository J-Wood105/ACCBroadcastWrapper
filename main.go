package main

import (
    "fmt"
    "time"

    "github.com/J-Wood105/ACCBroadcastWrapper/pkg/accapi"
)

func main() {
    client, err := accapi.NewACCUDPClient("127.0.0.1:9000")
    if err != nil {
        fmt.Printf("error creating ACC UDP client: %s\n", err)
        return
    }
    defer client.Close()

    for {
        packet, seq, err := client.ReadPacket()
        if err != nil {
            fmt.Printf("error reading packet: %s\n", err)
            continue
        }
        if packet == nil {
            // No packet received (e.g. read timeout).
            continue
        }

        // Process the packet (e.g. unpack the data).
        fmt.Printf("received packet %d")

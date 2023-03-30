package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"flag"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/charmbracelet/log"
	"github.com/checksum0/go-electrum/electrum"
)

func main() {
	ctx := context.Background()

	flag.Parse()
	peer := flag.Args()[0]

	client, err := electrum.NewClientTCP(ctx, peer)

	if err != nil {
		log.Fatal("Unable to setup tcp connection", "err", err)
	}

	// routine for check peer stills alive
	go func() {
		for {
			if err := client.Ping(ctx); err != nil {
				log.Fatal("Peer connection failed", "err", err)
			}
			time.Sleep(60 * time.Second)
		}
	}()

	// Create channel for listen new blocks emitted by peer
	ch, err := client.SubscribeHeaders(context.Background())
	if err != nil {
		log.Fatal("Unable to subscribe to block headers", "err", err)
	}

	// Once block header is received hash and height are logged
	for {
		blockHeader := <-ch

		d, err := hex.DecodeString(blockHeader.Hex)
		if err != nil {
			log.Error("Unable to decode block", "err", err)
		}

		msgBlock := &wire.MsgBlock{}
		err = msgBlock.Deserialize(bytes.NewBuffer(d))
		if err != nil {
			log.Error("Unable to deserialize block", "err", err)
		}

		blockHash := msgBlock.BlockHash()
		log.Infof("New block: Height %d, hash: %s", blockHeader.Height, blockHash.String())
	}
}

// Copyright (c) 2020 The Blocknet developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/blocknetdx/go-xrouter/blockcfg"
	"github.com/blocknetdx/go-xrouter/xrouter"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)

	client, err := xrouter.NewClient(blockcfg.MainnetParams)
	if err != nil {
		log.Println(err.Error())
		return
	}
	client.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()
	defer shutdown(client)

	if ready, err := client.WaitForXRouter(ctx); err != nil || !ready {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		log.Println("XRouter failed to connect and obtain service nodes", errStr)
		return
	}
	log.Printf("XRouter is ready")

	ctx2, cancel2 := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel2()
	queryCount := 1
	if err := client.WaitForService(ctx2, "xr::BLOCK", queryCount); err != nil {
		log.Printf("error: %v", err)
		return
	}
	if reply, err := client.GetBlockCount("xr::BLOCK", queryCount); err != nil {
		log.Printf("error: %v", err)
		return
	} else {
		log.Printf("result from %v: %v", hex.EncodeToString(reply.Pubkey), string(reply.Reply))
	}
}

func shutdown(client *xrouter.Client) {
	if err := client.Stop(); err != nil {
		fmt.Printf("error shutdown: %v", err)
	}
}

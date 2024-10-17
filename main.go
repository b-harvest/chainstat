package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/cometbft/cometbft/abci/types"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
)

func main() {
	// connect to a running node
	rpcAddress := "http://localhost:26657"
	client, err := rpchttp.New(rpcAddress, "/websocket")
	if err != nil {
		log.Fatalf("failed to create RPC client: %v", err)
	}

	args := os.Args[1:]
	// validate the number of arguments
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <start_block_height> <end_block_height>")
		os.Exit(1)
	}
	startBlockHeight, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Println("start_block_height must be a valid integer")
		os.Exit(1)
	}
	endBlockHeight, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Println("end_block_height must be a valid integer")
		os.Exit(1)
	}
	if startBlockHeight > endBlockHeight {
		fmt.Println("start_block_height must be less than end_block_height")
		os.Exit(1)
	}

	// Arrays to store the statistics for each 100-block range
	var avgTxsPer100Blocks []float64
	var avgBlockTimePer100Blocks []float64
	var avgTPSPer100Blocks []float64

	// go-routine related
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Loop through blocks in 100-block intervals
	for i := startBlockHeight; i <= endBlockHeight; i += 100 {
		var totalTxs int

		// Define the last block in the current 100-block range
		endRange := i + 99
		if endRange > endBlockHeight {
			endRange = endBlockHeight
		}

		// Process each block in the 100-block range
		for j := i; j <= endRange; j++ {
			wg.Add(1)
			go func(blockNum int64) {
				defer wg.Done()
				// Get the current block result
				currBlockResult, err := client.BlockResults(context.Background(), &blockNum)
				if err != nil {
					log.Fatalf("failed to get block results: %v", err)
				}
				// Increment the total number of transactions
				mu.Lock()
				totalTxs += len(currBlockResult.TxsResults)
				// check if there are any failed transactions
				for _, txResult := range currBlockResult.TxsResults {
					for _, ev := range txResult.Events {
						if ev.GetType() != "ethereum_tx" {
							continue
						}
						if isFailedEthTx(ev) {
							totalTxs--
							break
						}
					}
				}
				mu.Unlock()
			}(j)
		}
		wg.Wait()

		// Get block time difference between the start and end block
		startBlockInfo, err := client.Block(context.Background(), &i)
		if err != nil {
			log.Fatalf("failed to get block: %v", err)
		}
		endBlockInfo, err := client.Block(context.Background(), &endRange)
		if err != nil {
			log.Fatalf("failed to get block: %v", err)
		}
		totalTimeDiff := endBlockInfo.Block.Header.Time.Sub(startBlockInfo.Block.Header.Time).Seconds()

		// Calculate the averages for this 100-block range
		avgTxs := float64(totalTxs) / float64(endRange-i+1)
		avgBlockTime := totalTimeDiff / float64(endRange-i) // n-1 for block time differences
		tps := float64(totalTxs) / totalTimeDiff

		// Store the values in arrays for later use
		avgTxsPer100Blocks = append(avgTxsPer100Blocks, avgTxs)
		avgBlockTimePer100Blocks = append(avgBlockTimePer100Blocks, avgBlockTime)
		avgTPSPer100Blocks = append(avgTPSPer100Blocks, tps)

		// Output the statistics for this range
		fmt.Printf("Blocks %d-%d: AvgTxs: %.2f, AvgBlockTime: %.2fs, TPS=%.2f\n",
			i, endRange, avgTxs, avgBlockTime, tps)
	}

	// You can use avgTxsPer100Blocks, avgBlockTimePer100Blocks, avgTPSPer100Blocks for plotting graphs later
}

// isFailedEthTx checks if given event is related with failed Ethereum transaction
func isFailedEthTx(ev types.Event) bool {
	for _, attr := range ev.GetAttributes() {
		// https://github.com/b-harvest/ethermint/blob/2d35118b59d1ec73cf75b50e25008332f0f0867b/x/evm/types/events.go#L33
		// doesn't use AttributeKeyEthereumTxFailed because I want to keep minimum dependency (don't want to import x/evm/types)
		if attr.GetKey() == "ethereumTxFailed" {
			return true
		}
	}
	return false
}

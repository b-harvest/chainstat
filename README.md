# chainstat

**ChainStat** is an open-source tool designed to collect and analyze block statistics from a CometBFT-based blockchain. It efficiently processes blocks in specified ranges and provides detailed information such as average transactions per block, average block time, and transactions per second (TPS) across defined intervals. This tool leverages Go concurrency features like goroutines and mutexes to handle large datasets while ensuring performance and data integrity.

## Features
- Customizable Block Ranges: Allows you to specify a starting and ending block height to focus the analysis on the desired range of blocks.
- Statistics Collection: Computes detailed statistics for every 100-block range, including:
  - Average Transactions per Block (AvgTxs)
  - Average Block Time (AvgBlockTime)
  - Transactions per Second (TPS)

## Requirements
- **Go 1.23+**
- **CometBFT RPC Node:** The tool assumes you have access to a running CometBFT node’s RPC endpoint.

## Installation
First, clone the repository:
```bash
git clone https://github.com/b-harvest/chainstat.git
cd chainstat
```

Install the dependencies:
```bash
go mod tidy
```

## Usage
ChainStat requires the start and end block heights as arguments. It connects to a running CometBFT node and retrieves block data to compute statistics.

**Command**
```bash
go run main.go <start_block_height> <end_block_height>
```

**Example**
```bash
go run main.go 1000 5000
```

In this example, ChainStat will compute block statistics for blocks 1000 to 5000, processing them in 100-block intervals. For each interval, you’ll receive output similar to:
```bash
Blocks 1000-1099: AvgTxs: 12.34, AvgBlockTime: 1.23s, TPS=9.87
Blocks 1100-1199: AvgTxs: 11.56, AvgBlockTime: 1.45s, TPS=8.34
...
```

**Output**
- **AvgTxs:** The average number of transactions per block within the 100-block range.
- **AvgBlockTime:** The average time (in seconds) between blocks within the range.
- **TPS:** The average number of transactions per second processed in the 100-block range.

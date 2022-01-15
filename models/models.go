package main

import (
	"time"
)

// ====================================
// Database Tables
// ====================================

// A mining algorithm such as scrypt.
// Making this distinct will allow mapping between various pools and mining software.
// NOTE: These should only be algorithms supported by a pool provider. There is no reason to add
// algorithms that are not supported by a pool.
type Algorithm struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string
}

// A crypto coin.
type Coin struct {
	ID          uint64 `gorm:"primaryKey"`
	CoinGeckoID string
	Name        string
	Symbol      string
	Added       time.Time // The date/time added. Can be used to track new coins.
}

// A price for a coin in USD. This is not OHLC over a range.
type CoinPrice struct {
	ID      uint64 `gorm:"primaryKey"`
	CoinID  uint64
	Instant time.Time // The instant of the price
	Price   float64
}

// Mining hardware.
type Miner struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string // A unique name for the hardware
}

// Software executed on the miner to mine.
type MinerSoftware struct {
	ID      uint64 `gorm:"primaryKey"`
	Name    string // The general name of the miner, e.g. lolMiner, cpuminer-opt, etc.
	Website string // Preferably the Github repository website
	// This is important. The miner statistics program utilizes this to determine
	// which miner software is in use on-the-fly. It must match the executable.
	// For instance, if the executable is minerExecv1.2.3.exe, minerExec should
	// probably be the prefix as the version numbers may fluctuate over time.
	ExecutablePrefix string
	// The parameter for an algorithm which can be unique to each miner software.
	// Typically something like --algo. This does NOT include the value.
	AlgoParam string
	// Any other raw parameters and their values. This is just tacked onto the executable at runtime.
	OtherParams string
}

// Miner software can call the same algorithm differently. This is used as a map between the algorithms
// stored in the algorithms table and the miner's potentially unique naming. Also, this indicates which
// algorithms the software supports as that is often unique per software too.
type MinerSoftwareAlgos struct {
	ID              uint64 `gorm:"primaryKey"`
	MinerSoftwareID uint64
	AlgorithmID     uint64
	// The name utilized by the miner software. This will be applied to the algorithm parameter at runtime.
	// NOTE: If this is blank or null, the name from the algorithm table is used.
	Name string
	// If the algo requires additional parameters, that is handled here.
	ExtraParams string
}

// Statistics for a miner, specific software, and an algorithm. This can be utilized to make estimates
// on profitability in combination with pool statistics from a mining pool provider.
type MinerStats struct {
	ID              uint64 `gorm:"primaryKey"`
	MinerID         uint64
	MinerSoftwareID uint64
	AlgorithmID     uint64
	Instant         time.Time // The date/time the statistics were gathered.
	WorkPerSecond   float64   // The hash rate. Could be Mh/s, Kh/s, h/s etc. See MhFactor.
	// Used for mining profit calculations/estimates
	MhFactor float64 // 1 = Mh/s, 0.001 = Kh/s, 1000 = Gh/s
}

// A mining pool.
type Pool struct {
	ID          uint64 `gorm:"primaryKey"`
	ProviderID  uint64
	AlgorithmID uint64 // The name will not necessarily match the algorithm precisely
	Name        string // In some cases, the name of the pool could be different than the algo.
	URL         string // This is generated to the full address for easier automation
	Port        uint32
	// Used for mining profit calculations/estimates
	MhFactor float64 // 1 = Mh/s, 0.001 = Kh/s, 1000 = Gh/s
}

// Statistics at a certain point in time for a pool.
// Used to optimize mining operations by examining profit estimates/actuals.
type PoolStats struct {
	ID                  uint64 `gorm:"primaryKey"`
	PoolID              uint64
	Instant             time.Time // The date/time of the statistics
	CurrentHashrate     uint64    // The current shared hashrate for the pool
	Workers             uint32    // The current number of workers sharing the pool
	ProfitEstimate      float64   // An forward look at potential profit/day
	ProfitActual24Hours float64   // The actual profit/day for those sharing the pool
	CoinPriceID         uint64    // The line to the relevant coin price, if any. Bitcoin is usually used.
}

// A pool provider such as ZergPool.
type Provider struct {
	ID      uint64 `gorm:"primaryKey"`
	Name    string
	Website string
	Fee     float32
}

func main() {
	// Does nothing
}

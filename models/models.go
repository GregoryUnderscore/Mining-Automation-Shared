package models

import (
	"time"
)

// ====================================
// Database Tables
// ====================================

// This is used to determine if the schema should be updated.
const SchemaVersion = 1.05

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
	ID                  uint64 `gorm:"primaryKey"`
	Name                string // A unique name for the hardware
	MinerSoftwareAlgoID uint64 // The active software/algo mining on this miner (or last known active)
}

// Mining hardware / mining software bridge
// This is a bridge between mining hardware and mining software and can be used to store details specific
// to the miner/software such as file paths.
type MinerMinerSoftware struct {
	ID              uint64 `gorm:"primaryKey"`
	MinerID         uint64
	MinerSoftwareID uint64
	FilePath        string // The file path to the software on the miner.
}

// Software executed on the miner to mine.
type MinerSoftware struct {
	ID      uint64 `gorm:"primaryKey"`
	Name    string // The general name of the miner, e.g. lolMiner, cpuminer-opt, etc.
	Website string // Preferably the Github repository website
	// The parameter for an algorithm which can be unique to each miner software.
	// Typically something like --algo. This does NOT include the value.
	AlgoParam string
	// An optional parameter that will require the software to connect to a pool. Some software
	// does not have a benchmark mode and must actually connect.
	PoolParam string
	// An optional parameter that will pass a wallet to the pool. Some mining software requires
	// connecting to a pool, and some pools require a wallet to connect.
	WalletParam string
	// A parameter often used to identify miners on the pool. Can sometimes be used to set options too.
	PasswordParam string
	// Some software allows logging to a file. This is optional. If this is not set, the screen
	// output is saved to a file if possible.
	FileParam string
	// Any other raw parameters and their values. This is just tacked onto the executable at runtime.
	OtherParams string
	// How many lines to skip on the output. Some software outputs low hashrates initially, and those
	// can be skipped. Setting this to 1 skips 1 line and so forth.
	SkipLines uint8
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
	// If this is 1, do not use in automated optimization. Some software tends to be buggy
	// or produce invalid shares more than others. In those cases, it may makes sense to exclude.
	DoNotUse bool
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

// This is used to track the schema version of the database and possibly other things to identify
// if an update is necessary. Avoids unnecessary updates.
type Version struct {
	ID      uint8 `gorm:"primaryKey"`
	Name    string
	Version float32
}

func main() {
	// Does nothing
}

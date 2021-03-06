package database

// General database functions shared across various repositories.

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	. "github.com/GregoryUnderscore/Mining-Automation-Shared/models"
)

// Connect to a PostgreSQL database according to the passed parameters.
// @param host - The database server
// @param port - The port for the database server
// @param database - The database to use (must be manually created as of 1/14/22).
// @param user - The user to use for login
// @param password - The user's password for login
// @param timezone - The time zone where the program is executed, e.g. America/Chicago.
// @returns A pointer to the GORM database connection
func Connect(host string, port string, database string, user string, password string,
	timezone string) *gorm.DB {
	// Grab the configuration details for the database connection. These are stored in ZergPoolData.hcl.
	log.Println("Using the following configuration:")
	log.Println("Database Server: " + host + ":" + port)
	log.Println("Database: " + database)
	log.Println("User: " + user + "\n")

	dsn := "host=" + host + " "
	dsn += "port=" + port + " "
	dsn += "dbname=" + database + " "
	dsn += "user=" + user + " "
	dsn += "password=" + password + " "
	dsn += "TimeZone=" + timezone + " "
	dsn += "sslmode=disable"
	log.Println("Connecting to " + host + "...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to the database server.\n", err)
	}

	log.Println("Connected to " + host + ".")
	return db
}

// Create the database version record and store the version.
// @param db - The active database connection
// @returns Nothing
func createSchemaVersion(db *gorm.DB) {
	log.Println("Storing schema version as " + fmt.Sprint(SchemaVersion) + "...")
	var schemaVersion Version
	schemaVersion.Name = "database"
	schemaVersion.Version = SchemaVersion
	result := db.Create(&schemaVersion)
	if result.Error != nil {
		log.Fatalf("Issue creating schema version.\n", result.Error)
	}
}

// Create or update the database schema according to the shared models.
// @param db - The active database connection
func updateSchema(db *gorm.DB) {
	log.Println("Creating/updating schema...")
	// Create the schema if it does not exist. This also will perform alterations.
	// ==> Schema required for ZergPool statistics.
	db.AutoMigrate(&Version{}, &Provider{}, &Algorithm{}, &Pool{}, &PoolStats{},
		&Coin{}, &CoinPrice{})
	// ==> Schema required for miner statistics.
	db.AutoMigrate(&Miner{}, &MinerStats{}, &MinerSoftware{}, &MinerSoftwareAlgos{},
		&MinerMinerSoftware{})
}

// Verify the current schema contains all the appropriate tables, and if not, create/update them
// according to the current models.
// @param db - The active database connection
// @returns Nothing
func VerifyAndUpdateSchema(db *gorm.DB) {
	log.Println("Verifying schema...")
	if db.Migrator().HasTable(&Version{}) {
		var schemaVersion Version
		db.Where("name = ?", "database").Find(&schemaVersion)
		// If the database schema version is old, update it.
		if (Version{}) == schemaVersion || schemaVersion.Version < SchemaVersion {
			log.Println("Found older schema v" + fmt.Sprint(schemaVersion.Version) + "...")
			updateSchema(db) // Create/update the database schema
			// Ensure the schema version is up-to-date.
			if (Version{}) == schemaVersion {
				createSchemaVersion(db)
			} else {
				log.Println("Updating schema version to " + fmt.Sprint(SchemaVersion) + "...")
				schemaVersion.Version = SchemaVersion
				result := db.Save(&schemaVersion)
				if result.Error != nil {
					log.Fatalf("Issue updating schema version.\n", result.Error)
				}
			}
		}
	} else { // Table does not exist. An update is definitely needed.
		updateSchema(db) // Create/update the database schema
		createSchemaVersion(db)
	}
	log.Println("Schema verified.")
}

// Verify the miner exists in the database. If not, create it.
// @param tx - The active database session
// @param minerName - The name of the mining hardware
// @returns The ID associated with the miner.
func VerifyMiner(tx *gorm.DB, minerName string) uint64 {
	var miner Miner
	result := tx.Where("name = ?", minerName).Limit(1).Find(&miner)
	if result.RowsAffected == 0 {
		log.Println("Creating miner...")
		miner.Name = minerName
		result = tx.Create(&miner)
		if result.Error != nil {
			log.Fatalf("Issue creating miner.\n", result.Error)
		}
	} else if result.Error != nil {
		log.Fatalf("Unknown issue storing miner.\n", result.Error)
	} else {
		log.Println("Found existing miner.")
	}
	return miner.ID
}

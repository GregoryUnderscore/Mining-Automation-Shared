package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	. "Mining-Automation-Shared/models"
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

// Verify the current schema contains all the appropriate tables, and if not, create/update them
// according to the current models.
// @param db - The active database connection
// @returns Nothing
func VerifyAndUpdateSchema(db *gorm.DB) {
	log.Println("Verifying/updating schema")
	// Create the schema if it does not exist. This also will perform alterations.
	db.AutoMigrate(&Provider{}, &Algorithm{}, &Pool{}, &PoolStats{},
		&Coin{}, &CoinPrice{})
	log.Println("Schema verified.")
}

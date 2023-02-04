package HiveCmd

import (
    "crypto/rand"
    "database/sql"
    "encoding/base64"
    _ "github.com/gin-gonic/gin"
    _ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
    "log"
	"os"
)

// Function to initialize database
func initDB(db *sql.DB) {
    var err error

    log.Printf("Initializing Tables...")
    defer log.Printf("Tables successfully initialized")
    
    // Create "drones" table
    log.Printf("Creating \"drones\" table if not already present...")
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS drones (
                        id SERIAL PRIMARY KEY,
                        address INET,
                        port INTEGER,
                        name TEXT,
                        UNIQUE (address, port),
                        UNIQUE (name)
                      )`)
    if err != nil {
        log.Fatalln(err)
    }
    log.Printf("Success")

    // Create "permissions" table
    log.Printf("Creating \"permissions\" table if not already present...")
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS permissions (
                        id SERIAL PRIMARY KEY,
                        name TEXT
                      )`)
    if err != nil {
        log.Fatalln(err)
    }
    log.Printf("Success")

    // Create "groups" table
    log.Printf("Creating \"groups\" table if not already present...")
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS groups (
                        id SERIAL PRIMARY KEY,
                        name TEXT,
                        permissions_id SERIAL
                        REFERENCES permissions (id),
                        UNIQUE (name)
                      )`)
    if err != nil {
        log.Fatalln(err)
    }
    log.Printf("Success")

    // Create "swarms" table
    log.Printf("Creating \"swarms\" table if not already present...")
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS swarms (
                        id SERIAL PRIMARY KEY,
                        name TEXT,
                        drones_id SERIAL
                        REFERENCES drones (id),
                        UNIQUE (name)
                      )`)
    if err != nil {
        log.Fatalln(err)
    }
    log.Printf("Success")

    // Create "tokens" table
    log.Printf("Creating \"tokens\" table if not alrady present...")
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
                        id SERIAL PRIMARY KEY,
                        key TEXT,
                        created TIMESTAMP,
                        UNIQUE (key)
                      )`)
    if err != nil {
        log.Fatalln(err)
    }
    log.Printf("Success")

    // Create "users" table
    log.Printf("Creating \"users\" table if not already present...")
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
                        id SERIAL PRIMARY KEY,
                        name TEXT,
                        groups_id SERIAL
                        REFERENCES groups (id),
                        pass TEXT,
                        created TIMESTAMP,
                        UNIQUE (name)
                      )`)
    if err != nil {
        log.Fatalln(err)
    }
    log.Printf("Success")
}

// Function to generate a random alphanumeric string of set length
func RandStringBytes(n int) string {
    randomBytes := make([]byte, 64)
    _, err := rand.Read(randomBytes)
    if err != nil {
        log.Println(err)
    }
    return base64.StdEncoding.EncodeToString(randomBytes)[:n]
}

var (
    HiveCmd = &cobra.Command{
	    Use:   "Backstage-Hive",
	    Short: "Short Desc",
	    Long: `Long
               Desc`,
        PersistentPreRun: func(cmd *cobra.Command, args []string) {
        },
        Run: func(cmd *cobra.Command, args []string) {
        },
    }
)

func Execute() {
	err := HiveCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	HiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

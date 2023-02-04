package main
 
import (
    "crypto/rand"
    "database/sql"
    "encoding/base64"
    "fmt"
    _ "github.com/gin-gonic/gin"
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/spf13/viper"
    "log"
    _ "net/http"
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

// Function to generate an enrollment token
func genEnrollmentToken(db *sql.DB, host string, port int) {
    var key string = RandStringBytes(50)
    var enrollmentToken string = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"addr":"%s:%v","key":"%s"}`, host, port, key)))
    var execStr string = fmt.Sprintf(`INSERT INTO tokens (id, key, created)
                                      VALUES (DEFAULT, '%s', CURRENT_TIMESTAMP)`, key)
    db.Exec(execStr)
    log.Printf("Generated Token: \"%s\"", enrollmentToken)
    log.Printf("Generated Key: \"%s\"", key)
}

// Function to enroll a Drone into the Hive inventory
func enrollDrone(db *sql.DB) {
    var droneAddress string = "10.13.0.25"
    var dronePort int = 3802
    var droneName string = "drone-1"
    var execStr string = fmt.Sprintf(`INSERT INTO drones (id, address, port, name)
                                      VALUES (DEFAULT, '%s', %v, '%s')`,
                                      droneAddress, dronePort, droneName)
    db.Exec(execStr)
    log.Printf("drone \"%s\" Enrolled", droneName)
}

func main() {
    log.SetFlags(log.Lshortfile)
    log.SetPrefix("Backstage-Hive: ")

    vi := viper.New()
    vi.SetConfigFile("config.yaml")
    vi.ReadInConfig()

    vi.SetDefault("host", "localhost")
    vi.SetDefault("port", 6789)
    vi.SetDefault("dbHost", "localhost")
    vi.SetDefault("dbPort", 5432)
    vi.SetDefault("dbUser", "backstage")
    vi.SetDefault("dbPass", "backstage")
    vi.SetDefault("dbName", "backstage")

    psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
                             vi.GetString("dbHost"), vi.GetInt("dbPort"), vi.GetString("dbUser"),
                             vi.GetString("dbPass"), vi.GetString("dbName"))

    // Connect to database
    log.Printf("Connecting to database...")
    db, err := sql.Open("pgx", psqlconn)
    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }
    defer log.Printf("Database connection closed")
    defer db.Close()

    // Verify database connection
    err = db.Ping()
    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }
    log.Printf("Connection established")

    // Initialize database
    initDB(db)

    //genEnrollmentToken(db, vi.GetString("host"), vi.GetInt("port"))
}

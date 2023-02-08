package main

import (
    "crypto/rand"
    "database/sql"
    "encoding/base64"
    "encoding/json"
    "fmt"

    "github.com/Euvaz/Backstage-Hive/logger"
    "github.com/Euvaz/Backstage-Hive/models"
    "github.com/Euvaz/Backstage-Hive/pkg"
    "github.com/gin-gonic/gin"
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

func main() {
	viper := viper.New()
	viper.SetConfigFile("config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal(err.Error())
	}

	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 6789)
	viper.SetDefault("dbHost", "localhost")
	viper.SetDefault("dbPort", 5432)
	viper.SetDefault("dbUser", "backstage")
	viper.SetDefault("dbPass", "backstage")
	viper.SetDefault("dbName", "backstage")

	db := getDB(viper.GetString("dbHost"), viper.GetInt("dbPort"), viper.GetString("dbUser"), viper.GetString("dbPass"), viper.GetString("dbName"))

    // Add root command
	cmd := &cobra.Command {
		Use:   "Backstage-Hive",
		Short: "Short Desc",
		Long:  `Long Desc`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Starting server")
			initDB(db)
			defer closeDB(db)

            router := gin.Default()
            registerRoutes(router, db)
            addr, _ := pkg.ParseHost(viper.GetString("host"))
            port := viper.GetInt("port")
            // Attempt running on IPv4/Hostname
            err := router.Run(fmt.Sprintf("%s:%d", addr, port))
            if err != nil {
                // Attempt running on IPv6
                err := router.Run(fmt.Sprintf("[%s]:%d", addr, port))
                if err != nil {
                    logger.Fatal(err.Error())
                }
            }
		},
	}

    // Add command
    createCmd := &cobra.Command {
        Use:   "create",
        Short: "Short Desc",
        Long:  `Long Desc`,
    }

    // Add subcommand
    createTokenCmd := &cobra.Command {
        Use:   "token",
        Short: "Short Desc",
        Long:  `Long Desc`,
        Aliases: []string{"tok", "tokens"},
        Run: func(cmd *cobra.Command, args []string) {
            logger.Debug("Creating token...")
            key  := RandStringBytes(50)
            address, hostname := pkg.ParseHost(viper.GetString("host"))
            tokenBytes, err := json.Marshal(models.Token{Addr: address, Port: viper.GetInt("port"), Host: hostname, Key: key})
            enrollmentToken := base64.StdEncoding.EncodeToString(tokenBytes)

            _, err = db.Exec(`INSERT INTO tokens (id, key, created)
                               VALUES (DEFAULT, $1, CURRENT_TIMESTAMP)`, key)
            if err != nil {
                logger.Fatal(err.Error())
            }
            fmt.Println("Generated Token:", enrollmentToken)
            logger.Debug("Created token")
        },
    }

    // Add command
    getCmd := &cobra.Command {
        Use:   "get",
        Short: "Short Desc",
        Long:  `Long Desc`,
    }

    // Add subcommand
    getDroneCmd := &cobra.Command {
        Use:   "drone",
        Short: "Short Desc",
        Long:  `Long Desc`,
        Aliases: []string{"dr", "drones"},
        Run: func(cmd *cobra.Command, args []string) {
            rows, err := db.Query(`SELECT address, port, name, hostname FROM drones`)
            if err != nil {
                logger.Fatal(err.Error())
            }
            defer rows.Close()

            f := "%-15s %-6s %-13s, %s\n"
            fmt.Printf(f, "ADDRESS", "PORT", "NAME", "HOSTNAME")

            var address string
            var port string
            var name string
            var hostname string

            for rows.Next() {
                if err := rows.Scan(&address, &port, &name, &hostname); err != nil {
                    logger.Fatal(err.Error())
                }
                fmt.Printf(f, address, port, name, hostname)
            }
            if err = rows.Err(); err != nil {
                logger.Fatal(err.Error())
            }
        },
    }

    // Add subcommand
    getTokenCmd := &cobra.Command {
        Use:   "token",
        Short: "Short Desc",
        Long:  `Long Desc`,
        Aliases: []string{"tok", "tokens"},
        Run: func(cmd *cobra.Command, args []string) {
            rows, err := db.Query(`SELECT key, created FROM tokens`)
            if err != nil {
                logger.Fatal(err.Error())
            }
            defer rows.Close()

            f := "%-50s %s\n"
            fmt.Printf(f, "KEY", "CREATED")

            var key string
            var created string

            for rows.Next() {
                if err := rows.Scan(&key, &created); err != nil {
                    logger.Fatal(err.Error())
                }
                fmt.Printf(f, key, created)
            }
            if err = rows.Err(); err != nil {
                logger.Fatal(err.Error())
            }
        },
    }

    // Add commands
    cmd.AddCommand(createCmd)
    cmd.AddCommand(getCmd)

    // Add subcommands
    createCmd.AddCommand(createTokenCmd)
    getCmd.AddCommand(getDroneCmd)
    getCmd.AddCommand(getTokenCmd)


	err = cmd.Execute()
	if err != nil {
		logger.Fatal(err.Error())
	}
}

// Function to generate a random alphanumeric string of set length
func RandStringBytes(n int) string {
    randomBytes := make([]byte, 64)
    _, err := rand.Read(randomBytes)
    if err != nil {
        logger.Fatal(err.Error())
    }
    return base64.StdEncoding.EncodeToString(randomBytes)[:n]
}

func initDB(db *sql.DB) {
	var err error

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS drones (
                        id SERIAL PRIMARY KEY,
                        address INET,
                        port INTEGER,
                        name TEXT,
                        hostname TEXT,
                        UNIQUE (address, port),
                        UNIQUE (name),
                        UNIQUE (hostname)
                      )`)
	if err != nil {
		logger.Fatal(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS permissions (
                        id SERIAL PRIMARY KEY,
                        name TEXT
                      )`)
	if err != nil {
		logger.Fatal(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS groups (
                        id SERIAL PRIMARY KEY,
                        name TEXT,
                        permissions_id SERIAL
                        REFERENCES permissions (id),
                        UNIQUE (name)
                      )`)
	if err != nil {
		logger.Fatal(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS swarms (
                        id SERIAL PRIMARY KEY,
                        name TEXT,
                        drones_id SERIAL
                        REFERENCES drones (id),
                        UNIQUE (name)
                      )`)
	if err != nil {
		logger.Fatal(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
                        id SERIAL PRIMARY KEY,
                        key TEXT,
                        created TIMESTAMP,
                        UNIQUE (key)
                      )`)
	if err != nil {
		logger.Fatal(err.Error())
	}

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
		logger.Fatal(err.Error())
	}

	logger.Info("Tables successfully initialized")
}

func getDB(host string, port int, user string, pass string, name string) *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	logger.Debug("Connecting to database...")
	database, err := sql.Open("pgx", psqlconn)
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = database.Ping()
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("Connection established")

	return database
}

func closeDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Database connection closed")
}

// Function to verify authenticity of enrollment key
func enrollmentKeyIsValid(db *sql.DB, key string) bool {
    var count string
    rows := db.QueryRow(`SELECT COUNT (*) FROM tokens WHERE key = $1`, key)

    err := rows.Scan(&count)
    if err != nil {
        logger.Fatal(err.Error())
    }

    switch count {
    case "1":
        return true
    default:
        return false
    }
}

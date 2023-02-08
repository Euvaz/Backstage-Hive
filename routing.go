package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "io/ioutil"
    _ "net/http"

    "github.com/Euvaz/Backstage-Hive/logger"
    "github.com/Euvaz/Backstage-Hive/models"
    "github.com/gin-gonic/gin"
)

func registerRoutes (router *gin.Engine, db *sql.DB) {
    router.GET("/drones/:name", func(ctx *gin.Context) {
        logger.Info("Handling GET /drones")
    })

    router.POST("/drones/:name", func(ctx *gin.Context) {
        logger.Info("Handling POST /drones")
        name := ctx.Param("name")

        jsonData, err := ioutil.ReadAll(ctx.Request.Body)
        if err != nil {
            logger.Error(err.Error())
            ctx.AbortWithStatus(400)
            return
        }

        var token models.Token
        err = json.Unmarshal(jsonData, &token)
        if err != nil {
            logger.Error(err.Error())
            ctx.AbortWithStatus(400)
            return
        }

        if enrollmentKeyIsValid(db, token.Key) {
            _, err := db.Exec(`INSERT INTO drones (id, address, port, name, hostname)
                               VALUES (DEFAULT, $1, $2, $3, $4)`, token.Addr, token.Port, name, token.Host)
            if err != nil {
                logger.Error(err.Error())
                ctx.AbortWithStatus(500)
                return
            }
            logger.Info(fmt.Sprintf(`Drone "%s" Enrolled`, name))
        } else {
            logger.Warn(err.Error())
            ctx.AbortWithStatus(401)
            return
        }
    })
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/meta-node-blockchain/meta-node/pkg/bls"
	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/toolSign/config"
	"flag"
	"github.com/meta-node-blockchain/toolSign/api"
	"github.com/ethereum/go-ethereum/crypto"
	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/ethereum/go-ethereum/common"
	"encoding/hex"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"os"
)
const (
	defaultConfigPath = "config.yaml"
	defaultLogLevel   = logger.FLAG_INFO
)
var (
	CONFIG_FILE_PATH string
	LOG_LEVEL        int
)

// Struct to hold the response for domain verification
type SignRequest struct {
	PrivateKey   string 	`json:"privatekey"`
	Message		 string   	`json:"message"`
}


func main() {
	flag.StringVar(&CONFIG_FILE_PATH, "config", defaultConfigPath, "Config path")
	flag.IntVar(&LOG_LEVEL, "log-level", defaultLogLevel, "Log level")
	flag.IntVar(&LOG_LEVEL, "ll", defaultLogLevel, "Log level (shorthand)")

	flag.Parse()
	var loggerConfig = &logger.LoggerConfig{
		Flag:    LOG_LEVEL,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)
	config, err := config.LoadConfig(CONFIG_FILE_PATH)
	if err != nil {
		log.Fatal("can not load config", err)
	}
	engine := gin.Default()
	r := engine.Group("/get-sign")
	{
		r.POST("", GetSign)
	}

	fmt.Printf("Starting server on %s...",config.API_PORT)
	engine.Run(config.API_PORT)
}
func GetSign(c *gin.Context) {
	var queryParam SignRequest
	err := c.ShouldBindJSON(&queryParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Unable to bind query %v", err),
		})
		return
	}

	response := getSign(queryParam.PrivateKey,queryParam.Message)
	message := gin.H{
		"message": "successful request",
		"sign":    hex.EncodeToString(response.Bytes()),
	}
	api.ResponseWithStatusAndData(http.StatusOK, message, c)
}

func getSign(privateKey string, message string) cm.Sign {
	keyPair := bls.NewKeyPair(common.FromHex(privateKey))
	prikey := keyPair.PrivateKey()
	sign := bls.Sign(prikey, crypto.Keccak256([]byte(message)))  
	return sign
}
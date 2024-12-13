package app

import (
	"fmt"
	"log"
	"os"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/config"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/database"
	"github.com/meta-node-blockchain/meta-node-baccarat/route"
	"github.com/meta-node-blockchain/meta-node/cmd/client"
	c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/services"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/controller"
)
type App struct {
	Config *config.AppConfig
	ApiApp *gin.Engine
	ChainClient *client.Client
	StopChan    chan bool
}

func NewApp(
	configPath string,
	loglevel int,
)(*App, error) {
	var loggerConfig = &logger.LoggerConfig{
		Flag:    loglevel,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("can not load config", err)
	}
	app := &App{}
	engine := gin.Default()
	app.ChainClient, err = client.NewStorageClient(
		&c_config.ClientConfig{
			Version_:                config.MetaNodeVersion,
			PrivateKey_:             config.PrivateKey_,
			ParentAddress:           config.NodeAddress,
			ParentConnectionAddress: config.NodeConnectionAddress,
			DnsLink_:                config.DnsLink(),
		},
		[]common.Address{
			common.HexToAddress(config.BaccaratAddress),
		},
	)

	if err != nil {
		logger.Error(fmt.Sprintf("error when create chain client %v", err))
		return nil, err
	}
	leveldb, err :=database.Open(config.PathLevelDB)
	// create card abi
	reader, err := os.Open(config.BaccaratABIPath)
	if err != nil {
		logger.Error("Error occured while read baccarat abi")
		return nil, err
	}
	defer reader.Close()

	abi, err := abi.JSON(reader)
	if err != nil {
		logger.Error("Error occured while parse baccarat smart contract abi")
		return nil, err
	}
	// Initialize services

	servs := services.NewSendTransactionService(
		app.ChainClient, 
		&abi, 
		common.HexToAddress(config.BaccaratAddress),
	)
	controller := controller.NewController(leveldb,servs)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	//
	engine.Use(cors.New(corsConfig))
	route.InitialRoutes(
		engine,
		controller,
	)
	app.Config = config
	app.ApiApp = engine
	return app, nil
}

func (app *App) Run () {
	app.StopChan = make(chan bool)
	go func() {
		app.ApiApp.Run(app.Config.API_PORT)
	}()
	for {
		select {
			case <-app.StopChan:
				return
		}
	}
}
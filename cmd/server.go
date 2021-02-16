package main

import (
	"context"
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"project-evredika/internal/eventer/channels"
	receiver "project-evredika/internal/eventer/event_receiver"
	"project-evredika/internal/eventer/event_sender"
	"project-evredika/internal/eventer/models"
	"project-evredika/internal/indexer/common"
	"project-evredika/internal/indexer/lister"
	"project-evredika/internal/storage/data_saver"
	"project-evredika/internal/storage/transaction"
	user_handler "project-evredika/internal/user/net/http"
	user_storage "project-evredika/internal/user/storage"
	user_usecase "project-evredika/internal/user/usecase"
	service_cfg "project-evredika/pkg/cfg"
	"project-evredika/pkg/echo/query_params"
	"project-evredika/pkg/eventer/event_receiver"
	"project-evredika/pkg/storage/configurer"
	"project-evredika/pkg/storage/data_saver/middlewares"
	user_indexer "project-evredika/pkg/storage/index_maker/user"
)

const (
	swapValue    = 1
	defaultValue = 0
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	var cfg service_cfg.Config
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Error("failed to process service env", "error", err)
		os.Exit(1)
	}
	// creating storage
	storageCfgrer := configurer.NewStorageConfigurer(
		map[string]configurer.CfgHandler{
			configurer.S3: configurer.CreateS3Storage,
			configurer.OS: configurer.CreateOSStorage,
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := storageCfgrer.Configure(ctx, cfg.StorageType, cfg.Bucket)
	if err != nil {
		logger.Error("failed to create storage", "error", err)
		os.Exit(1)
	}

	// creating data - indexers
	indexer := user_indexer.NewIndexMaker(cfg.KeyFormat, logger)

	// creating indexers
	userIndexer := common.NewUserIndexer(lister.NewLister())

	// creating event handlers
	putHandler := event_receiver.NewUserPutIndexHandler(userIndexer, logger)
	deleteHandler := event_receiver.NewUserDeleteIndexHandler(userIndexer)

	// creating event-receivers
	receiverDelete := receiver.NewReceiver(
		channels.NewMetaChannel(cfg.ChannelSize),
		[]event_receiver.Handler{deleteHandler},
	)
	receiverPut := receiver.NewReceiver(
		channels.NewMetaChannel(cfg.ChannelSize),
		[]event_receiver.Handler{putHandler},
	)
	// creating event-senders
	sender := event_sender.NewSender(channels.NewEventChannel(cfg.ChannelSize))
	sender.Subscribe(models.EventPut, receiverPut)
	sender.Subscribe(models.EventDelete, receiverDelete)
	// adopting storage with middlewares
	storage = middlewares.NewStorageEventer(sender, indexer, storage)
	storage = middlewares.NewStorageLogger(logger, storage)
	// creating storage with transaction
	storageWithTransaction := transaction.NewTransaction(storage, swapValue, defaultValue)
	// starting eventers
	wg := &sync.WaitGroup{}

	wg.Add(3)
	defer wg.Wait()

	go receiverPut.Start(ctx, wg)
	go receiverDelete.Start(ctx, wg)
	go sender.Start(ctx, wg)

	// prepare indexes
	if _, err = storage.ListData(ctx, &data_saver.Metadata{Bucket: cfg.Bucket}); err != nil {
		logger.Error("failed to prepare indexes", "error", err)
		os.Exit(1)
	}
	// creating user storage
	userStorage := user_storage.NewUserStorage(
		storageWithTransaction,
		userIndexer,
		indexer,
		logger,
		cfg.Bucket,
	)
	userStorage = user_storage.NewUserStorageLogging(
		logger,
		userStorage,
	)
	// creating user usecase
	userUsecase := user_usecase.NewUserUsecase(userStorage)
	// creating user http/handler
	userHandler := user_handler.NewUserHttpDelivery(
		logger,
		userUsecase,
		query_params.NewQueryGetter(defaultValue),
	)

	// initializing server
	server := echo.New()
	userHandler.Initiate(server)

	// starting server
	logger.Debug("starting server at port", cfg.Port)
	logger.Fatal(server.Start(cfg.Port))
}

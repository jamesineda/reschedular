package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	db2 "github.com/jamesineda/reschedular/app/db"
	"github.com/jamesineda/reschedular/app/event"
	"github.com/jamesineda/reschedular/app/queue"
	"github.com/jamesineda/reschedular/app/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	ERROR      = "ERROR"
	ConfigPath = "CONFIG_PATH"
	SqsQueue   = "SQS_QUEUE"
)

func BindCommandLineArgs() {
	flag.String(ConfigPath, "config/development.yml", "path to config file")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(fmt.Errorf("failed to bind commandline %s", err))
	}
}

func HandleRequest(ctx context.Context, e event.IncomingEvent) (string, error) {
	if err := e.HandleEvent(ctx); err != nil {
		log.Fatalf("received unhandled incoming event: %s", e.FunctionName())
		return ERROR, err
	}

	return fmt.Sprintf("Hello %s!", e.FunctionName()), nil
}

func main() {
	BindCommandLineArgs()
	queueUrl := viper.GetString(SqsQueue)
	configPath := viper.GetString(ConfigPath)

	config, err := utils.NewConfig(configPath)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to parse configPath %s", err))
		return
	}

	db, err := db2.NewDatabaseConn(config.Database)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to establish connection to database %s", err))
		return
	}

	// Not used AWS SQS before, so I'm going to assume the default config store is fine to use for this demo, as per the docs
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sqs.New(sess)

	eventsQueue := queue.NewEventsQueue(100)

	var wg sync.WaitGroup
	//lambdaChannel := make(chan bool) // would be used if lambda.Start() blocks
	sqsQueueChannel := make(chan bool)
	sigC := make(chan os.Signal)

	// set the database on the context
	ctx := context.Background()
	context.WithValue(ctx, "db", db)
	context.WithValue(ctx, "timer", &utils.RealTimer{})
	context.WithValue(ctx, "idGenny", &utils.UUIDID{})
	context.WithValue(ctx, "scs", &svc)
	context.WithValue(ctx, "scsQueueUrl", &queueUrl)
	context.WithValue(ctx, "eventsQueue", &eventsQueue)
	event.StartAsynchronousEventProcessor(ctx, &wg, sqsQueueChannel)

	// I'm not really sure how lambda.Start() behaves, so I'm making the huge assumption that is doesn't block due to
	// the lack of a Stop() or Close() like function exposed. If it DOES block, then I would move the function call into
	// a go routine and pass the waitGroup and a channel to the handler, so that I can shut down the process on a OS
	// interrupt.
	lambda.Start(HandleRequest)

	// wait here until a TERM signal is received
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
	<-sigC
	signal.Stop(sigC)

	// shuts down the Event process queue
	sqsQueueChannel <- true

	log.Println("Reschedular service has shutdown.")
}

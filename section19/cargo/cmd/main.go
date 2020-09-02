package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/longjoy/micro-go-course/section19/cargo/component"
	"github.com/longjoy/micro-go-course/section19/cargo/endpoint"
	"github.com/longjoy/micro-go-course/section19/cargo/transport"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"gopkg.in/mgo.v2"

	"github.com/longjoy/micro-go-course/section19/cargo/dao/inmem"
	"github.com/longjoy/micro-go-course/section19/cargo/dao/mongo"
	"github.com/longjoy/micro-go-course/section19/cargo/inspection"
	shipping "github.com/longjoy/micro-go-course/section19/cargo/model"
	"github.com/longjoy/micro-go-course/section19/cargo/service/booking"
	"github.com/longjoy/micro-go-course/section19/cargo/service/handling"
)

const (
	defaultPort       = "8080"
	defaultHost       = "localhost"
	defaultMongoDBURL = "mongodb://129.211.63.96:27117"
	defaultDBName     = "cargo"
)

func main() {
	var (
		port   = envString("PORT", defaultPort)
		addr   = envString("HOST", defaultHost)
		dburl  = envString("MONGODB_URL", defaultMongoDBURL)
		dbname = envString("DB_NAME", defaultDBName)

		//consulHost = flag.String("consul.host", "114.67.98.210", "consul server ip address")
		consulHost = flag.String("consul.host", "106.15.233.99", "consul server ip address")
		consulPort = flag.String("consul.port", "8500", "consul server port")

		serviceHost  = flag.String("service.host", addr, "service ip address")
		servicePort  = flag.String("service.port", port, "service port")
		mongoDBURL   = flag.String("db.url", dburl, "MongoDB URL")
		databaseName = flag.String("db.name", dbname, "MongoDB database name")
		inmemory     = flag.Bool("inmem", false, "use in-memory repositories")
	)

	flag.Parse()
	ctx := context.Background()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// Setup repositories
	var (
		cargos         shipping.CargoRepository
		locations      shipping.LocationRepository
		voyages        shipping.VoyageRepository
		handlingEvents shipping.HandlingEventRepository
	)

	if *inmemory {
		cargos = inmem.NewCargoRepository()
		locations = inmem.NewLocationRepository()
		voyages = inmem.NewVoyageRepository()
		handlingEvents = inmem.NewHandlingEventRepository()
	} else {
		session, err := mgo.Dial(*mongoDBURL)
		if err != nil {
			panic(err)
		}
		defer session.Close()

		session.SetMode(mgo.Monotonic, true)

		cargos, _ = mongo.NewCargoRepository(*databaseName, session)
		locations, _ = mongo.NewLocationRepository(*databaseName, session)
		voyages, _ = mongo.NewVoyageRepository(*databaseName, session)
		handlingEvents = mongo.NewHandlingEventRepository(*databaseName, session)
	}

	// Configure some questionable dependencies.
	var (
		handlingEventFactory = shipping.HandlingEventFactory{
			CargoRepository:    cargos,
			VoyageRepository:   voyages,
			LocationRepository: locations,
		}
		handlingEventHandler = handling.NewEventHandler(
			inspection.NewService(cargos, handlingEvents, nil),
		)
	)

	// Facilitate testing by adding some cargos.
	storeTestData(cargos)

	//fieldKeys := []string{"method"}

	var bs booking.Service
	bs = booking.NewService(cargos, locations, handlingEvents)
	bs = booking.NewLoggingService(log.With(logger, "component", "booking"), bs)

	var hs handling.Service
	hs = handling.NewService(handlingEvents, handlingEventFactory, handlingEventHandler)
	hs = handling.NewLoggingService(log.With(logger, "component", "handling"), hs)

	endpoints := &endpoint.CargoEndpoints{
		endpoint.MakeBookCargoEndpoint(bs),
		endpoint.MakeLoadCargoEndpoint(bs),
		endpoint.AssignCargoToRouteEndpoint(bs),
		endpoint.ChangeDestinationEndpoint(bs),
		endpoint.CargosEndpoint(bs),
		endpoint.LocationsEndpoint(bs),
		endpoint.RegisterHandlingEventEndpoint(hs),
	}

	r := transport.MakeHttpHandler(ctx, endpoints)

	//创建注册对象
	//TODO replace with pkg consul
	registar := component.Register(*consulHost, *consulPort, *serviceHost, *servicePort, logger)

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", servicePort, "msg", "listening")
		//启动前执行注册
		registar.Register()
		errs <- http.ListenAndServe(":"+*servicePort, r)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	error := <-errs

	//服务退出取消注册
	registar.Deregister()
	logger.Log("terminated", error)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func storeTestData(r shipping.CargoRepository) {
	test1 := shipping.NewCargo("FTL456", shipping.RouteSpecification{
		Origin:          shipping.AUMEL,
		Destination:     shipping.SESTO,
		ArrivalDeadline: time.Now().AddDate(0, 0, 7),
	})
	if _, err := r.Store(test1); err != nil {
		panic(err)
	}

	test2 := shipping.NewCargo("ABC123", shipping.RouteSpecification{
		Origin:          shipping.SESTO,
		Destination:     shipping.CNHKG,
		ArrivalDeadline: time.Now().AddDate(0, 0, 14),
	})
	if _, err := r.Store(test2); err != nil {
		panic(err)
	}
}

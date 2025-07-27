package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/febriandani/ecommerce-be-system-service/cmd/server/api"
	"github.com/febriandani/ecommerce-be-system-service/cmd/server/db"
	dbinfra "github.com/febriandani/ecommerce-be-system-service/cmd/server/infra"
	proto "github.com/febriandani/ecommerce-be-system-service/protogen/golang/proto/system"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgrpc"
	"go.elastic.co/apm/module/apmhttp"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func initConfig() {
	viper.SetConfigName("app")                  // file: app.yml
	viper.AddConfigPath("./cmd/server/config/") // path to app.yml
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
}

func getConfigDB() *dbinfra.DatabaseSystem {
	return &dbinfra.DatabaseSystem{
		ReadUser: dbinfra.DBSystem{
			Username:     viper.GetString("DATABASE.READ.USERNAME"),
			Password:     viper.GetString("DATABASE.READ.PASSWORD"),
			URL:          viper.GetString("DATABASE.READ.URL"),
			Port:         viper.GetString("DATABASE.READ.PORT"),
			DBName:       viper.GetString("DATABASE.READ.DB_NAME"),
			MaxIdleConns: viper.GetInt("DATABASE.READ.MAXIDLECONNS"),
			MaxOpenConns: viper.GetInt("DATABASE.READ.MAXOPENCONNS"),
			MaxLifeTime:  viper.GetInt("DATABASE.READ.MAXLIFETIME"),
			Timeout:      viper.GetString("DATABASE.READ.TIMEOUT"),
			SSLMode:      viper.GetString("DATABASE.READ.SSL_MODE"),
		},
		WriteUser: dbinfra.DBSystem{
			Username:     viper.GetString("DATABASE.WRITE.USERNAME"),
			Password:     viper.GetString("DATABASE.WRITE.PASSWORD"),
			URL:          viper.GetString("DATABASE.WRITE.URL"),
			Port:         viper.GetString("DATABASE.WRITE.PORT"),
			DBName:       viper.GetString("DATABASE.WRITE.DB_NAME"),
			MaxIdleConns: viper.GetInt("DATABASE.WRITE.MAXIDLECONNS"),
			MaxOpenConns: viper.GetInt("DATABASE.WRITE.MAXOPENCONNS"),
			MaxLifeTime:  viper.GetInt("DATABASE.WRITE.MAXLIFETIME"),
			Timeout:      viper.GetString("DATABASE.WRITE.TIMEOUT"),
			SSLMode:      viper.GetString("DATABASE.WRITE.SSL_MODE"),
		},
	}
}

func main() {
	initConfig()
	logger := dbinfra.NewLogger()

	// Optional: set from ENV or Viper
	apm.DefaultTracer.Service.Name = viper.GetString("APP.NAME")
	apm.DefaultTracer.Service.Environment = viper.GetString("APP.ENV")

	// Database connection
	dbRead := dbinfra.NewDB(logger)
	dbRead.ConnectDB(&getConfigDB().ReadUser)
	if dbRead.Err != nil {
		logger.Fatalf("Error connecting to Read DB: %v", dbRead.Err)
	}

	dbWrite := dbinfra.NewDB(logger)
	dbWrite.ConnectDB(&getConfigDB().WriteUser)
	if dbWrite.Err != nil {
		logger.Fatalf("Error connecting to Write DB: %v", dbWrite.Err)
	}

	dbList := &dbinfra.DatabaseList{
		Backend: dbinfra.DatabaseType{
			Read:  dbRead,
			Write: dbWrite,
		},
	}
	dbConfig := db.NewDatabaseConfig(dbList, logger)

	// Inisialisasi server
	grpcPort := viper.GetInt("APP.PORT_SERVER")
	restPort := viper.GetString("APP.PORT_CLIENT")

	// Listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	// gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(apmgrpc.NewUnaryServerInterceptor()),
		grpc.StreamInterceptor(apmgrpc.NewStreamServerInterceptor()),
	)
	systemServer := api.NewSystemServer(dbConfig, logger)
	proto.RegisterSystemsServer(grpcServer, systemServer)

	// Jalankan REST gateway
	go startGatewayServer(fmt.Sprintf("localhost:%d", grpcPort), restPort)

	logger.Infof("gRPC server started at :%d", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func startGatewayServer(grpcAddr string, restPort string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption("*", &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := proto.RegisterSystemsHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway handler: %v", err)
	}

	// Basic auth + Elastic APM Middleware
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !isValidUser(username, password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		mux.ServeHTTP(w, r)
	})

	apmWrapped := apmhttp.Wrap(baseHandler) // ✅ wrap REST gateway

	log.Printf("REST gateway started at :%s", restPort)
	if err := http.ListenAndServe(":"+restPort, apmWrapped); err != nil {
		log.Fatalf("Failed to serve REST gateway: %v", err)
	}
}

// Contoh validasi username dan password
func isValidUser(username, password string) bool {
	// Kamu bisa ganti ini dengan pengecekan ke database, env, atau secret manager
	return username == "admin" && password == "secret123"
}

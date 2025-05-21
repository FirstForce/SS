package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"mqtt-streaming-server/broker"
	"mqtt-streaming-server/routes"
)

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := os.ReadFile("/run/secrets/ca.crt")
	if err != nil {
		panic(err)
	}
	certpool.AppendCertsFromPEM(pemCerts)

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair("/run/secrets/web.crt", "/run/secrets/web.key")
	if err != nil {
		panic(err)
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: certpool,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

func main() {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s@mongo-db:27017/?authSource=admin", os.Getenv("MONGO_INITDB_ROOT_USERNAME"), os.Getenv("MONGO_INITDB_ROOT_PASSWORD"))
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		panic(err)
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	db := mongoClient.Database("mqtt-streaming-server")

	fmt.Println("Connected to MongoDB!")

	// Initialize user routes
	handler := routes.InitRoutes(db)

	go func() {
		fmt.Println("Starting HTTP server on port 8080...")
		if err := http.ListenAndServe("0.0.0.0:8080", handler); err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	brokerHandler := broker.NewBrokerHandler(db)

	tlsconfig := NewTLSConfig()

	opts := mqtt.NewClientOptions()
	opts.AddBroker("ssl://broker:8883")
	opts.SetClientID("web").SetTLSConfig(tlsconfig)

	// Start the connection
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to a Topic
	if token := client.Subscribe("photos/#", 0, brokerHandler.HandlePhoto); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe("register/#", 0, brokerHandler.RegisterDevice); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	<-c
}

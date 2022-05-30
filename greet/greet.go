package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"

	"greet/internal/config"
	"greet/internal/handler"
	"greet/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	_ "github.com/joho/godotenv"
)

var configFile = flag.String("f", "etc/greet-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

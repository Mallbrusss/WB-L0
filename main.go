package main

import (
	consumer "L0/Consumer"
	dbm "L0/DataBaseManager"
	httpServ "L0/HttpServer"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dbm.Connect()
	consumer.Consume()
	httpServ.Serv()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	
	consumer.Disconnect()
	dbm.Disconnect()
	fmt.Println("Завершение программы...")
}

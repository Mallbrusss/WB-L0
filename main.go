package main

import (
	consumer "L0/Consumer"
	dbm "L0/DataBaseManager"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dbm.DbConnect()
	consumer.Consume()

	// Ждем сигнала для завершения программы (например, Ctrl+C).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	consumer.Disconnect()
	dbm.DbDisconnect()
	fmt.Println("Завершение программы...")
}

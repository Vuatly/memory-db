package main

import (
	"bufio"
	"fmt"
	"memory-db/cmd"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	app, err := cmd.NewApp()
	if err != nil {
		panic(err)
	}

	go signalsHandler()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Database is ready")
	for {
		scanner.Scan()
		if err = scanner.Err(); err != nil {
			app.Logger.Error("error reading input", zap.Error(err))
			continue
		}

		query := scanner.Text()
		output := app.Database.HandleQuery(query)
		fmt.Println(output)
	}
}

func signalsHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigs
	fmt.Printf("received signal: %s. shutting down...", sig)
	signal.Stop(sigs)
	close(sigs)

	os.Exit(0)
}

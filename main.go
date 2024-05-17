package main

import (
	"context"
	"fmt"
	"github.com/LukeWinikates/homekit-alertsink/alertmanger"
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	alertManager := accessory.NewSecuritySystem(accessory.Info{Name: "AlertManager"})
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./db"
	}
	fs := hap.NewFsStore(dataDir)

	server, err := hap.NewServer(fs, alertManager.A)
	if err != nil {
		log.Panic(err)
	}
	server.Addr = "localhost:9009"

	newHandlerFunc := alertmanger.NewHandlerFunc(func(payload *alertmanger.WebhookPayload) error {
		status, err := convertStatus(payload.Status)
		if err != nil {
			return err
		}
		return alertManager.SecuritySystem.SecuritySystemCurrentState.SetValue(status)
	})
	server.ServeMux().HandleFunc("/alerts", newHandlerFunc)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		signal.Stop(c)
		cancel()
	}()
	log.Fatal(server.ListenAndServe(ctx))
}

func convertStatus(status string) (int, error) {
	switch status {
	case "resolved":
		return characteristic.SecuritySystemCurrentStateStayArm, nil
	case "alerting":
		return characteristic.SecuritySystemCurrentStateAlarmTriggered, nil
	default:
		return characteristic.SecuritySystemCurrentStateDisarmed,
			fmt.Errorf("unexpected webhook status (not 'alerting' or 'resolved'): %s", status)
	}
}

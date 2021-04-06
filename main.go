package main

import (
	"CallerIdBot/engine"
	"fmt"
	"log"
	"time"
)

func main() {



	engine.Conf.Init()
	log.Println("Configuration file init OK")
	fmt.Println("Cursors ->", engine.DevicesCursor)
	fmt.Println("Devices ->", engine.Devices)
	fmt.Println("Customers ->", engine.Customers)
	engine.Redis.Init()
	log.Println("Redit init OK")



	engine.SMSch = make(chan int)
	engine.DCch = make(chan int)
	//engine.Redis.ClearByKeys([]string{"SMSBOTFROMTX","SMSBOTFROMDC"})
	//log.Println("All zombie records cleared")


	go engine.ListenSMS()
	log.Println("ListenSMS started")
	go engine.ListenDC()
	log.Println("ListenDC started")
	go engine.RefrechDCKey()
	log.Println("RefrechDCKey started")

	for {
		select {
		case <-engine.SMSch:
			//smsRq.Gateway, smsRq.Number, "CID: Testing in progress...", "", "SMS",
			log.Println("SMS incoming message ->", engine.SMS)
			engine.TelnyxIncoming()
			break
		case <-engine.DCch:
			log.Println("DC incoming message ->", engine.DC)
			engine.DCIncoming()
			break
		default:
			time.Sleep(time.Second)
			fmt.Print(".")

		}

	}

}

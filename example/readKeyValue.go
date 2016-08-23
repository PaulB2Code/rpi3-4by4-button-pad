package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/PaulB2Code/rpi3-4by4-button-pad"
)

func main() {
	log.Println("Start Read Key Value Programm")

	kp, err := KeyPad.New()
	if err != nil {
		log.Panic("[ERROR - FATAL], Impossible de démarrer, ", err)
	}

	defer kp.Close()

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, os.Interrupt)
	go func() {
		for _ = range closeSignal {
			log.Println("\nClose Display.\n")
			kp.Close()
			os.Exit(0)
		}
	}()

	c := make(chan int)
	kp.TrackClicked(100*time.Millisecond, c)
	defer kp.StopTracking()

	for {
		val := <-c
		kp.Reading = false

		letter, err := kp.GetValueWithColumn(val)
		if err != nil {
			log.Println("[ERROR], ", err)
		}
		log.Println("Letter clicked", letter)

	}

}

package undertaker

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func NewUndertaker(log logrus.FieldLogger) *Undertaker {
	return &Undertaker{
		wishlist:  make([]func() error, 0),
		errorChan: make(chan error),
		log:       log.WithField("service", "Undertaker"),
	}
}

// Handle the app final wishlist before taking it down.
type Undertaker struct {
	wishlist  []func() error
	errorChan chan error
	log       logrus.FieldLogger
}

func (u *Undertaker) Add(fn func() error) {
	u.wishlist = append(u.wishlist, fn)
}

func (u *Undertaker) WaitForDeath() {
	sig := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	var sigReceived uint

	// Handle SIGINT and SIGTERM (CTRL + C)
	go func() {
		for {
			select {
			case s := <-sig:
				u.log.Info("Signal received: %v", s)
				sigReceived++

				if sigReceived < 2 {
					// Quit gracefully
					u.log.Info("Executing the wishlist.")
					for _, fn := range u.wishlist {
						if err := fn(); err != nil {
							u.errorChan <- err
						}
					}

					u.log.Info("Done")
					done <- true
				} else {
					u.log.Info("Too many attempts, please wait.")
				}
			}
		}
	}()

	<-done
}

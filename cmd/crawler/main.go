package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/seggga/gb_best_go1/internal/crawler"
	"github.com/seggga/gb_best_go1/internal/domain"
	"github.com/seggga/gb_best_go1/internal/requester"
	"github.com/seggga/gb_best_go1/internal/service"
)

func main() {
	// read config file
	cfg, err := ReadConfig()
	if err != nil {
		log.Printf("could not read config file: %v", err)
		return
	}

	r, err := requester.NewRequester(cfg.ReqTimeout, nil)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	cr, err := crawler.NewCrawler(r, cfg.MaxDepth)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	srv, err := service.NewService(*cfg, cr, cancel, ctx)
	if err != nil {
		log.Printf("error starting service, %v", err)
		return
	}
	log.Printf("Crawler started with PID: %d", os.Getpid())

	srv.Run()

	sigInt := make(chan os.Signal) // Создаем канал для приема сигналов // gocritic
	// Подписываемся на сигнал SIGINT // gocritic
	signal.Notify(sigInt, syscall.SIGINT) //nolint:govet // syscall.SIGINT implements required interface

	sigUsr := make(chan os.Signal) // Создаем канал для приема сигналов // gocritic
	// Подписываемся на сигнал SIGUSR1 // gocritic
	signal.Notify(sigUsr, syscall.SIGUSR1) //nolint:govet // syscall.SIGUSR1 implements required interface

	var next = true
	for next {
		select {
		case <-ctx.Done(): // Если всё завершили - выходим // gocritic
			next = false

		// got INT signal
		case <-sigInt:
			log.Println("got INTERRUPT signal")
			cancel() // Если пришёл сигнал SigInt - завершаем контекст // gocritic

		// total timeout
		case <-time.After(time.Second * time.Duration(cfg.CrawlTimeout)):
			log.Printf("Crawler stops on timeout: %d sec", cfg.CrawlTimeout)
			cancel()

		// add 2 to max depth
		case <-sigUsr:
			log.Println("got USR1 signal")
			srv.IncreaseDepth() // sigUsr1 - increase maxDepth
		}
	}
	log.Println("program exit")
}

// ReadConfig implements filling config from yaml-file
func ReadConfig() (*domain.Config, error) {
	// read config file
	configData, err := ioutil.ReadFile("./configs/config.yaml")
	if err != nil {
		return nil, err
	}

	// decode config
	cfg := new(domain.Config)
	err = yaml.Unmarshal(configData, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil

}

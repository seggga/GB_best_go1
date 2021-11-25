package main

import (
	"context"
	"io/ioutil"
	"lesson1/internal/requester"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/seggga/GB_best_go1/internal/crawler"
)

//Config - структура для конфигурации
type Config struct {
	MaxDepth     uint64 `yaml:"maxdepth"`
	MaxResults   int    `yaml:"maxresults"`
	MaxErrors    int    `yaml:"maxerrors"`
	Url          string `yaml:"url"`
	ReqTimeout   int    `yaml:"reqtimeout"`
	CrawlTimeout int    `yaml:"crawltimeout"`
}

// ReadConfig implements filling config from yaml-file
func ReadConfig() (*Config, error) {
	// read config file
	configData, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	// decode config
	cfg := new(Config)
	err = yaml.Unmarshal(configData, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	// read config file
	cfg, err := ReadConfig()
	if err != nil {
		log.Printf("could not read config file: %v", err)
		return
	}

	var cr Crawler
	var r Requester

	r = requester.NewRequester(time.Duration(cfg.ReqTimeout)*time.Second, nil)
	cr = crawler.NewCrawler(r, cfg.MaxDepth)
	log.Printf("Crawler started with PID: %d", os.Getpid())

	ctx, cancel := context.WithCancel(context.Background())
	go cr.Scan(ctx, cfg.Url, 1)             //Запускаем краулер в отдельной рутине
	go processResult(ctx, cancel, cr, *cfg) //Обрабатываем результаты в отдельной рутине

	sigInt := make(chan os.Signal)        //Создаем канал для приема сигналов
	signal.Notify(sigInt, syscall.SIGINT) //Подписываемся на сигнал SIGINT

	sigUsr := make(chan os.Signal)         //Создаем канал для приема сигналов
	signal.Notify(sigUsr, syscall.SIGUSR1) //Подписываемся на сигнал SIGUSR1

	for {
		select {
		case <-ctx.Done(): //Если всё завершили - выходим
			return

		// got INT signal
		case <-sigInt:
			log.Println("got INTERRUPT signal")
			cancel() //Если пришёл сигнал SigInt - завершаем контекст

		// total timeout
		case <-time.After(time.Second * time.Duration(cfg.CrawlTimeout)):
			log.Printf("Crawler stops on timeout: %d sec", cfg.CrawlTimeout)
			cancel()

		// add 2 to max depth
		case <-sigUsr:
			log.Println("got USR1 signal")
			cr.IncreaseDepth() // sigUsr1 - increase maxDepth
		}
	}
}

func processResult(ctx context.Context, cancel func(), cr Crawler, cfg Config) {
	var maxResult, maxErrors = cfg.MaxResults, cfg.MaxErrors
	for {
		select {
		case <-ctx.Done():
			return

		// got message in the channel
		case msg := <-cr.ChanResult():
			if msg.Err != nil {
				// message contains error
				maxErrors--
				log.Printf("crawler result return err: %s\n", msg.Err.Error())
				if maxErrors <= 0 {
					log.Println("Maximum number of errors occured.")
					cancel()
					return
				}
			} else {
				// message contains data
				maxResult--
				log.Printf("crawler result: [url: %s] Title: %s\n", msg.Url, msg.Title)
				if maxResult <= 0 {
					log.Println("Maximum number of results obtained.")
					cancel()
					return
				}
			}
		}
	}
}

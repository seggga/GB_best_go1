package service

import (
	"context"
	"log"

	"github.com/seggga/gb_best_go1/internal/domain"
)

type Service struct {
	config  domain.Config
	crawler domain.Crawler
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewService creates a structure Service
func NewService(cfg domain.Config, cr domain.Crawler, cancel context.CancelFunc, ctx context.Context) (*Service, error) {

	// TODO - validate config
	return &Service{
		config:  cfg,
		crawler: cr,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (s *Service) Run() {
	go s.crawler.Scan(s.ctx, s.config.URL, 1)              //Запускаем краулер в отдельной рутине
	go processResult(s.ctx, s.cancel, s.crawler, s.config) //Обрабатываем результаты в отдельной рутине
}

func (s *Service) IncreaseDepth() {
	s.crawler.IncreaseDepth()
}

func processResult(ctx context.Context, cancel func(), cr domain.Crawler, cfg domain.Config) {
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

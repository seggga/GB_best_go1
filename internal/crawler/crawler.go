package crawler

import (
	"context"
	"log"
	"sync"
	"sync/atomic"

	"github.com/seggga/gb_best_go1/internal/domain"
)

// crawler is a structure that has all the items to control whole process
type crawler struct {
	r        domain.Requester        // a thing that queries pages
	res      chan domain.CrawlResult // a channel to pass results from r
	visited  map[string]struct{}     // a map to hold visited URLs
	mu       sync.RWMutex            // a mutex to share "visited"-map between multibple go-routines
	maxDepth uint64                  // limits scanning depth
}

// NewCrawler creates a new crawler structure
func NewCrawler(r domain.Requester, maxDepth uint64) *crawler {
	return &crawler{
		r:        r,
		res:      make(chan domain.CrawlResult),
		visited:  make(map[string]struct{}),
		mu:       sync.RWMutex{},
		maxDepth: maxDepth,
	}
}

// Scan fills crawler's map with visited URLs and calls Get-method to scan webpages
func (c *crawler) Scan(ctx context.Context, url string, depth uint64) {
	//Проверяем то, что есть запас по глубине
	c.mu.RLock()
	maxDepthAchieved := depth > c.maxDepth
	c.mu.RUnlock()
	if maxDepthAchieved {
		return
	}
	c.mu.RLock()
	_, ok := c.visited[url] //Проверяем, что мы ещё не смотрели эту страницу
	c.mu.RUnlock()
	if ok {
		return
	}
	select {
	case <-ctx.Done(): //Если контекст завершен - прекращаем выполнение
		return
	default:
		page, err := c.r.Get(ctx, url) //Запрашиваем страницу через Requester
		if err != nil {
			c.res <- domain.CrawlResult{Err: err} //Записываем ошибку в канал
			return
		}
		c.mu.Lock()
		c.visited[url] = struct{}{} //Помечаем страницу просмотренной
		c.mu.Unlock()
		c.res <- domain.CrawlResult{ //Отправляем результаты в канал
			Title: page.GetTitle(),
			Url:   url,
		}
		for _, link := range page.GetLinks() {
			go c.Scan(ctx, link, depth+1) //На все полученные ссылки запускаем новую рутину сборки
		}
	}
}

func (c *crawler) ChanResult() <-chan domain.CrawlResult {
	return c.res
}

// IncreaseDpeth adds 2 to the property 'maxDepth' atomically
func (c *crawler) IncreaseDepth() {
	newDepth := atomic.AddUint64(&c.maxDepth, 2)
	log.Printf("MaxDepth increased via SIGUSR1, new value is %d", newDepth)
}

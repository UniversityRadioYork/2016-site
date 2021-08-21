package models

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/UniversityRadioYork/2016-site/structs"
	"log"
	"net"
	"sync"
	"time"
)

type ShortURL struct {
	ShortUrlID uint   `toml:"id"`
	Slug       string `toml:"slug"`
	RedirectTo string `toml:"redirectTo"`
}

type ShortURLModel struct {
	c          *structs.Config
	urlsBySlug map[string]*ShortURL
	urlsLock   sync.RWMutex
}

func NewShortURLsModel(cfg *structs.Config) *ShortURLModel {
	return &ShortURLModel{
		c:          cfg,
		urlsBySlug: make(map[string]*ShortURL),
		urlsLock:   sync.RWMutex{},
	}
}

func (m *ShortURLModel) Match(slug string) *ShortURL {
	m.urlsLock.RLock()
	defer m.urlsLock.RUnlock()
	if url, ok := m.urlsBySlug[slug]; ok {
		return url
	} else {
		return nil
	}
}

func (m *ShortURLModel) TrackClick(slug string, visitorIP *net.IPAddr) error {
	// TODO
	return nil
}

func (m *ShortURLModel) doTickUpdate() {
	var decoded struct {
		URLs []ShortURL `toml:"urls"`
	}
	if _, err := toml.DecodeFile("short_urls.toml.example", &decoded); err != nil {
		log.Println(fmt.Errorf("while parsing short_urls: %w", err))
		return
	}
	indexed := make(map[string]*ShortURL)
	for _, url := range decoded.URLs {
		indexed[url.Slug] = &url
	}

	m.urlsLock.Lock()
	defer m.urlsLock.Unlock()
	m.urlsBySlug = indexed
	log.Println("Short URLs updated successfully")
}

func (m *ShortURLModel) UpdateTimer() {
	if m.c.ShortURLs.UpdateInterval == 0 {
		log.Fatal("Tried to UpdateTimer but update interval is zero!")
	}
	for {
		m.doTickUpdate()
		time.Sleep(time.Duration(m.c.ShortURLs.UpdateInterval) * time.Second)
	}
}

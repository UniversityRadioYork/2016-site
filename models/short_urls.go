package models

import (
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/myradio-go"
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
	s          *myradio.Session
	urlsBySlug map[string]*ShortURL
	urlsLock   sync.RWMutex
}

func NewShortURLsModel(c *structs.Config, s *myradio.Session) *ShortURLModel {
	return &ShortURLModel{
		c:          c,
		s:          s,
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

func (m *ShortURLModel) TrackClick(id uint, visitorUserAgent string, visitorIP net.IP) error {
	ip := ""
	if visitorIP != nil {
		ip = visitorIP.String()
	}
	return m.s.LogShortURLClick(id, visitorUserAgent, ip)
}

func (m *ShortURLModel) doTickUpdate() {
	urls, err := m.s.GetAllShortURLs()
	if err != nil {
		log.Printf("when getting short URLs: %w", err)
		return
	}

	indexed := make(map[string]*ShortURL)
	for _, url := range urls {
		indexed[url.Slug] = &ShortURL{
			ShortUrlID: url.ShortURLID,
			Slug:       url.Slug,
			RedirectTo: url.RedirectTo,
		}
	}

	m.urlsLock.Lock()
	defer m.urlsLock.Unlock()
	m.urlsBySlug = indexed
	log.Println("Short URLs updated successfully")
}

func (m *ShortURLModel) UpdateTimer() {
	if m.c.ShortURLs.UpdateInterval == 0 {
		panic("Tried to UpdateTimer but update interval is zero!")
	}
	for {
		m.doTickUpdate()
		time.Sleep(time.Duration(m.c.ShortURLs.UpdateInterval) * time.Second)
	}
}

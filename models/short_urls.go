package models

import (
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/myradio-go"
	"log"
	"net"
	"sync"
	"time"
)

type ShortURLModel struct {
	c          *structs.Config
	s          *myradio.Session
	urlsBySlug map[string]*myradio.ShortURL
	urlsLock   sync.RWMutex
}

func NewShortURLsModel(c *structs.Config, s *myradio.Session) *ShortURLModel {
	return &ShortURLModel{
		c:          c,
		s:          s,
		urlsBySlug: make(map[string]*myradio.ShortURL),
		urlsLock:   sync.RWMutex{},
	}
}

func (m *ShortURLModel) Match(slug string) *myradio.ShortURL {
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
		log.Printf("when getting short URLs: %v", err)
		return
	}

	indexed := make(map[string]*myradio.ShortURL)
	for idx, url := range urls {
		indexed[url.Slug] = &urls[idx]
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

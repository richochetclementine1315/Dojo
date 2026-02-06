package service

import (
	"log"
	"time"
)

// ContestSyncService handles periodic contest synchronization
type ContestSyncService struct {
	contestService *ContestService
	ticker         *time.Ticker
	stopChan       chan bool
}

// NewContestSyncService creates a new contest sync service
func NewContestSyncService(contestService *ContestService) *ContestSyncService {
	return &ContestSyncService{
		contestService: contestService,
		stopChan:       make(chan bool),
	}
}

// Start begins the periodic contest synchronization
func (s *ContestSyncService) Start(interval time.Duration) {
	log.Println("Starting contest sync service...")

	// Sync immediately on start
	s.syncContests()

	// Then sync periodically
	s.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.syncContests()
			case <-s.stopChan:
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop halts the periodic synchronization
func (s *ContestSyncService) Stop() {
	log.Println("Stopping contest sync service...")
	s.stopChan <- true
}

// syncContests fetches and stores contests from all platforms
func (s *ContestSyncService) syncContests() {
	log.Println("Syncing contests from all platforms...")

	count, err := s.contestService.SyncContestsFromPlatform("all")
	if err != nil {
		log.Printf("Error syncing contests: %v\n", err)
		return
	}

	log.Printf("Successfully synced %d contests\n", count)
}

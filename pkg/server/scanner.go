package server

import (
	"time"

	"go.uber.org/zap"
)

func (s *Server) runScanner() {
	ticker := time.NewTicker(s.ScanPeriod)

	for {
		select {
		case <-ticker.C:
			s.logger.Debug("Scanning")
			endpoints, err := s.allEndpoints()
			if err != nil {
				s.logger.Error("Failed to get endpoints", zap.Error(err))
			}
			s.logger.Debug("Scan finished", zap.Any("endpoints", endpoints))
			err = s.db.Store(time.Now().Unix(), endpoints)
			if err != nil {
				s.logger.Error("Failed to save snapshot", zap.Error(err))
			}
		case <-s.ctx.Done():
			ticker.Stop()
		}
	}
}

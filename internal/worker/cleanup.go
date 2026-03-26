package worker

import (
	"QuickSlot/pkg/database/mysql"
	"context"
	"log"
	"time"
)

// CleanExpiredSlots deletes unbooked slots that are in the past.
// It runs in a goroutine with the given interval and stops when done channel is closed.
func CleanExpiredSlots(db *mysql.Dialect, interval time.Duration, done <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Println("background worker started: expired slot cleanup")

	for {
		select {
		case <-done:
			log.Println("background worker stopped: expired slot cleanup")
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			res, err := db.DB.ExecContext(ctx,
				"DELETE FROM time_slots WHERE is_booked = FALSE AND end_time < NOW()")
			cancel()

			if err != nil {
				log.Printf("expired slot cleanup error: %v", err)
				continue
			}

			rows, _ := res.RowsAffected()
			if rows > 0 {
				log.Printf("cleaned %d expired slots", rows)
			}
		}
	}
}

package sessions

import "time"

type Session struct {
	ID        string                 // Session ID
	UserID    string                 // Associated User ID
	ExpiresAt time.Time              // Expiration time of the session
	Data      map[string]interface{} // Additional session data
}

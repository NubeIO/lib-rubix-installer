package installer

import (
	"time"
)

// timestamp *
func timestamp() string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return t
}

package types

import (
	"fmt"
	"time"
)

// FormatTime formats the given time.
func FormatTime(t time.Time) string {
	return fmt.Sprint("'", t.UTC().Format("2006-01-02 15:04:05.999999"), "+00'")
}

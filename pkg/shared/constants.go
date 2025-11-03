package shared

import "time"

// Timeout constants for database operations
const (
	// TimeoutForReadOperation - timeout for read operations (Find)
	TimeoutForReadOperation = 5 * time.Second
	// TimeoutForWriteOperation - timeout for write operations (Create, Update)
	TimeoutForWriteOperation = 10 * time.Second
	// TimeoutForComplexOperation - timeout for complex operations (multiple DB calls)
	TimeoutForComplexOperation = 15 * time.Second
)

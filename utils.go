package main

import (
	"fmt"
	"time"
)

func generateMessageID() string {

	return fmt.Sprintf(
		"%d",
		time.Now().UnixNano(),
	)
}

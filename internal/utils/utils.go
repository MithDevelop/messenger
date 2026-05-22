package utils

import (
	"fmt"
	"time"
)

func GenerateMessageID() string {

	return fmt.Sprintf(
		"%d",
		time.Now().UnixNano(),
	)
}

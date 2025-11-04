package generate

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("req-%x-%d", b, time.Now().UnixNano())
}

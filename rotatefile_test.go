package rotatefile

import (
	"fmt"
	"testing"
	"time"
)

func TestRotate(t *testing.T) {
	file, _ := New("logs/day.log", PerMinute)

	for i := 0; i < 1000; i++ {
		now := time.Now().Format(time.DateTime)
		_, _ = file.Write([]byte(fmt.Sprintf("%d: %s\n", i, now)))
		fmt.Println(now)
		time.Sleep(time.Second)
	}
}

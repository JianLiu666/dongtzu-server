package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestConvert(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(10 * time.Minute)

	fmt.Println((t2.Unix() - t1.Unix()) / 60)
}

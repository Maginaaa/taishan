package middleware

import (
	"fmt"
	"time"
)

var Location *time.Location

func init() {
	Location, _ = time.LoadLocation("Asia/Shanghai")
	fmt.Println("Location initialized")
}

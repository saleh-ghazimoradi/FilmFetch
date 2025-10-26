package helper

import (
	"fmt"
	"log"
	"sync"
)

var WG sync.WaitGroup

func Background(fn func()) {
	WG.Go(func() {
		defer func() {
			pv := recover()
			if pv != nil {
				log.Println(fmt.Sprintf("%v", pv))
			}
		}()
		fn()
	})
}

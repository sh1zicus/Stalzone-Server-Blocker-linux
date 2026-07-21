package ping

import (
	"sync"
	"time"

	"github.com/sh1zicus/stalzone-server-blocker/internal/model"
)

func Refresh(pools []model.Pool) {

	var wg sync.WaitGroup

	for pi := range pools {

		for ti := range pools[pi].Tunnels {

			wg.Add(1)

			go func(pi, ti int) {

				defer wg.Done()

				ping, ok := TCP(
					pools[pi].Tunnels[ti].Address,
		    1500*time.Millisecond,
				)

				pools[pi].Tunnels[ti].Ping = ping
				pools[pi].Tunnels[ti].PingOK = ok

			}(pi, ti)
		}
	}

	wg.Wait()
}

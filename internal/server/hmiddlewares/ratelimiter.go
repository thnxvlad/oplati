package hmiddlewares

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type TimeReq struct {
	timeRequest time.Time
	countReques int
}

const maxRequestPerInterval = 3

const timeIntervalRequests time.Duration = 60 * time.Second

func RateLimiterMiddleware(next http.Handler) http.Handler {
	var timeRequest map[string]TimeReq = make(map[string]TimeReq)
	var mutex sync.Mutex
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userAgent := r.UserAgent()
		key := fmt.Sprintf("%s:%s", userIP, userAgent)
		mutex.Lock()
		defer mutex.Unlock()
		timeReq, ok := timeRequest[key]
		if !ok || timeReq.timeRequest.Add(timeIntervalRequests).Before(time.Now()){
			timeRequest[key] = TimeReq{
			timeRequest: time.Now(),
			countReques: 1,
			}
		} else {
			if timeReq.countReques < maxRequestPerInterval {
				timeReq.countReques += 1
				timeRequest[key] = timeReq
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
		}
	})
}

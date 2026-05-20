package hmiddlewares

// TODO: реализовать ограничение на количество запросов в секунду по  ip + user-agent

import (
	"fmt"
	"net/http"
	"time"
)

var timeRequest map[string]TimeReq = make(map[string]TimeReq)

type TimeReq struct {
	timeRequest time.Time
	countReques int
}

const maxRequestPerSecond = 3

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(AccountIdContextKey{}).(string)
		userAgent := r.UserAgent()
		key := fmt.Sprintf("%s:%s", userID, userAgent)
		timeReq, ok := timeRequest[key]
		if ok {
			if timeReq.timeRequest.Add(time.Second).Before(time.Now()) {
				timeRequest[key] = TimeReq{
					timeRequest: time.Now(),
					countReques: 1,
				}
			} else {
				if timeReq.countReques < maxRequestPerSecond {
					cnt := timeReq.countReques + 1
					timeRequest[key] = TimeReq{
						timeRequest: time.Now(),
						countReques: cnt,
					}
				} else {
					w.WriteHeader(http.StatusTooManyRequests)
					return
				}
			}
		} else {
			timeRequest[key] = TimeReq{
				timeRequest: time.Now(),
				countReques: 1,
			}
		}
	})
}

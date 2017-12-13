package ratelimit

import (
	"fmt"

	"golang.org/x/net/context"

	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

var (
// TODO: use unique identifier for user
//	a     = "a"
//users = &sync.Map{}
)

// RateLimiter limits the amount of request for a user.
func RateLimiter(ctx context.Context, info *tap.Info) (context.Context, error) {
	/*	if _, ok := users.Load(a); !ok {
			fmt.Println("STORE")
			users.Store(a,
				//rate.NewLimiter(config.RateLimit, config.Burst))
				rate.NewLimiter(config.RateLimit, 100))


		// HOW DO WE EXECUTE AN INTERCEPTOR PER HANDLER???????
		if err := f.Wait(ctx); err != nil {
			fmt.Printf("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA: %v -- %v -- %v\n", f.Limit(), f.Burst(), err)
			return nil, status.Errorf(codes.ResourceExhausted, "service is over rate limit")
		}
		}*/
	t := rate.NewLimiter(1, 1)
	//if u, ok := users.Load(a); ok {
	//		l, ok := u.(rate.Limiter)
	if !t.Allow() {

		return nil, status.Errorf(codes.ResourceExhausted, "service is over rate limit")
	}
	fmt.Printf("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA: %v -- %v\n", t.Limit(), t.Burst())
	// Wait errors out if the request cannot be processed within
	// the deadline. This is preemptive, instead of waiting the
	// entire duration.
	if err := t.Wait(ctx); err != nil {
		fmt.Printf("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA: %v -- %v -- %v\n", t.Limit(), t.Burst(), err)
		return nil, status.Errorf(codes.ResourceExhausted, "service is over rate limit")
	}
	//	}

	return ctx, nil
}

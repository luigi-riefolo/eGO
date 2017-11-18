package metrics

/*
var (
	// How often our /hello request durations fall into one of the defined buckets.
	// We can use default buckets or set ones we are interested in.
	duration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "hello_request_duration_seconds",
		Help:    "Histogram of the /hello request duration.",
		Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	})
	// Counter vector to which we can attach labels. That creates many key-value
	// label combinations. So in our case we count requests by status code separetly.
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hello_requests_total",
			Help: "Total number of /hello requests.",
		},
		[]string{"status"},
	)
)

// init registers Prometheus metrics.
func init() {
	prometheus.MustRegister(duration)
	prometheus.MustRegister(counter)
}

func main() {
	addr := flag.String("http", "0.0.0.0:8000", "HTTP server address")
	flag.Parse()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	var status int

	defer func(begun time.Time) {
		duration.Observe(time.Since(begun).Seconds())

		// hello_requests_total{status="200"} 2385
		counter.With(prometheus.Labels{
			"status": fmt.Sprint(status),
		}).Inc()
	}(time.Now())

	status = doSomeWork()
	w.WriteHeader(status)
	w.Write([]byte("Hello, World!\n"))
}

func doSomeWork() int {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	statusCodes := [...]int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusInternalServerError,
	}
	return statusCodes[rand.Intn(len(statusCodes))]
}
*/

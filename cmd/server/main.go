package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

var metricStorage = MemStorage{
	gauges:   make(map[string]float64),
	counters: make(map[string]int64),
}

func metricUpdate(res http.ResponseWriter, req *http.Request) {
	slashes := strings.Split(req.URL.Path, "/")
	if len(slashes) != 5 {
		http.Error(res, "invalid request path", http.StatusNotFound)
		return
	}
	for i := range slashes {
		slashes[i] = strings.TrimSpace(slashes[i])
	}
	_, metricType, name, value := slashes[1], slashes[2], slashes[3], slashes[4]

	if name == "" || name == " " {
		http.Error(res, "metric name is required", http.StatusNotFound)
		return
	}
	switch metricType {
	case "gauge":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(res, "invalid gauge value", http.StatusBadRequest)
			return
		}
		metricStorage.gauges[name] = val
	case "counter":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(res, "invalid counter value", http.StatusBadRequest)
			return
		}
		metricStorage.counters[name] += val
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, metricUpdate)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

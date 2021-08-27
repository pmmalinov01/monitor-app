// https://medium.com/@gsisimogang/instrumenting-golang-server-in-5-min-c1c32489add3

package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Websitedata struct {
	responseCode float64
	responseTime time.Duration
}

func monitorWebsite(url string) *Websitedata {
	now := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	foo := Websitedata{}
	if resp.StatusCode == 200 {
		foo.responseCode = 0
		foo.responseTime = time.Since(now)
		return &foo
	}
	foo.responseCode = 1
	foo.responseTime = time.Since(now)
	return &foo
}

func VarOrDefault(vName, defValue string) string {
	res, ok := os.LookupEnv(vName)
	if !ok {
		return defValue
	}

	return res

}

var (
	sample_external_url_up = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_external_url_up",
		Help: "Help me if needed.",
	},
		[]string{"url"})

	response_ms_gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_external_url_response_ms",
		Help: "Help me if needed.",
	}, []string{"url"})
)

func init() {
	prometheus.MustRegister(sample_external_url_up)
	prometheus.MustRegister(response_ms_gauge)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	interval, err := strconv.Atoi(VarOrDefault("monitorInterval", "10"))
	url := VarOrDefault("PROM_URL", "https://httpstat.us/")
	if err != nil {
		log.Fatal("Specify monitor interval")
	}
	log.Printf("Starting to monitor [%s], interval [%d].\n", url, interval)

	sCodes := []string{"200", "503"}
	go func() {
		for {
			for i := 0; i < len(sCodes); i++ {

				bar := monitorWebsite(url + sCodes[0])
				okStatus := monitorWebsite(url + sCodes[1])
				sample_external_url_up.WithLabelValues("https://httpstat.us/503").Set(bar.responseCode)
				response_ms_gauge.WithLabelValues("http://httpstat.us/503").Set(float64(bar.responseTime.Milliseconds()))
				sample_external_url_up.WithLabelValues("https://httpstat.us/200").Set(okStatus.responseCode)
				response_ms_gauge.WithLabelValues("http://httpstat.us/200").Set(float64(okStatus.responseTime.Milliseconds()))
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}
	}()

	http.ListenAndServe(":8001", nil)
}

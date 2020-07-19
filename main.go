package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
)

type (
	// Target is target
	Target struct {
		Prefecture string `bigquery:"prefecture"`
		City       string `bigquery:"city"`
		URL        string `bigquery:"url"`
	}
	// Tracking is tracking
	Tracking struct {
		Target
		ResCode   bigquery.NullInt64 `bigquery:"res_code"`
		ResMsec   bigquery.NullInt64 `bigquery:"res_msec"`
		Timestamp time.Time          `bigquery:"timestamp"`
	}
)

var (
	targets = []*Target{
		{
			Prefecture: "yamaguchi",
			City:       "yamaguchi",
			URL:        "https://www.city.yamaguchi.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "shimonoseki",
			URL:        "http://www.city.shimonoseki.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "ube",
			URL:        "https://www.city.ube.yamaguchi.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "hagi",
			URL:        "https://www.city.hagi.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "hofu",
			URL:        "https://www.city.hofu.yamaguchi.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "kudamatsu",
			URL:        "https://www.city.kudamatsu.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "iwakuni",
			URL:        "https://www.city.iwakuni.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "hikari",
			URL:        "https://www.city.hikari.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "yanai",
			URL:        "https://www.city-yanai.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "mine",
			URL:        "http://www2.city.mine.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "nagato",
			URL:        "https://www.city.nagato.yamaguchi.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "shunan",
			URL:        "http://www.city.shunan.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "sanyoonoda",
			URL:        "https://www.city.sanyo-onoda.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "kaminoseki",
			URL:        "https://www.town.kaminoseki.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "tabuse",
			URL:        "https://www.town.tabuse.lg.jp/www/index.html",
		},
		{
			Prefecture: "yamaguchi",
			City:       "hirao",
			URL:        "http://www.town.hirao.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "suooshima",
			URL:        "https://www.town.suo-oshima.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "waki",
			URL:        "http://www.town.waki.lg.jp/",
		},
		{
			Prefecture: "yamaguchi",
			City:       "abu",
			URL:        "http://www.town.abu.lg.jp/",
		},
	}
)

func main() {
	http.HandleFunc("/yamaguchi", handler)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	log.Printf("sever shutdown")
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("gov-site-tracker started\n")

	for _, t := range targets {
		target := t
		log.Printf("gov: %s-%s, check\n", target.Prefecture, target.City)

		tracking, err := track(target)
		if err != nil {
			log.Printf("gov: %s-%s, track err: %s\n", target.Prefecture, target.City, err)
		}
		err = put(tracking)
		if err != nil {
			log.Printf("gov: %s-%s, put err: %s\n", target.Prefecture, target.City, err)
		}
	}

	log.Printf("gov-site-tracker finished\n")
	w.WriteHeader(200)
}

func track(target *Target) (*Tracking, error) {
	start := time.Now()
	tracking := &Tracking{
		Target:    *target,
		Timestamp: start,
	}

	req, err := http.NewRequest("GET", target.URL, nil)
	if err != nil {
		return tracking, err
	}

	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return tracking, err
	}
	defer res.Body.Close()

	code := res.StatusCode
	msec := time.Now().Sub(start).Milliseconds()

	tracking.ResCode = bigquery.NullInt64{
		Int64: int64(code),
		Valid: true,
	}
	tracking.ResMsec = bigquery.NullInt64{
		Int64: msec,
		Valid: true,
	}
	return tracking, nil
}

func put(tracking *Tracking) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, "lustrous-bus-243613")
	if err != nil {
		return err
	}
	defer client.Close()

	u := client.Dataset("gov_site").Table("tracking").Uploader()
	err = u.Put(ctx, tracking)
	if err != nil {
		return err
	}

	return nil
}

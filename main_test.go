package main

import "testing"

func TestTrack(t *testing.T) {

	// hikari
	{
		target := &Target{
			Prefecture: "yamaguchi",
			City:       "hikari",
			URL:        "https://www.city.hikari.lg.jp/",
		}
		tracking, err := track(target)
		if err != nil {
			t.Errorf("track error: %w", err)
		}
		if tracking == nil {
			t.Errorf("track nil")
		}
	}
}

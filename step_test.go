package main

import "testing"

func TestFetchEvents(t *testing.T) {
	client := Github{"username", "accesstoken"}
	events := client.fetchMergeEvents("user", "repo")

	if len(events) < 1 {
		t.Error("Expected some merge events, got none. Maybe the repo the test was run on is boring?")
	}
}

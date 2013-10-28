package cached_rss2

import (
	rss "github.com/ungerik/go-rss"
	"testing"
	"time"
)

// testFeed is a sample FetchedFeed defined here for convenience
var testFeed FetchedFeed = FetchedFeed{
	time.Now(),
	"http://test/feed",
	&rss.Channel{
		Title:       "test feed",
		Link:        "http://test/feed",
		Description: "test feed for testing",
		Item: []rss.Item{
			rss.Item{
				Title: "test one",
				Link:  "http://test/1",
				GUID:  "test-1",
			},
		},
	},
}

func TestHashURL(t *testing.T) {
	path := "http://www.google.com"
	knownCorrectHash := "ed646a3334ca891fd3467db131372140"

	hashed := hashURL(path)
	if hashed != knownCorrectHash {
		t.Error("expected: ", knownCorrectHash, "actual: ", hashed)
	}
}

func TestWriteFeedToCache(t *testing.T) {
	// call writeFeedToCache on temp dir with testFeed
	// verify that an output file matching the hash of the URL is there
	// verify that it's readable by encodings/gob
	// verify that it decodes to testFeed

	// verify that if given a bad path, or a bad feed, it returns the
	// right error
}

func TestReadFeedFromCache(t *testing.T) {
	// verify that if given a bad path, it returns the right error
	// write some junk to a temp file named with a hash URL, and verify
	//   that it can't read the file
	// call writeFeedToCache and then readFeedFromCache and verify that
	//   it roundtrips correctly
}

func TestFetchFeed(t *testing.T) {
	// fetch a well-known feed
	// verify that the insides are as expected, the FetchDate etc is right
	// fetch a bogus feed & verify that it returns error
}

func TestCachedFeed(t *testing.T) {
	// fetch a bogus feed & verify it returns error
	// write some junk to the filesystem named with a hash URL, and verify
	//   it chokes on fetching that URL
	// fetch a well known feed and verify that it creates an output file
	// with the right name
}

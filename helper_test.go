package cached_rss2

import "testing"

func TestHashURL(t *testing.T) {
	path := "http://www.google.com"
	knownCorrectHash := "ed646a3334ca891fd3467db131372140"

	hashed := hashURL(path)
	if hashed != knownCorrectHash {
		t.Error("expected: ", knownCorrectHash, "actual: ", hashed)
	}
}

// todo test writeFeedToCache

// todo test readFeedFromCache

// todo test roundtrip ^^
//   except with a bit in between where we deliberately b0rk the file?

// todo test FetchFeed

// todo test CachedFeed

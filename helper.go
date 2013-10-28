/*
 Package cached_rss2 provides a simple RSS 2.0 reader with a cache layer,
 using `github.com/ungerik/go-rss` for the RSS reader and `encodings/gob`
 for serialization.
 */
package cached_rss2

import (
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	rss "github.com/ungerik/go-rss"
	"io"
	"os"
	"path/filepath"
	"time"
)

// FetchedFeed describes a fetched RSS feed with fetch time and URL metadata
// as well as the embedded rss.Channel
type FetchedFeed struct {
	FetchDate time.Time
	URL       string
	*rss.Channel
}

// hashURL MD5-hashes the given URL and returns a hexadecimal representation
func hashURL(url string) string {
	h := md5.New()
	io.WriteString(h, url)
	baseName := hex.EncodeToString(h.Sum(nil))
	return baseName
}

// writeFeedToCache serialize a FetchedFeed as a binary file (encodings/gob)
// to the given target directory
func writeFeedToCache(feed *FetchedFeed, cacheDir string) error {
	// Hash the feed URL to provide a unique cache file name. Add the file
	// type extension too.
	cacheFileName := hashURL(feed.URL) + ".gob"
	cachePath := filepath.Join(cacheDir, cacheFileName)
	outFile, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(feed)
	if err != nil {
		return err
	}

	return nil
}

// readFeedFromCache deserializes a binary file representing a FetchedFeed
// from the given target directory. This presumes that the file was written
// by the writeFeedToCache sister function.
func readFeedFromCache(feedURL string, cacheDir string) (*FetchedFeed, error) {
	cacheFileName := hashURL(feedURL) + ".gob"
	cachePath := filepath.Join(cacheDir, cacheFileName)
	inFile, err := os.Open(cachePath)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	decoder := gob.NewDecoder(inFile)
	var feed FetchedFeed
	err = decoder.Decode(&feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

// FetchFeed retrieves an rss.Channel from the given URL, wrapping it in a
// FetchedFeed structure to add metadata.
func FetchFeed(feedURL string) (*FetchedFeed, error) {
	c, err := rss.Read(feedURL)
	if err != nil {
		return nil, err
	}

	f := FetchedFeed{time.Now(), feedURL, c}
	return &f, nil
}

// CachedFeed retrieves a FetchedFeed from the given URL, trying the cache
// first. It also writes anything it fetches to the cache for later retrieval.
func CachedFeed(feedURL string, cacheDir string, cacheTimeout time.Duration) (*FetchedFeed, error) {
	feed, err := readFeedFromCache(feedURL, cacheDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if feed != nil {
		// Check the feed's age
		timeSince := time.Since(feed.FetchDate)
		if timeSince > cacheTimeout {
			// Cache is too old
			feed = nil
		}
	}
	if feed == nil {
		// Needs fetching
		feed, err = FetchFeed(feedURL)
		if err != nil {
			return nil, err
		}
		// Bail if we can't cache the feed -- this is a *cached* feed after all
		err = writeFeedToCache(feed, cacheDir)
		if err != nil {
			return nil, err
		}
	}
	return feed, nil
}

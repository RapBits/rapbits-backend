package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
)

var songCache *bigcache.BigCache
var songsCache *bigcache.BigCache

// InitCache will inititiate our cache
func InitCache() {

	config := bigcache.Config{
		// number of shards must be power of 2
		Shards: 1024,
		// time after which entry can be evicted
		LifeWindow: 60 * time.Minute,
		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 10,
		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,
	}

	var initErr error

	songCache, initErr = bigcache.NewBigCache(config)
	if initErr != nil {
		log.Fatal("songCache", initErr)
	}

	songsCache, initErr = bigcache.NewBigCache(config)
	if initErr != nil {
		log.Fatal("songsCache", initErr)
	}

}

// RetrieveSong will return the song with songID from the cache if it exists.
// Otherwise it will make a database request and fill in the cache with the
// result, and return it
func RetrieveSong(songID string) (*Song, error) {
	var err error
	// try to retrieve from cache first
	if entry, err := songCache.Get("songID"); err == nil {
		cachedSong := new(Song)
		if err = json.Unmarshal(entry, cachedSong); err == nil {
			return cachedSong, nil
		}
	}

	// if song isn't in cache retrieve from db and put in cache
	if song, err := GetSong(songID); err == nil {
		songBytes := new(bytes.Buffer)
		json.NewEncoder(songBytes).Encode(song)
		songCache.Set(songID, songBytes.Bytes())
		return song, nil
	}

	return nil, err
}

// RetrieveSongs will return 30 songs starting from offset from the cache if it
// exists.  Otherwise it will make a database request and fill in the cache
// with the result, and return it
func RetrieveSongs(offset int) ([]*Song, error) {

	var err error
	// try to retrieve from cache first
	if entry, err := songsCache.Get(strconv.Itoa(offset)); err == nil {
		cachedSongs := make([]*Song, 0)
		if err = json.Unmarshal(entry, cachedSongs); err == nil {
			return cachedSongs, nil
		}
	}

	if songsFromOffset, err := GetSongs(offset); err == nil {
		songsBytes := new(bytes.Buffer)
		json.NewEncoder(songsBytes).Encode(songsFromOffset)
		songsCache.Set(strconv.Itoa(offset), songsBytes.Bytes())
		return songsFromOffset, nil
	}

	return nil, err
}

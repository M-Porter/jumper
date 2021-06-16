package core

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"time"
)

var (
	cacheOpenFileOptions             = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	cacheOpenFileMode    os.FileMode = 0666

	// one week cache
	staleCacheTime = time.Hour * 168
)

type Cache struct {
	LastUpdate  time.Time
	Directories []string
}

// updates the cache if its stale
func isCacheStale(fromPath string) (bool, error) {
	c, err := readFromCache(fromPath)
	if err != nil {
		return false, err
	}
	return c.LastUpdate.Sub(time.Now().UTC()).Seconds() > staleCacheTime.Seconds(), nil
}

func writeToCache(path string, dirs []string) error {
	c := Cache{
		LastUpdate:  time.Now().UTC(),
		Directories: dirs,
	}

	b, err := encodeList(c)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, cacheOpenFileOptions, cacheOpenFileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	return err
}

func readFromCache(path string) (*Cache, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	c, err := decodeList(bf)
	if err != nil {
		// todo - delete the cache file and start over
	}
	return c, nil
}

func encodeList(c Cache) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(c); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decodeList(b []byte) (*Cache, error) {
	var c *Cache
	d := gob.NewDecoder(bytes.NewReader(b))
	if err := d.Decode(&c); err != nil {
		return nil, err
	}
	return c, nil
}

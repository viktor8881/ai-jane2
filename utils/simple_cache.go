package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type Data struct {
	Exp     time.Time   `json:"Exp"`
	Content interface{} `json:"Content"`
}

type simpleCache struct {
	dirname string
	mu      sync.RWMutex
}

func NewSimpleCache(dirname string) *simpleCache {
	return &simpleCache{dirname: dirname}
}

func (sc *simpleCache) Get(key string) (interface{}, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	fileName := sc.genFileName(key)

	_, err := os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	fContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	data := Data{}
	if err = json.Unmarshal(fContent, &data); err != nil {
		return nil, err
	}

	if time.Now().After(data.Exp) {
		return nil, errors.New("The value by key: `" + key + "` is expired.")
	}

	return data.Content, nil
}

func (sc *simpleCache) Set(key string, value interface{}, secExp int) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	fileName := sc.genFileName(key)

	data := Data{
		Exp:     time.Now().Add(time.Second * time.Duration(secExp)),
		Content: value,
	}

	rec, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(rec)
	return err
}

func (sc *simpleCache) Remove(key string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	return os.Remove(sc.genFileName(key))
}

func (sc *simpleCache) genFileName(key string) string {
	hash := sha256.Sum256([]byte(key))
	return sc.dirname + hex.EncodeToString(hash[:]) + ".cache"
}

func (sc *simpleCache) Clear() error {
	files, err := os.ReadDir(sc.dirname)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			if file.Name() == ".gitignore" {
				continue
			}
			filePath := sc.dirname + "/" + file.Name()
			err := os.Remove(filePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

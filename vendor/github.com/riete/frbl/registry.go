package frbl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

const REGISTRY = "/tmp/registry-frbl.json"

var mutex sync.Mutex

func offsetUpdate(path string, offset int64) error {
	mutex.Lock()
	defer mutex.Unlock()

	f, err := os.OpenFile(REGISTRY, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	record := make(map[string]int64)
	if string(data) != "" {
		if err := json.Unmarshal(data, &record); err != nil {
			return err
		}
	}

	record[path] = offset
	newRecord, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if err := f.Truncate(0); err != nil {
		return err
	}
	_, err = f.WriteAt(newRecord, 0)
	return err
}

func offsetGet(path string) int64 {
	f, err := os.Open(REGISTRY)
	if err != nil {
		return 0
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return 0
	}

	record := make(map[string]int64)
	if err := json.Unmarshal(data, &record); err != nil {
		return 0
	}
	return record[path]
}

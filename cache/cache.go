package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/initsuj/gomc/mcauth"
)

const (
	filename = "accounts.json"
	timeout = time.Duration(100) * time.Millisecond
)

var (
	accounts = make(map[string]mcauth.Account)
	state = &sync.RWMutex{}

	watcher *fsnotify.Watcher
	timer   *time.Timer
	done = make(chan bool)
)

func Init() error {
	var err error

	if err = load(); err != nil {
		return nil
	}

	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err = watcher.Add(filename); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op & fsnotify.Write == fsnotify.Write {
					changed()
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			case <-done:
				break
			}
		}
	}()

	return nil
}

func Store(a mcauth.Account) error {
	state.RLock()
	accounts[a.Login] = a
	state.RUnlock()

	if err := save(); err != nil{
		return err
	}

	if watcher == nil{
		return load()
	}
	return nil
}

func Find(s string) (a mcauth.Account) {
	state.RLock()
	a = accounts[s]
	state.RUnlock()

	return a
}

func Close() {
	if watcher != nil{
		done <- true
		watcher.Close()
	}
}

func load() error {
	c, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	state.Lock()
	if err := json.Unmarshal(c, &accounts); err != nil {
		return err
	}
	state.Unlock()

	return nil
}

func save() error {
	state.RLock()
	contents, err := json.MarshalIndent(accounts, "", "    ")
	if err != nil {
		return err
	}
	state.RUnlock()
	runtime.Gosched()

	return write(contents)

}

func write(b []byte) error {
	err := ioutil.WriteFile(filename, b, 7777)
	if err != nil {
		return err
	}

	return nil
}

func changed() {
	if timer == nil {
		go func() {
			timer = time.NewTimer(timeout)
			<-timer.C
			timer = nil
			if err := load(); err != nil {
				log.Println("Could not update cache: ", err)
			}
		}()
	} else {
		timer.Reset(timeout)
	}
}


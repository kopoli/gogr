package gogr

import "sync"

// Options is an interface to get and set string-like options for components
type Options interface {
	Set(key string, value string)
	Get(key string, fallback string) string
}

func GetOptions() Options {
	return &options
}

var options optionMap

// optionMap implements the Options interface with a map
type optionMap struct {
	values map[string]string
	mutex   sync.Mutex
}

func (o *optionMap) initialize() {
	if o.values == nil {
		o.values = make(map[string]string)
	}
}

// Set sets the option key with value
func (o *optionMap) Set(key string, value string) {
	o.initialize()
	o.mutex.Lock()
	o.values[key] = value
	o.mutex.Unlock()
}

// Get gets the value of a key or if not available, returns the fallback
func (o *optionMap) Get(key string, fallback string) string {
	var ret string
	var ok bool
	o.mutex.Lock()
	if ret, ok = o.values[key]; !ok {
		o.mutex.Unlock()
		return fallback
	}
	o.mutex.Unlock()
	return ret
}

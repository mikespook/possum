package possum

import "sync"

// Application keeps data globally.
type Application struct {
	sync.RWMutex
	data map[string]interface{}
}

// NewApplication returns an Application instance.
func NewApplication() *Application {
	return &Application{
		data: make(map[string]interface{}),
	}
}

// Get returns value for a specific key with interface{} type.
func (app *Application) Get(key string) interface{} {
	defer app.RUnlock()
	app.RLock()
	return app.data[key]
}

// Set saves value to a specific key.
func (app *Application) Set(key string, value interface{}) {
	defer app.Unlock()
	app.Lock()
	app.data[key] = value
}

// Delete return value for a specific key,
// and remove it from the global map.
func (app *Application) Delete(key string) (value interface{}) {
	defer app.Unlock()
	app.Lock()
	value = app.data[key]
	delete(app.data, key)
	return value
}

// Singleton
var application = NewApplication()

func Get(key string) interface{} {
	return application.Get(key)
}

func Set(key string, value interface{}) {
	application.Set(key, value)
}

func Delete(key string) (value interface{}) {
	return application.Delete(key)
}

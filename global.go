package possum

import "sync"

type Application struct {
	sync.RWMutex
	data map[string]interface{}
}

func NewApplication() *Application {
	return &Application{
		data: make(map[string]interface{}),
	}
}

func (app *Application) Get(key string) interface{} {
	defer app.RUnlock()
	app.RLock()
	return app.data[key]
}

func (app *Application) Set(key string, value interface{}) {
	defer app.Unlock()
	app.Lock()
	app.data[key] = value
}

func (app *Application) Delete(key string) (value interface{}) {
	defer app.Unlock()
	app.Lock()
	value = app.data[key]
	delete(app.data, key)
	return value
}

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

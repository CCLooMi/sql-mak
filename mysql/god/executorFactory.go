package god

import (
	"reflect"
)

type SQLExecutorFactory struct {
}

var executorProviders map[string]reflect.Value = make(map[string]reflect.Value)

func RegisterExecutorProvider(name string, providerValue reflect.Value) {
	executorProviders[name] = providerValue
}

func GetExecutor(name string) *reflect.Value {
	providerValue, ok := executorProviders[name]
	if ok {
		return &providerValue
	}
	return nil
}

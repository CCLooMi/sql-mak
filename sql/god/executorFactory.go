package god

import (
	"reflect"
)

type SQLExecutorFactory struct {
}

var executorProviders map[reflect.Type]reflect.Type = make(map[reflect.Type]reflect.Type)

func RegisterExecutorProvider(executorType, providerType reflect.Type) {
	executorProviders[executorType] = providerType
}

func GetExecutor(god *SQLGod, executorType reflect.Type) interface{} {
	providerType, ok := executorProviders[executorType]
	if ok {
		provider := reflect.New(providerType)
		provider.FieldByName("god").Set(reflect.ValueOf(god))
		return provider
	}
	return nil
}

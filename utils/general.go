package utils

// Must panic if err != nil
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

type ContextProvider interface {
	Get(key string) (interface{}, bool)
	MustGet(key string) interface{}
}

package types

// Map is a key/value map.
type Map map[interface{}]interface{}

// Pair is a key/value pair.
type Pair struct {
	Key   interface{}
	Value interface{}
}

// PairFunc is a function that returns a Pair instance.
type PairFunc func(key interface{}, value interface{}) Pair

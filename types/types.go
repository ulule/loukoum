package types

// Map is a map of empty interfaces (typically used for UPDATE).
type Map map[interface{}]interface{}

// Pair is a pair of key and value.
type Pair struct {
	Key   interface{}
	Value interface{}
}

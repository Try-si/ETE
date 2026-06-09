package ETEHelper

func GetAllKeys[T any](m map[[2]int]T) [][2]int {
	keys := make([][2]int, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

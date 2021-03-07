package vivian

func makeStringSlice(generics interface{}) []string {
	genericsSlice := generics.([]interface{})
	result := make([]string, len(genericsSlice))
	for i, generic := range genericsSlice {
		result[i] = generic.(string)
	}
	return result
}

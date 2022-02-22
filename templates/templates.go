package templates

// Generates a slice of numbers
func Iterate(count int) []int {
	var i int
	var Items []int
	for i = 1; i < count+1; i++ {
		Items = append(Items, i)
	}
	return Items
}

// Checks if the value is not nil
func NotNil(v interface{}) bool {
	return v != nil
}

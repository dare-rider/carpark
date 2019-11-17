package utils

// HandleError panics an error
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

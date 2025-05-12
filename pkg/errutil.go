package pkg

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Require[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

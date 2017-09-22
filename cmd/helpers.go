package cmd

import "fmt"

func log(message ...string) {
	s := make([]interface{}, len(message)-1)
	for i := 1; i < len(message); i++ {
		s[i-1] = message[i]
	}
	fmt.Println(fmt.Sprintf(message[0], s...))
}

// abort: aborts this program on any error
func abort(err error) {
	if err != nil {
		panic(err)
	}
}

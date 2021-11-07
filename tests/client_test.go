package tests

import "io/ioutil"

func getToken() string {
	b, _ := ioutil.ReadFile("token.ini")
	if len(b) == 0 {
		return ""
	}
	return string(b)
}

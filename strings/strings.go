package strings

import (
	"math/rand"
	"time"
)

//ArrayContains returns boolean if the string array contains the given string parameter
func ArrayContains(arr []string, s string) bool {
	if arr == nil {
		return false
	}
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

//FindFirst returns the first value in the array that satisfies the given function
func FindFirst(arr []string, finder func(string) bool) (found string, index int) {
	if arr == nil {
		return "", -1
	}
	for i, a := range arr {
		if finder(a) {
			return a, i
		}
	}
	return "", -1
}

//ArraysIntersect returns intersection of two arrays
func ArraysIntersect(arr1, arr2 []string) []string {
	var intersect = make([]string, 0)
	for _, el := range arr1 {
		if ArrayContains(arr2, el) {
			intersect = append(intersect, el)
		}
	}
	return intersect
}

//GetKeys returns the list of keys of the given map
func GetKeys(m map[string]string) []string {
	var retArray = make([]string, 0)
	for k := range m {
		retArray = append(retArray, k)
	}
	return retArray
}

func RandOfSize(n int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}


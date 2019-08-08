package utils

import (
	"log"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func Intersect(list1 []string, list2 []string) []string {
	hash := make(map[string]bool)
	ret := []string{}

	for i := 0; i < len(list1); i++ {
		hash[list1[i]] = true
	}

	for i := 0; i < len(list2); i++ {
		if _, found := hash[list2[i]]; found {
			ret = append(ret, list2[i])
		}
	}
	return ret
}

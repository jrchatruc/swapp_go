package utils

import (
	"log"
	"time"
)

func Intersect(list1 []string, list2 []string) []string {
	ret := []string{}

	for _, s := range list1 {
		if contains(list2, s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func contains(a []string, s string) bool {
	for _, i := range a {
		if i == s {
			return true
		}
	}
	return false
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

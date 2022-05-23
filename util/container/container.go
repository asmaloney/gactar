package container

import (
	"sort"
)

// Contains returns whether the "value" exists in the "list".
// Case-sensitive.
func Contains(value string, list []string) bool {
	if value == "" {
		return false
	}

	for _, v := range list {
		if v == value {
			return true
		}
	}

	return false
}

// GetIndex1 returns the index (indexed from 1) of "value" or -1 if not found.
// Case-sensitive.
func GetIndex1(value string, list []string) int {
	for i, v := range list {
		if v == value {
			return i + 1
		}
	}

	return -1
}

// UniqueAndSorted returns the list "s" sorted with duplicates removed.
// Case-sensitive.
func UniqueAndSorted(s []string) (list []string) {
	unique := make(map[string]bool, len(s))
	list = make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				list = append(list, elem)
				unique[elem] = true
			}
		}
	}

	sort.Strings(list)

	return
}

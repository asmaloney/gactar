// Package container implements some routines for working with slices.
package container

import (
	"sort"
)

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

func AppendUnique[T comparable](slice []T, item T) []T {
	for _, i := range slice {
		if i == item {
			return slice
		}
	}

	return append(slice, item)
}

func FindAndDelete[T comparable](s []T, item T) []T {
	index := 0
	for _, i := range s {
		if i != item {
			s[index] = i
			index++
		}
	}
	return s[:index]
}

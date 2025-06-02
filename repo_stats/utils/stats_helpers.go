package utils

import (
	"fmt"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
)

func topnMapStrInt(x map[string]int, n int) map[string]int {
	keys := make([]string, 0, len(x))
	for key := range x {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return x[keys[i]] > x[keys[j]]
	})

	if n > len(keys) {
		n = len(keys)
	}

	result := make(map[string]int)
	for i := 0; i < n && i < len(keys); i++ {
		key := keys[i]
		result[key] = x[key]
	}
	return result
}

func (x *Stats) filterFiles(fileMap map[string]int) map[string]int {
	result := make(map[string]int)
	for file, _ := range fileMap {
		if slices.Contains(x.ignoreFiles, filepath.Base(file)) {
			continue
		}

		if slices.Contains(x.ignoreExtensions, filepath.Ext(file)) {
			continue
		}

		result[file] = fileMap[file]
	}
	return result
}

func printTop(items map[string]int) {
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return items[keys[i]] > items[keys[j]] // Descending order
	})

	// Print in sorted order
	for index, key := range keys {
		OutputFrom([]string{strconv.Itoa(index + 1), key, strconv.Itoa(items[key])},
			[]Color{Subtle, Highlight, Subtle})
	}
	fmt.Println()
}

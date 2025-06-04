package utils

import (
	"fmt"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// topnMapStrInt
// Gets the top n items from a map, ranked by their integer value
//
// Parameters:
//   - x: map of strings to integers
//
// Returns map of only the top 5 items ranked by value
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

func (x *Stats) isInsideDirectory(filePath string) (bool, error) {
	cleanFile := filepath.Clean(filePath)

	for _, dirPath := range x.ignoreDirs {
		cleanDir := filepath.Clean(dirPath)

		// Add separator to ensure we match full directory names
		if !strings.HasSuffix(cleanDir, string(filepath.Separator)) {
			cleanDir += string(filepath.Separator)
		}

		// Check if file path starts with directory path
		if strings.HasPrefix(cleanFile+string(filepath.Separator), cleanDir) {
			return true, nil
		}
	}
	return false, nil
}

// filterFiles
// Filters files based on filepath and file extension filtering rules
//
// Parameters:
//   - fileMap: map of file paths to some arbitrary integer value (size, changes, etc)
//
// Returns: filtered fileMap
func (x *Stats) filterFiles(fileMap map[string]int) map[string]int {
	result := make(map[string]int)
	for file, _ := range fileMap {
		if slices.Contains(x.ignoreFiles, filepath.Base(file)) {
			continue
		}

		if slices.Contains(x.ignoreExtensions, filepath.Ext(file)) {
			continue
		}

		ignoreDir, err := x.isInsideDirectory(file)
		if err != nil {
			continue
		}
		if ignoreDir {
			continue
		}

		result[file] = fileMap[file]
	}
	return result
}

// printTop
// Prints the items in order of an integer value, prefixed by a number
//
// Parameters:
//   - items: map of items to output with integer value denoting relative value
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

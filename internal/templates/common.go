package templates

import "strings"


func ExtractReviewers(records [][]string) []string {
	var allRevIDs []string
	for _, rec := range records {
		if len(rec) < 1 {
			continue
		}
		refID := ExtractReviewerID(rec[0])
		if SearchSlice(allRevIDs, refID) == -1 {
			allRevIDs = append(allRevIDs, refID)
		}
	}
	return allRevIDs
}

func ExtractReviewerID(fullID string) string {
	return strings.Split(strings.Split(strings.Split(strings.Split(fullID, ".")[0], "-")[0], ":")[0], " ")[0]
}

// SearchSlice is a generic linear search function that works for any slice type
func SearchSlice[T comparable](slice []T, element T) int {
	for i, v := range slice {
		if v == element {
			return i
		}
	}
	return -1
}
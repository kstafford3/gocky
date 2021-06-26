package gocky

import (
	"testing"
)

func TestContains(t *testing.T) {
	type test struct {
		searchSpace  []string
		members      []string
		notMembers   []string
	}

	testCases := []test{
		{
			searchSpace: []string{ "one", "two", "three" },
			members: []string{ "three", "two", "one" },
			notMembers: []string{ "thre", "four", "five" },
		},
	}

	for _, testCase := range testCases {
		for _, member := range testCase.members {
			if !contains(testCase.searchSpace, member) {
				t.Errorf("List [%v] should contain \"%s\"", testCase.searchSpace, member)
			}
		}

		for _, notMember := range testCase.notMembers {
			if contains(testCase.searchSpace, notMember) {
				t.Errorf("List [%v] should not contain \"%s\"", testCase.searchSpace, notMember)
			}
		}
	}
}
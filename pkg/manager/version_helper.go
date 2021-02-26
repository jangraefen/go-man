package manager

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
)

func toVersionName(versionNumber *version.Version) string {
	segments := versionNumber.Segments()
	segmentsLen := len(segments)

	// nolint:prealloc
	var (
		segmentNames      []string
		encounteredNumber bool
	)

	for index := range segments {
		segment := segments[segmentsLen-index-1]
		if !encounteredNumber {
			if segment == 0 {
				continue
			} else {
				encounteredNumber = true
			}
		}

		segmentNames = append([]string{strconv.Itoa(segment)}, segmentNames...)
	}

	return strings.Join(segmentNames, ".")
}

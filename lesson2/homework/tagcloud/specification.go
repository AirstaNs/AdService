package tagcloud

import (
	"sort"
)

var tagCloud TagCloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	StatMap map[string]int
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() TagCloud {
	tagCloud.StatMap = make(map[string]int)
	return TagCloud{StatMap: tagCloud.StatMap}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (TagCloud) AddTag(tag string) {
	if elem, ok := tagCloud.StatMap[tag]; ok {
		elem++
		tagCloud.StatMap[tag] = elem
	} else {
		tagCloud.StatMap[tag] = 1
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (TagCloud) TopN(n int) []TagStat {
	statMap := tagCloud.StatMap
	stats := make([]TagStat, 0, len(statMap))
	for k, v := range statMap {
		stats = append(stats, TagStat{Tag: k, OccurrenceCount: v})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].OccurrenceCount > stats[j].OccurrenceCount
	})
	if n > len(tagCloud.StatMap) {
		return stats
	} else {
		return stats[:n]
	}
}

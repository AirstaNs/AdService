package tagcloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	StatMap  map[string]int
	StatList []TagStat
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() *TagCloud {
	return &TagCloud{
		StatMap:  make(map[string]int),
		StatList: make([]TagStat, 0),
	}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (tagCloud *TagCloud) AddTag(tag string) {
	// Получаем индекс по которому хранится tag в TagList из tagMap O(1)
	//	Либо добавляет в конец tag с частотой 1
	if currIndex, ok := tagCloud.StatMap[tag]; ok {
		newStat := tagCloud.StatList[currIndex]
		newStat.OccurrenceCount++

		tagCloud.StatList[currIndex] = newStat
		statList := tagCloud.StatList

		// Ищем с помощью бинарного поиска первое вхождение большего числа и свапаем от currIndex до indexGreat,
		//либо пока  statList[i].OccurrenceCount < statList[more].OccurrenceCount (next элемент не станет больше currIndex),
		// если нету большего элемента от currIndex до indexGreat не свапаем и выходим из метода. O(logN)
		indexGreat := binarySearchTagStat(statList, 0, currIndex, statList[currIndex].OccurrenceCount)
		if indexGreat == -1 || indexGreat == currIndex {
			return
		}
		for i := currIndex; i > indexGreat; i-- {
			more := i - 1
			if statList[i].OccurrenceCount < statList[more].OccurrenceCount {
				break
			}
			if statList[i].OccurrenceCount > statList[more].OccurrenceCount {
				swap(i, more, tagCloud.StatList)
				tagCloud.StatMap[tag] = more
				tagCloud.StatMap[statList[i].Tag] = i
			}
		}
	} else {
		tagCloud.StatList = append(tagCloud.StatList, TagStat{Tag: tag, OccurrenceCount: 1})
		tagCloud.StatMap[tag] = len(tagCloud.StatList) - 1
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (tagCloud *TagCloud) TopN(n int) []TagStat {
	stats := tagCloud.StatList
	if n > len(tagCloud.StatList) {
		return stats
	} else {
		return stats[:n]
	}
}

func swap(i, j int, list []TagStat) {
	list[i], list[j] = list[j], list[i]
}

func binarySearchTagStat(arr []TagStat, low, high, currElement int) int {
	for low <= high {
		mid := (low + high) / 2
		if arr[mid].OccurrenceCount > currElement {
			return mid
		} else if arr[mid].OccurrenceCount < currElement {
			return mid
		} else {
			low = mid + 1
		}
	}
	return -1
}

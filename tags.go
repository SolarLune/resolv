package resolv

import "strconv"

// Tags represents one or more bitwise tags contained within a single uint64.
// You can use a tag to easily identify a type of object (e.g. player, solid, ramp, platform, etc).
// The maximum number of tags one can define is 64 (to match the uint size).
type Tags uint64

// Set sets the tag value indicated to the Tags object.
// Note that you can combine tags using the bitwise operator `|` (e.g. `TagSolidWall | TagPlatform`).
func (t *Tags) Set(tagValue Tags) {
	(*t) = (*t) | tagValue
}

// Unset clears the tag value indicated in the Tags object.
// Note that you can combine tags using the bitwise operator `|` (e.g. `TagSolidWall | TagPlatform`).
func (t *Tags) Unset(tagValue Tags) {
	(*t) = (*t) ^ tagValue
}

// Clear clears the Tags object.
func (t *Tags) Clear() {
	(*t) = 0
}

// Has returns if the Tags object has the tags indicated by tagValue set.
// Note that you can combine tags using the bitwise operator `|` (e.g. `TagSolidWall | TagPlatform`).
func (t Tags) Has(tagValue Tags) bool {
	return t&tagValue > 0
}

// IsEmpty returns if the Tags object has no tags set.
func (t Tags) IsEmpty() bool {
	return t == 0
}

// String prints out the tags set in the Tags object as a human-readable string.
func (t Tags) String() string {
	result := "Tags : [ "

	tagIndex := 0

	for i := 0; i < 64; i++ {
		possibleTag := Tags(1 << i)
		if t.Has(possibleTag) {
			if tagIndex > 0 {
				result += "| "
			}

			value, ok := tagDirectory[possibleTag]

			if !ok {
				value = strconv.Itoa(int(possibleTag))
			}

			result += value + " "
			tagIndex++
		}
	}
	result += "]"

	return result
}

var tagDirectory = map[Tags]string{}
var currentTagIndex = Tags(1)

// Creates a new tag with the given human-readable name associated with it.
// You can also create tags using bitwise representation directly (`const myTag = resolv.Tags << 1`).
// Be sure to use either method, rather than both; if you do use both, NewTag()'s internal tag index would be mismatched.
// The maximum number of tags one can define is 64.
func NewTag(tagName string) Tags {
	t := Tags(currentTagIndex)
	tagDirectory[currentTagIndex] = tagName
	currentTagIndex = currentTagIndex << 1
	return t
}

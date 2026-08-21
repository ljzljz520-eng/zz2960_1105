package archive

import (
	"sort"
	"strings"
)

type Index struct {
	Files   []string
	Digests map[string]string
}

func NewIndex() Index { return Index{Files: make([]string, 0), Digests: make(map[string]string)} }
func (index *Index) Add(file, digest string) {
	if index.Digests == nil {
		index.Digests = make(map[string]string)
	}
	if _, exists := index.Digests[file]; exists {
		return
	}
	index.Files = append(index.Files, file)
	index.Digests[file] = digest
	sort.Strings(index.Files)
}
func (index Index) Has(file string) bool { _, ok := index.Digests[file]; return ok }
func (index Index) Search(term string) []string {
	needle := strings.ToLower(strings.TrimSpace(term))
	result := make([]string, 0)
	for _, file := range index.Files {
		if needle == "" || strings.Contains(strings.ToLower(file), needle) {
			result = append(result, file)
		}
	}
	return result
}
func (index Index) Remove(file string) Index {
	result := NewIndex()
	for _, existing := range index.Files {
		if existing != file {
			result.Add(existing, index.Digests[existing])
		}
	}
	return result
}

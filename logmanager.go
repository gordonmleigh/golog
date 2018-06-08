package golog

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

// NameFilter represents a filter for a log, e.g. DEBUG:*
type NameFilter struct {
	Parts    []string
	Wildcard bool
}

// NameFilter has the following interface:
var _ interface {
	Compare(other NameFilter) int
	Match(name string) bool
} = new(NameFilter)

// ParseNameFilter parses a log name filter.
func ParseNameFilter(str string) (NameFilter, error) {
	f := NameFilter{
		Parts: strings.Split(str, ":"),
	}

	for i, p := range f.Parts {
		if p == "*" && i < len(f.Parts)-1 {
			return NameFilter{}, errors.New("wildcard can only appear in last position")
		}
	}

	if f.Parts[len(f.Parts)-1] == "*" {
		f.Parts = f.Parts[:len(f.Parts)-1]
		f.Wildcard = true
	}

	return f, nil
}

// MustParseNameFilter parses a name filter and panics on error.
func MustParseNameFilter(str string) NameFilter {
	f, err := ParseNameFilter(str)
	if err != nil {
		panic(err.Error())
	}
	return f
}

// Compare compares the current instance to the other one.
func (f NameFilter) Compare(other NameFilter) int {
	if len(f.Parts) != len(other.Parts) {
		// most specific wins
		return len(f.Parts) - len(other.Parts)
	}
	if f.Wildcard != other.Wildcard {
		if f.Wildcard {
			return -1
		}
		return 1
	}
	// fallback to lexicographic by part
	for i, p := range f.Parts {
		cmp := strings.Compare(p, other.Parts[i])
		if cmp != 0 {
			return cmp
		}
	}
	return 0
}

// Match returns true if the given name matches the pattern.
func (f NameFilter) Match(name string) bool {
	nameParts := strings.Split(name, ":")

	for i := 0; i < len(f.Parts) && i < len(nameParts); i++ {
		if f.Parts[i] == "*" {
			return true
		}
		if f.Parts[i] != nameParts[i] {
			return false
		}
	}

	return true
}

type writerPattern struct {
	filter NameFilter
	writer LogWriterFunc
}

// LogManager manages instances of loggers.
type LogManager struct {
	synchro  sync.Mutex
	logs     []*Logger
	patterns []writerPattern
}

// NewLogger registers a new logger.
func (m *LogManager) NewLogger(name string) *Logger {
	logger := &Logger{
		Name: name,
	}

	m.synchro.Lock()
	defer m.synchro.Unlock()

	m.logs = append(m.logs, logger)
	logger.SetWriter(m.getWriterNoLock(name))
	return logger
}

// SetWriter sets the writer for all logs matching the pattern.
func (m *LogManager) SetWriter(pattern NameFilter, writer LogWriterFunc) {
	m.synchro.Lock()
	defer m.synchro.Unlock()

	for _, logger := range m.logs {
		if pattern.Match(logger.Name) {
			logger.writer.Store(writer)
		}
	}

	m.addPatternNoLock(pattern, writer)
}

func (m *LogManager) getWriterNoLock(name string) LogWriterFunc {
	for _, p := range m.patterns {
		if p.filter.Match(name) {
			return p.writer
		}
	}
	return nil
}

func (m *LogManager) addPatternNoLock(pattern NameFilter, writer LogWriterFunc) {
	p := writerPattern{
		filter: pattern,
		writer: writer,
	}
	i := sort.Search(len(m.patterns), func(i int) bool {
		return m.patterns[i].filter.Compare(p.filter) > 0
	})

	m.patterns = append(m.patterns, p)

	if i < len(m.patterns)-1 {
		copy(m.patterns[i+1:], m.patterns[i:])
		m.patterns[i] = p
	}
}

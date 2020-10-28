package executor_test

import (
	"Nani/internal/app/executor"
	"Nani/internal/app/inhuman"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortKeywords_ShouldReturnSortedKeywords_NotError(t *testing.T) {
	m := inhuman.Keywords{
		"key": 3,
		"key1": 1,
		"key2": 10,
	}
	keys := executor.SortKeywords(m)
	assert.Equal(t, []string{"key1", "key", "key2"}, keys)
}

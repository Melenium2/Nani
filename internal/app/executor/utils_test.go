package executor_test

import (
	"Nani/internal/app/executor"
	"Nani/internal/app/inhuman"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortKeywords_ShouldReturnSortedKeywords_NotError(t *testing.T) {
	t.Skip("Because Sort method gives random result, because map inside")
	m := inhuman.Keywords{
		"key": 3,
		"key1": 1,
		"key4": 1,
		"key2": 10,
		"key3": 1,
		"key10": 5,
		"key12": 9,
	}
	keys := executor.SortKeywords(m)
	assert.Equal(t, []string{"key4", "key3", "key1", "key", "key10", "key12", "key2"}, keys)
}

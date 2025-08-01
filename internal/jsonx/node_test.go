package jsonx

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNode_children(t *testing.T) {
	n, err := Parse([]byte(`{"a": 1, "b": {"f": 2}, "c": [3, 4]}`))
	require.NoError(t, err)

	paths, _ := n.Children()
	assert.Equal(t, []string{"a", "b", "c"}, paths)
}

func TestNode_expandRecursively(t *testing.T) {
	n, err := Parse([]byte(`{"a": {"b": {"c": 1}}}`))
	require.NoError(t, err)

	n.CollapseRecursively()
	n.ExpandRecursively(0, 3)
	assert.Equal(t, `"c"`, n.Next.Next.Next.Key)
}

func TestNode_Paths(t *testing.T) {
	n, err := Parse([]byte(`{"a": 1, "b": {"f": 2}, "c": [3, {"d": 4}]}`))
	require.NoError(t, err)

	paths := make([]string, 0, 10)
	nodes := make([]*Node, 0, 10)
	n.Paths(&paths, &nodes)
	assert.Equal(t, []string{
		".a",
		".b",
		".c",
		".b.f",
		".c[0]",
		".c[1]",
		".c[1].d",
	}, paths)
}

func TestNode_Paths_Collapsed(t *testing.T) {
	n, err := Parse([]byte(`{"a": 1, "b": {"f": 2}, "c": [3, {"d": 4}]}`))
	require.NoError(t, err)
	n.CollapseRecursively()

	paths := make([]string, 0, 10)
	nodes := make([]*Node, 0, 10)
	n.Paths(&paths, &nodes)
	assert.Equal(t, []string{
		".a",
		".b",
		".c",
		".b.f",
		".c[0]",
		".c[1]",
		".c[1].d",
	}, paths)
}

func TestNode_ForEach(t *testing.T) {
	n, err := Parse([]byte(`{"a": 1, "b": 2, "c": 3}`))
	require.NoError(t, err)

	var keys []string
	n.ForEach(func(node *Node) {
		if k, err := strconv.Unquote(node.Key); err == nil {
			keys = append(keys, k)
		}
	})
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestNode_ForEach_Empty(t *testing.T) {
	n, err := Parse([]byte(`{}`))
	require.NoError(t, err)

	called := false
	n.ForEach(func(node *Node) {
		called = true
	})
	assert.False(t, called)
}

func TestNode_ForEach_SkipsNested(t *testing.T) {
	n, err := Parse([]byte(`{"a": {"b": 1}, "c": [2, {"d": 3}]}`))
	require.NoError(t, err)

	var keys []string
	n.ForEach(func(node *Node) {
		if k, err := strconv.Unquote(node.Key); err == nil {
			keys = append(keys, k)
		}
	})
	assert.Equal(t, []string{"a", "c"}, keys)
}

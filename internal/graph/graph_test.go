package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// 1 -> 2 -> 4
// 1 -> 2 -> 5
// 1 -> 2 -> 7 -> 6
// 1 -> 3 -> 8 -> 9 -> 6
// 7 -> 9
func TestDirected_Score(t *testing.T) {

	directed := NewDirected()
	directed.AddEdge("1", true, "2", false)
	directed.AddEdge("2", false, "4", false)
	directed.AddEdge("2", false, "5", false)
	directed.AddEdge("2", false, "7", false)
	directed.AddEdge("7", false, "6", false)
	directed.AddEdge("1", true, "3", false)
	directed.AddEdge("3", false, "8", false)
	directed.AddEdge("8", false, "9", false)
	directed.AddEdge("9", false, "6", false)
	directed.AddEdge("7", false, "9", false)

	directed.Score()

	require.Equal(t, FlagScore, directed.Get("2").Score)
	require.Equal(t, FlagScore, directed.Get("3").Score)
	require.Equal(t, 1, directed.Get("4").Score)
	require.Equal(t, 1, directed.Get("5").Score)
	require.Equal(t, 2, directed.Get("6").Score)
	require.Equal(t, 1, directed.Get("7").Score)
	require.Equal(t, 1, directed.Get("8").Score)
	require.Equal(t, 2, directed.Get("9").Score)
}

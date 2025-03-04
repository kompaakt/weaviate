//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2021 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package hnsw

import (
	"testing"

	"github.com/semi-technologies/weaviate/adapters/repos/db/vector/hnsw/distancer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHnswIndex(t *testing.T) {
	// mock out commit logger before adding data so we don't leave a disk
	// footprint. Commit logging and deserializing from a (condensed) commit log
	// is tested in a separate integration test that takes care of providing and
	// cleaning up the correct place on disk to write test files
	makeCL := MakeNoopCommitLogger
	index, err := New(Config{
		RootPath:              "doesnt-matter-as-committlogger-is-mocked-out",
		ID:                    "unittest",
		MakeCommitLoggerThunk: makeCL,
		DistanceProvider:      distancer.NewCosineProvider(),
		VectorForIDThunk:      testVectorForID,
	}, UserConfig{
		MaxConnections: 30,
		EFConstruction: 60,
	})
	require.Nil(t, err)

	for i, vec := range testVectors {
		err := index.Add(uint64(i), vec)
		require.Nil(t, err)
	}

	t.Run("searching within cluster 1", func(t *testing.T) {
		position := 0
		res, _, err := index.knnSearchByVector(testVectors[position], 3, 36, nil)
		require.Nil(t, err)
		assert.ElementsMatch(t, []uint64{0, 1, 2}, res)
	})

	t.Run("searching within cluster 2", func(t *testing.T) {
		position := 3
		res, _, err := index.knnSearchByVector(testVectors[position], 3, 36, nil)
		require.Nil(t, err)
		assert.ElementsMatch(t, []uint64{3, 4, 5}, res)
	})

	t.Run("searching within cluster 3", func(t *testing.T) {
		position := 6
		res, _, err := index.knnSearchByVector(testVectors[position], 3, 36, nil)
		require.Nil(t, err)
		assert.ElementsMatch(t, []uint64{6, 7, 8}, res)
	})

	t.Run("searching within cluster 2 with a scope larger than the cluster", func(t *testing.T) {
		position := 3
		res, _, err := index.knnSearchByVector(testVectors[position], 50, 36, nil)
		require.Nil(t, err)
		assert.Equal(t, []uint64{
			3, 5, 4, // cluster 2
			7, 8, 6, // cluster 3
			2, 1, 0, // cluster 1
		}, res)
	})
}

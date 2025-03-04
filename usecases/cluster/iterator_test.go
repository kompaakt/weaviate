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

package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeIteration(t *testing.T) {
	source := fakeNodeSource{[]string{"node1", "node2", "node3", "node4"}}
	it, err := NewNodeIterator(source, StartRandom)
	require.Nil(t, err)

	found := map[string]int{}

	for i := 0; i < 20; i++ {
		host := it.Next()
		found[host]++
	}

	// each host must be contained 5 times
	assert.Equal(t, found["node1"], 5)
	assert.Equal(t, found["node2"], 5)
	assert.Equal(t, found["node3"], 5)
	assert.Equal(t, found["node4"], 5)
}

type fakeNodeSource struct {
	hostnames []string
}

func (f fakeNodeSource) AllNames() []string {
	return f.hostnames
}

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

package traverser

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/aggregation"
	"github.com/semi-technologies/weaviate/entities/filters"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/search"
	"github.com/semi-technologies/weaviate/usecases/config"
	"github.com/semi-technologies/weaviate/usecases/schema"
	"github.com/sirupsen/logrus"
)

type locks interface {
	LockConnector() (func() error, error)
	LockSchema() (func() error, error)
}

type authorizer interface {
	Authorize(principal *models.Principal, verb, resource string) error
}

// Traverser can be used to dynamically traverse the knowledge graph
type Traverser struct {
	config         *config.WeaviateConfig
	locks          locks
	logger         logrus.FieldLogger
	authorizer     authorizer
	vectorSearcher VectorSearcher
	explorer       explorer
	schemaGetter   schema.SchemaGetter
}

type VectorSearcher interface {
	VectorSearch(ctx context.Context, vector []float32,
		offset, limit int, filters *filters.LocalFilter) ([]search.Result, error)
	Aggregate(ctx context.Context, params aggregation.Params) (*aggregation.Result, error)
}

type explorer interface {
	GetClass(ctx context.Context, params GetParams) ([]interface{}, error)
	Concepts(ctx context.Context, params ExploreParams) ([]search.Result, error)
}

// NewTraverser to traverse the knowledge graph
func NewTraverser(config *config.WeaviateConfig, locks locks,
	logger logrus.FieldLogger, authorizer authorizer,
	vectorSearcher VectorSearcher,
	explorer explorer, schemaGetter schema.SchemaGetter) *Traverser {
	return &Traverser{
		config:         config,
		locks:          locks,
		logger:         logger,
		authorizer:     authorizer,
		vectorSearcher: vectorSearcher,
		explorer:       explorer,
		schemaGetter:   schemaGetter,
	}
}

// TraverserRepo describes the dependencies of the Traverser UC to the
// connected database
type TraverserRepo interface {
	GetClass(context.Context, *GetParams) (interface{}, error)
	Aggregate(context.Context, *aggregation.Params) (interface{}, error)
}

// SearchResult is a single search result. See wrapping Search Results for the Type
type SearchResult struct {
	Name      string
	Certainty float32
}

// SearchResults is grouping of SearchResults for a SchemaSearch
type SearchResults struct {
	Type    SearchType
	Results []SearchResult
}

// Len of the result set
func (r SearchResults) Len() int {
	return len(r.Results)
}

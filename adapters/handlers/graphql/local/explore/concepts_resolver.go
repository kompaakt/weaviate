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

package explore

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/semi-technologies/weaviate/adapters/handlers/graphql/local/common_filters"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/search"
	"github.com/semi-technologies/weaviate/usecases/traverser"
)

// Resolver is a local interface that can be composed with other interfaces to
// form the overall GraphQL API main interface. All data-base connectors that
// want to support the Meta feature must implement this interface.
type Resolver interface {
	Explore(ctx context.Context, principal *models.Principal,
		params traverser.ExploreParams) ([]search.Result, error)
}

// RequestsLog is a local abstraction on the RequestsLog that needs to be
// provided to the graphQL API in order to log Local.Fetch queries.
type RequestsLog interface {
	Register(requestType string, identifier string)
}

type resources struct {
	resolver Resolver
}

func newResources(s interface{}) (*resources, error) {
	source, ok := s.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected source to be a map, but was %T", source)
	}

	resolver, ok := source["Resolver"].(Resolver)
	if !ok {
		return nil, fmt.Errorf("expected source to contain a usable Resolver, but was %#v", source)
	}

	return &resources{
		resolver: resolver,
	}, nil
}

type resolver struct {
	modulesProvider ModulesProvider
}

func newResolver(modulesProvider ModulesProvider) *resolver {
	return &resolver{modulesProvider}
}

func (r *resolver) resolve(p graphql.ResolveParams) (interface{}, error) {
	resources, err := newResources(p.Source)
	if err != nil {
		return nil, err
	}

	params := traverser.ExploreParams{}

	if param, ok := p.Args["nearVector"]; ok {
		extracted := common_filters.ExtractNearVector(param.(map[string]interface{}))
		params.NearVector = &extracted
	}

	if param, ok := p.Args["nearObject"]; ok {
		extracted := common_filters.ExtractNearObject(param.(map[string]interface{}))
		params.NearObject = &extracted
	}

	if param, ok := p.Args["offset"]; ok {
		params.Offset = param.(int)
	}

	if param, ok := p.Args["limit"]; ok {
		params.Limit = param.(int)
	}

	if r.modulesProvider != nil {
		extractedParams := r.modulesProvider.CrossClassExtractSearchParams(p.Args)
		if len(extractedParams) > 0 {
			params.ModuleParams = extractedParams
		}
	}

	return resources.resolver.Explore(p.Context,
		principalFromContext(p.Context), params)
}

func principalFromContext(ctx context.Context) *models.Principal {
	principal := ctx.Value("principal")
	if principal == nil {
		return nil
	}

	return principal.(*models.Principal)
}

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

// Package migrate provides a simple composer tool, which implements the
// Migrator interface and can take in any number of migrators which themselves
// have to implement the interface
package migrate

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/schema"
	"github.com/semi-technologies/weaviate/usecases/sharding"
)

// Migrator represents both the input and output interface of the Composer
type Migrator interface {
	AddClass(ctx context.Context, class *models.Class, shardingState *sharding.State) error
	DropClass(ctx context.Context, className string) error
	UpdateClass(ctx context.Context, className string,
		newClassName *string) error

	AddProperty(ctx context.Context, className string,
		prop *models.Property) error
	UpdateProperty(ctx context.Context, className string,
		propName string, newName *string) error
	ValidateVectorIndexConfigUpdate(ctx context.Context,
		old, updated schema.VectorIndexConfig) error
	UpdateVectorIndexConfig(ctx context.Context, className string,
		updated schema.VectorIndexConfig) error
}

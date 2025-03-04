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

package modspellcheck

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/modulecapabilities"
	"github.com/semi-technologies/weaviate/entities/moduletools"
	"github.com/semi-technologies/weaviate/entities/schema"
)

func (m *SpellCheckModule) ClassConfigDefaults() map[string]interface{} {
	return map[string]interface{}{}
}

func (m *SpellCheckModule) PropertyConfigDefaults(
	dt *schema.DataType) map[string]interface{} {
	return map[string]interface{}{}
}

func (m *SpellCheckModule) ValidateClass(ctx context.Context,
	class *models.Class, cfg moduletools.ClassConfig) error {
	return nil
}

var _ = modulecapabilities.ClassConfigurator(New())

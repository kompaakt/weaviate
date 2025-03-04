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
	"testing"

	"github.com/semi-technologies/weaviate/entities/aggregation"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/schema"
	"github.com/semi-technologies/weaviate/usecases/config"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Traverser_Aggregate(t *testing.T) {
	t.Run("with aggregation only", func(t *testing.T) {
		principal := &models.Principal{}
		logger, _ := test.NewNullLogger()
		locks := &fakeLocks{}
		authorizer := &fakeAuthorizer{}
		vectorRepo := &fakeVectorRepo{}
		explorer := &fakeExplorer{}
		schemaGetter := &fakeSchemaGetter{aggregateTestSchema}

		traverser := NewTraverser(&config.WeaviateConfig{}, locks, logger, authorizer,
			vectorRepo, explorer, schemaGetter)

		params := aggregation.Params{
			ClassName: "MyClass",
			Properties: []aggregation.ParamProperty{
				aggregation.ParamProperty{
					Name:        "label",
					Aggregators: []aggregation.Aggregator{aggregation.NewTopOccurrencesAggregator(nil)},
				},
				aggregation.ParamProperty{
					Name:        "number",
					Aggregators: []aggregation.Aggregator{aggregation.SumAggregator},
				},
				aggregation.ParamProperty{
					Name:        "int",
					Aggregators: []aggregation.Aggregator{aggregation.SumAggregator},
				},
				aggregation.ParamProperty{
					Name:        "date",
					Aggregators: []aggregation.Aggregator{aggregation.NewTopOccurrencesAggregator(nil)},
				},
			},
		}

		agg := aggregation.Result{
			Groups: []aggregation.Group{
				aggregation.Group{
					Properties: map[string]aggregation.Property{
						"label": aggregation.Property{
							TextAggregation: aggregation.Text{
								Items: []aggregation.TextOccurrence{
									aggregation.TextOccurrence{
										Value:  "Foo",
										Occurs: 200,
									},
								},
							},
							Type: aggregation.PropertyTypeText,
						},
						"date": aggregation.Property{
							TextAggregation: aggregation.Text{
								Items: []aggregation.TextOccurrence{
									aggregation.TextOccurrence{
										Value:  "Bar",
										Occurs: 100,
									},
								},
							},
							Type: aggregation.PropertyTypeText,
						},
						"number": aggregation.Property{
							Type: aggregation.PropertyTypeNumerical,
							NumericalAggregations: map[string]float64{
								"sum": 200,
							},
						},
						"int": aggregation.Property{
							Type: aggregation.PropertyTypeNumerical,
							NumericalAggregations: map[string]float64{
								"sum": 100,
							},
						},
					},
				},
			},
		}

		vectorRepo.On("Aggregate", params).Return(&agg, nil)
		res, err := traverser.Aggregate(context.Background(), principal, &params)
		require.Nil(t, err)
		assert.Equal(t, &agg, res)
	})

	t.Run("with a mix of aggregation and type inspection", func(t *testing.T) {
		principal := &models.Principal{}
		logger, _ := test.NewNullLogger()
		locks := &fakeLocks{}
		authorizer := &fakeAuthorizer{}
		vectorRepo := &fakeVectorRepo{}
		explorer := &fakeExplorer{}
		schemaGetter := &fakeSchemaGetter{aggregateTestSchema}

		traverser := NewTraverser(&config.WeaviateConfig{}, locks, logger, authorizer,
			vectorRepo, explorer, schemaGetter)

		params := aggregation.Params{
			ClassName: "MyClass",
			Properties: []aggregation.ParamProperty{
				aggregation.ParamProperty{
					Name: "label",
					Aggregators: []aggregation.Aggregator{
						aggregation.TypeAggregator,
						aggregation.NewTopOccurrencesAggregator(nil),
					},
				},
				aggregation.ParamProperty{
					Name: "number",
					Aggregators: []aggregation.Aggregator{
						aggregation.TypeAggregator,
						aggregation.SumAggregator,
					},
				},
				aggregation.ParamProperty{
					Name: "int",
					Aggregators: []aggregation.Aggregator{
						aggregation.TypeAggregator,
						aggregation.SumAggregator,
					},
				},
				aggregation.ParamProperty{
					Name: "date",
					Aggregators: []aggregation.Aggregator{
						aggregation.TypeAggregator,
						aggregation.NewTopOccurrencesAggregator(nil),
					},
				},
				aggregation.ParamProperty{
					Name:        "a ref",
					Aggregators: []aggregation.Aggregator{aggregation.TypeAggregator},
				},
			},
		}

		agg := aggregation.Result{
			Groups: []aggregation.Group{
				aggregation.Group{
					Properties: map[string]aggregation.Property{
						"label": aggregation.Property{
							TextAggregation: aggregation.Text{
								Items: []aggregation.TextOccurrence{
									aggregation.TextOccurrence{
										Value:  "Foo",
										Occurs: 200,
									},
								},
							},
							Type: aggregation.PropertyTypeText,
						},
						"date": aggregation.Property{
							TextAggregation: aggregation.Text{
								Items: []aggregation.TextOccurrence{
									aggregation.TextOccurrence{
										Value:  "Bar",
										Occurs: 100,
									},
								},
							},
							Type: aggregation.PropertyTypeText,
						},
						"number": aggregation.Property{
							Type: aggregation.PropertyTypeNumerical,
							NumericalAggregations: map[string]float64{
								"sum": 200,
							},
						},
						"int": aggregation.Property{
							Type: aggregation.PropertyTypeNumerical,
							NumericalAggregations: map[string]float64{
								"sum": 100,
							},
						},
					},
				},
			},
		}

		expectedResult := aggregation.Result{
			Groups: []aggregation.Group{
				aggregation.Group{
					Properties: map[string]aggregation.Property{
						"label": aggregation.Property{
							TextAggregation: aggregation.Text{
								Items: []aggregation.TextOccurrence{
									aggregation.TextOccurrence{
										Value:  "Foo",
										Occurs: 200,
									},
								},
							},
							Type:       aggregation.PropertyTypeText,
							SchemaType: string(schema.DataTypeString),
						},
						"date": aggregation.Property{
							TextAggregation: aggregation.Text{
								Items: []aggregation.TextOccurrence{
									aggregation.TextOccurrence{
										Value:  "Bar",
										Occurs: 100,
									},
								},
							},
							SchemaType: string(schema.DataTypeDate),
							Type:       aggregation.PropertyTypeText,
						},
						"number": aggregation.Property{
							Type:       aggregation.PropertyTypeNumerical,
							SchemaType: string(schema.DataTypeNumber),
							NumericalAggregations: map[string]float64{
								"sum": 200,
							},
						},
						"int": aggregation.Property{
							Type:       aggregation.PropertyTypeNumerical,
							SchemaType: string(schema.DataTypeInt),
							NumericalAggregations: map[string]float64{
								"sum": 100,
							},
						},
						"a ref": aggregation.Property{
							Type: aggregation.PropertyTypeReference,
							ReferenceAggregation: aggregation.Reference{
								PointingTo: []string{"AnotherClass"},
							},
							SchemaType: string(schema.DataTypeCRef),
						},
					},
				},
			},
		}

		vectorRepo.On("Aggregate", params).Return(&agg, nil)
		res, err := traverser.Aggregate(context.Background(), principal, &params)
		require.Nil(t, err)
		assert.Equal(t, &expectedResult, res)
	})
}

var aggregateTestSchema = schema.Schema{
	Objects: &models.Schema{
		Classes: []*models.Class{
			&models.Class{
				Class: "AnotherClass",
			},
			&models.Class{
				Class: "MyClass",
				Properties: []*models.Property{
					&models.Property{
						Name:     "label",
						DataType: []string{string(schema.DataTypeString)},
					},
					&models.Property{
						Name:     "number",
						DataType: []string{string(schema.DataTypeNumber)},
					},
					&models.Property{
						Name:     "int",
						DataType: []string{string(schema.DataTypeInt)},
					},
					&models.Property{
						Name:     "date",
						DataType: []string{string(schema.DataTypeDate)},
					},
					&models.Property{
						Name:     "a ref",
						DataType: []string{"AnotherClass"},
					},
				},
			},
		},
	},
}

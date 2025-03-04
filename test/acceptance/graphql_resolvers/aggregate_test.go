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

package test

import (
	"encoding/json"
	"testing"

	"github.com/semi-technologies/weaviate/test/acceptance/helper"
	"github.com/stretchr/testify/assert"
)

func aggregatesWithoutGroupingOrFilters(t *testing.T) {
	result := AssertGraphQL(t, helper.RootAuth, `
		{
			Aggregate{
				City {
					meta {
						count
					}
					isCapital {
						count
						percentageFalse
						percentageTrue
						totalFalse
						totalTrue
						type
					}
					population {
						mean
						count
						maximum
						minimum
						sum
						type
					}
					inCountry {
						pointingTo
						type
					}
					name {
						topOccurrences {
							occurs
							value
						}
						type
						count
					}
				}
			}
		}
	`)

	t.Run("meta count", func(t *testing.T) {
		meta := result.Get("Aggregate", "City").AsSlice()[0].(map[string]interface{})["meta"]
		count := meta.(map[string]interface{})["count"]
		expected := json.Number("5")
		assert.Equal(t, expected, count)
	})

	t.Run("boolean props", func(t *testing.T) {
		isCapital := result.Get("Aggregate", "City").AsSlice()[0].(map[string]interface{})["isCapital"]
		expected := map[string]interface{}{
			"count":           json.Number("5"),
			"percentageTrue":  json.Number("0.4"),
			"percentageFalse": json.Number("0.6"),
			"totalTrue":       json.Number("2"),
			"totalFalse":      json.Number("3"),
			"type":            "boolean",
		}
		assert.Equal(t, expected, isCapital)
	})

	t.Run("int/number props", func(t *testing.T) {
		isCapital := result.Get("Aggregate", "City").AsSlice()[0].(map[string]interface{})["population"]
		expected := map[string]interface{}{
			"mean":    json.Number("1294000"),
			"count":   json.Number("5"),
			"maximum": json.Number("3470000"),
			"minimum": json.Number("0"),
			"sum":     json.Number("6470000"),
			"type":    "int",
		}
		assert.Equal(t, expected, isCapital)
	})

	t.Run("ref prop", func(t *testing.T) {
		inCountry := result.Get("Aggregate", "City").AsSlice()[0].(map[string]interface{})["inCountry"]
		expected := map[string]interface{}{
			"pointingTo": []interface{}{"Country"},
			"type":       "cref",
		}
		assert.Equal(t, expected, inCountry)
	})

	t.Run("string prop", func(t *testing.T) {
		name := result.Get("Aggregate", "City").
			AsSlice()[0].(map[string]interface{})["name"].(map[string]interface{})
		typeField := name["type"]
		topOccurrences := name["topOccurrences"]

		assert.Equal(t, "string", typeField)

		expectedTopOccurrences := []interface{}{
			map[string]interface{}{
				"value":  "Amsterdam",
				"occurs": json.Number("1"),
			},
			map[string]interface{}{
				"value":  "Dusseldorf",
				"occurs": json.Number("1"),
			},
			map[string]interface{}{
				"value":  "Rotterdam",
				"occurs": json.Number("1"),
			},
			map[string]interface{}{
				"value":  "Berlin",
				"occurs": json.Number("1"),
			},
			map[string]interface{}{
				"value":  "Null Island",
				"occurs": json.Number("1"),
			},
		}
		assert.ElementsMatch(t, expectedTopOccurrences, topOccurrences)
	})
}

func localMetaWithFilters(t *testing.T) {
	result := AssertGraphQL(t, helper.RootAuth, `
		{
			Aggregate{
				City (where: {
					valueBoolean: true,
					operator: Equal,
					path: ["isCapital"]
				}){
					meta {
						count
					}
					isCapital {
						count
						percentageFalse
						percentageTrue
						totalFalse
						totalTrue
						type
					}
					population {
						mean
						count
						maximum
						minimum
						sum
						type
					}
					inCountry {
						pointingTo
						type
					}
					name {
						topOccurrences {
							occurs
							value
						}
						type
						count
					}
				}
			}
		}
	`)

	t.Run("meta count", func(t *testing.T) {
		meta := result.Get("Aggregate", "City").AsSlice()[0].(map[string]interface{})["meta"]
		count := meta.(map[string]interface{})["count"]
		expected := json.Number("2")
		assert.Equal(t, expected, count)
	})

	t.Run("boolean props", func(t *testing.T) {
		isCapital := result.Get("Aggregate", "City").AsSlice()[0].(map[string]interface{})["isCapital"]
		expected := map[string]interface{}{
			"count":           json.Number("2"),
			"percentageTrue":  json.Number("1"),
			"percentageFalse": json.Number("0"),
			"totalTrue":       json.Number("2"),
			"totalFalse":      json.Number("0"),
			"type":            "boolean",
		}
		assert.Equal(t, expected, isCapital)
	})

	t.Run("int/number props", func(t *testing.T) {
		population := result.Get("Aggregate", "City").AsSlice()[0].(map[string]interface{})["population"]
		expected := map[string]interface{}{
			"mean":    json.Number("2635000"),
			"count":   json.Number("2"),
			"maximum": json.Number("3470000"),
			"minimum": json.Number("1800000"),
			"sum":     json.Number("5270000"),
			"type":    "int",
		}
		assert.Equal(t, expected, population)
	})

	t.Run("ref prop", func(t *testing.T) {
		inCountry := result.Get("Aggregate", "City").AsSlice()[0].(map[string]interface{})["inCountry"]
		expected := map[string]interface{}{
			"pointingTo": []interface{}{"Country"},
			"type":       "cref",
		}
		assert.Equal(t, expected, inCountry)
	})

	t.Run("string prop", func(t *testing.T) {
		name := result.Get("Aggregate", "City").
			AsSlice()[0].(map[string]interface{})["name"].(map[string]interface{})
		typeField := name["type"]
		topOccurrences := name["topOccurrences"]

		assert.Equal(t, "string", typeField)

		expectedTopOccurrences := []interface{}{
			map[string]interface{}{
				"value":  "Amsterdam",
				"occurs": json.Number("1"),
			},
			map[string]interface{}{
				"value":  "Berlin",
				"occurs": json.Number("1"),
			},
		}
		assert.ElementsMatch(t, expectedTopOccurrences, topOccurrences)
	})
}

// This test prevents a regression on the fix for
// https://github.com/semi-technologies/weaviate/issues/824
func localMeta_StringPropsNotSetEverywhere(t *testing.T) {
	AssertGraphQL(t, helper.RootAuth, `
		{
			Aggregate {
				City {
					name {
						topOccurrences {
							occurs
							value
						}
					}
				}
			}
		}
	`)
}

func aggregatesArrayClassWithoutGroupingOrFilters(t *testing.T) {
	result := AssertGraphQL(t, helper.RootAuth, `
		{
			Aggregate{
				ArrayClass {
					meta {
						count
					}
					numbers {
						mean
						count
						maximum
						minimum
						sum
						type
					}
					strings {
						topOccurrences {
							occurs
							value
						}
						type
						count
					}
					booleans {
						count
						percentageFalse
						percentageTrue
						totalFalse
						totalTrue
						type
					}
				}
			}
		}
	`)

	t.Run("meta count", func(t *testing.T) {
		meta := result.Get("Aggregate", "ArrayClass").AsSlice()[0].(map[string]interface{})["meta"]
		count := meta.(map[string]interface{})["count"]
		expected := json.Number("3")
		assert.Equal(t, expected, count)
	})

	t.Run("int[]/number[] props", func(t *testing.T) {
		isCapital := result.Get("Aggregate", "ArrayClass").AsSlice()[0].(map[string]interface{})["numbers"]
		expected := map[string]interface{}{
			"mean":    json.Number("1.6666666666666667"),
			"count":   json.Number("6"),
			"maximum": json.Number("3"),
			"minimum": json.Number("1"),
			"sum":     json.Number("10"),
			"type":    "number[]",
		}
		assert.Equal(t, expected, isCapital)
	})

	t.Run("string[]/text[] prop", func(t *testing.T) {
		name := result.Get("Aggregate", "ArrayClass").
			AsSlice()[0].(map[string]interface{})["strings"].(map[string]interface{})
		typeField := name["type"]
		topOccurrences := name["topOccurrences"]

		assert.Equal(t, "string[]", typeField)

		expectedTopOccurrences := []interface{}{
			map[string]interface{}{
				"value":  "a",
				"occurs": json.Number("3"),
			},
			map[string]interface{}{
				"value":  "b",
				"occurs": json.Number("2"),
			},
			map[string]interface{}{
				"value":  "c",
				"occurs": json.Number("1"),
			},
		}
		assert.ElementsMatch(t, expectedTopOccurrences, topOccurrences)
	})

	t.Run("boolean props", func(t *testing.T) {
		isCapital := result.Get("Aggregate", "ArrayClass").AsSlice()[0].(map[string]interface{})["booleans"]
		expected := map[string]interface{}{
			"count":           json.Number("6"),
			"percentageTrue":  json.Number("0.5"),
			"percentageFalse": json.Number("0.5"),
			"totalTrue":       json.Number("3"),
			"totalFalse":      json.Number("3"),
			"type":            "boolean[]",
		}
		assert.Equal(t, expected, isCapital)
	})
}

func aggregatesArrayClassWithGrouping(t *testing.T) {
	result := AssertGraphQL(t, helper.RootAuth, `
	{
		Aggregate{
			ArrayClass(groupBy:["numbers"]){
				meta{
					count
				}
				groupedBy{
					value
				}
			}
		}
	}
	`)

	t.Run("groupedBy result", func(t *testing.T) {
		groupedByResults := result.Get("Aggregate", "ArrayClass").AsSlice()
		assert.Equal(t, 3, len(groupedByResults))

		for _, res := range groupedByResults {
			meta := res.(map[string]interface{})["meta"]
			count := meta.(map[string]interface{})["count"]
			groupedBy := res.(map[string]interface{})["groupedBy"]
			value := groupedBy.(map[string]interface{})["value"]
			valueString := value.(string)

			if valueString == "1" {
				assert.Equal(t, json.Number("3"), count)
			}
			if valueString == "2" {
				assert.Equal(t, json.Number("2"), count)
			}
			if valueString == "3" {
				assert.Equal(t, json.Number("1"), count)
			}
		}
	})
}

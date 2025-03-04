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

package ask

func (g *GraphQLArgumentsProvider) extractAskFn(source map[string]interface{}) interface{} {
	var args AskParams

	question, ok := source["question"].(string)
	if ok {
		args.Question = question
	}

	// autocorrect is an optional arg, so it could be nil
	autocorrect, ok := source["autocorrect"]
	if ok {
		args.Autocorrect = autocorrect.(bool)
	}

	// if there's text transformer present and autocorrect set to true
	// perform text transformation operation
	if args.Autocorrect && g.askTransformer != nil {
		if transformedValues, err := g.askTransformer.Transform([]string{args.Question}); err == nil && len(transformedValues) == 1 {
			args.Question = transformedValues[0]
		}
	}

	certainty, ok := source["certainty"]
	if ok {
		args.Certainty = certainty.(float64)
	}

	properties, ok := source["properties"].([]interface{})
	if ok {
		args.Properties = make([]string, len(properties))
		for i, value := range properties {
			args.Properties[i] = value.(string)
		}
	}

	return &args
}

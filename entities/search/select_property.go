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

package search

import (
	"fmt"
	"regexp"

	"github.com/semi-technologies/weaviate/entities/additional"
	"github.com/semi-technologies/weaviate/entities/schema"
)

type SelectProperty struct {
	Name string `json:"name"`

	IsPrimitive bool `json:"isPrimitive"`

	// Include the __typename in all the Refs below.
	IncludeTypeName bool `json:"includeTypeName"`

	// Not a primitive type? Then select these properties.
	Refs []SelectClass `json:"refs"`
}

type SelectClass struct {
	ClassName            string                `json:"className"`
	RefProperties        SelectProperties      `json:"refProperties"`
	AdditionalProperties additional.Properties `json:"additionalProperties"`
}

// FindSelectClass by specifying the exact class name
func (sp SelectProperty) FindSelectClass(className schema.ClassName) *SelectClass {
	for _, selectClass := range sp.Refs {
		if selectClass.ClassName == string(className) {
			return &selectClass
		}
	}

	return nil
}

// HasPeer returns true if any of the referenced classes are from the specified
// peer
func (sp SelectProperty) HasPeer(peerName string) bool {
	r := regexp.MustCompile(fmt.Sprintf("^%s__", peerName))
	for _, selectClass := range sp.Refs {
		if r.MatchString(selectClass.ClassName) {
			return true
		}
	}

	return false
}

type SelectProperties []SelectProperty

func (sp SelectProperties) HasRefs() bool {
	for _, p := range sp {
		if !p.IsPrimitive {
			return true
		}
	}

	return false
}

func (sp SelectProperties) ShouldResolve(path []string) (bool, error) {
	if len(path)%2 != 0 || len(path) == 0 {
		return false, fmt.Errorf("used incorrectly: path must have even number of segments in the form of " +
			"refProp, className, refProp, className, etc.")
	}

	// the above gives us the guarantuee that path contains at least two elements
	property := path[0]
	class := schema.ClassName(path[1])

	for _, p := range sp {
		if p.IsPrimitive {
			continue
		}

		if p.Name != property {
			continue
		}

		selectClass := p.FindSelectClass(class)
		if selectClass == nil {
			continue
		}

		if len(path) > 2 {
			// we're not done yet, this one's nested
			return selectClass.RefProperties.ShouldResolve(path[2:])
		}

		// we are done and found the path
		return true, nil
	}

	return false, nil
}

func (sp SelectProperties) FindProperty(propName string) *SelectProperty {
	for _, prop := range sp {
		if prop.Name == propName {
			return &prop
		}
	}

	return nil
}

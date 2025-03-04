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

//go:build integrationTest
// +build integrationTest

package db

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate/adapters/repos/db/vector/hnsw"
	"github.com/semi-technologies/weaviate/entities/additional"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/schema"
	"github.com/semi-technologies/weaviate/entities/schema/crossref"
	"github.com/semi-technologies/weaviate/usecases/objects"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MergingObjects(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	dirName := fmt.Sprintf("./testdata/%d", rand.Intn(10000000))
	os.MkdirAll(dirName, 0o777)
	defer func() {
		err := os.RemoveAll(dirName)
		fmt.Println(err)
	}()

	logger := logrus.New()
	schemaGetter := &fakeSchemaGetter{shardState: singleShardState()}
	repo := New(logger, Config{RootPath: dirName}, &fakeRemoteClient{},
		&fakeNodeResolver{})
	repo.SetSchemaGetter(schemaGetter)
	err := repo.WaitForStartup(testCtx())
	require.Nil(t, err)
	migrator := NewMigrator(repo, logger)

	schema := schema.Schema{
		Objects: &models.Schema{
			Classes: []*models.Class{
				&models.Class{
					Class:               "MergeTestTarget",
					VectorIndexConfig:   hnsw.NewDefaultUserConfig(),
					InvertedIndexConfig: invertedConfig(),
					Properties: []*models.Property{
						{
							Name:     "name",
							DataType: []string{"string"},
						},
					},
				},
				&models.Class{
					Class:               "MergeTestSource",
					VectorIndexConfig:   hnsw.NewDefaultUserConfig(),
					InvertedIndexConfig: invertedConfig(),
					Properties: []*models.Property{ // tries to have "one of each property type"
						{
							Name:     "string",
							DataType: []string{"string"},
						},
						{
							Name:     "text",
							DataType: []string{"text"},
						},
						{
							Name:     "number",
							DataType: []string{"number"},
						},
						{
							Name:     "int",
							DataType: []string{"int"},
						},
						{
							Name:     "date",
							DataType: []string{"date"},
						},
						{
							Name:     "geo",
							DataType: []string{"geoCoordinates"},
						},
						{
							Name:     "toTarget",
							DataType: []string{"MergeTestTarget"},
						},
					},
				},
			},
		},
	}

	t.Run("add required classes", func(t *testing.T) {
		for _, class := range schema.Objects.Classes {
			t.Run(fmt.Sprintf("add %s", class.Class), func(t *testing.T) {
				err := migrator.AddClass(context.Background(), class, schemaGetter.shardState)
				require.Nil(t, err)
			})
		}
	})

	schemaGetter.schema = schema

	target1 := strfmt.UUID("897be7cc-1ae1-4b40-89d9-d3ea98037638")
	target2 := strfmt.UUID("5cc94aba-93e4-408a-ab19-3d803216a04e")
	target3 := strfmt.UUID("81982705-8b1e-4228-b84c-911818d7ee85")
	target4 := strfmt.UUID("7f69c263-17f4-4529-a54d-891a7c008ca4")
	sourceID := strfmt.UUID("8738ddd5-a0ed-408d-a5d6-6f818fd56be6")

	t.Run("add objects", func(t *testing.T) {
		err := repo.PutObject(context.Background(), &models.Object{
			ID:    sourceID,
			Class: "MergeTestSource",
			Properties: map[string]interface{}{
				"string": "only the string prop set",
			},
		}, []float32{0.5})
		require.Nil(t, err)

		targets := []strfmt.UUID{target1, target2, target3, target4}

		for i, target := range targets {
			err = repo.PutObject(context.Background(), &models.Object{
				ID:    target,
				Class: "MergeTestTarget",
				Properties: map[string]interface{}{
					"name": fmt.Sprintf("target item %d", i),
				},
			}, []float32{0.5})
			require.Nil(t, err)
		}
	})

	t.Run("merge other previously unset properties into it", func(t *testing.T) {
		md := objects.MergeDocument{
			Class: "MergeTestSource",
			ID:    sourceID,
			PrimitiveSchema: map[string]interface{}{
				"number": 7.0,
				"int":    int64(9),
				"geo": &models.GeoCoordinates{
					Latitude:  ptFloat32(30.2),
					Longitude: ptFloat32(60.2),
				},
				"text": "some text",
			},
		}

		err := repo.Merge(context.Background(), md)
		assert.Nil(t, err)
	})

	t.Run("check that the object was successfully merged", func(t *testing.T) {
		source, err := repo.ObjectByID(context.Background(), sourceID, nil, additional.Properties{})
		require.Nil(t, err)

		schema := source.Object().Properties.(map[string]interface{})
		expectedSchema := map[string]interface{}{
			// from original
			"string": "only the string prop set",

			// from merge
			"number": 7.0,
			"int":    float64(9),
			"geo": &models.GeoCoordinates{
				Latitude:  ptFloat32(30.2),
				Longitude: ptFloat32(60.2),
			},
			"text": "some text",
		}

		assert.Equal(t, expectedSchema, schema)
	})

	t.Run("trying to merge from unexisting index", func(t *testing.T) {
		md := objects.MergeDocument{
			Class: "WrongClass",
			ID:    sourceID,
			PrimitiveSchema: map[string]interface{}{
				"number": 7.0,
			},
		}

		err := repo.Merge(context.Background(), md)
		assert.Equal(t, fmt.Errorf(
			"merge from non-existing index for WrongClass"), err)
	})
	t.Run("add a reference and replace one prop", func(t *testing.T) {
		source, err := crossref.ParseSource(fmt.Sprintf(
			"weaviate://localhost/MergeTestSource/%s/toTarget", sourceID))
		require.Nil(t, err)
		targets := []strfmt.UUID{target1}
		refs := make(objects.BatchReferences, len(targets), len(targets))
		for i, target := range targets {
			to, err := crossref.Parse(fmt.Sprintf("weaviate://localhost/%s", target))
			require.Nil(t, err)
			refs[i] = objects.BatchReference{
				Err:  nil,
				From: source,
				To:   to,
			}
		}
		md := objects.MergeDocument{
			Class: "MergeTestSource",
			ID:    sourceID,
			PrimitiveSchema: map[string]interface{}{
				"string": "let's update the string prop",
			},
			References: refs,
		}
		err = repo.Merge(context.Background(), md)
		assert.Nil(t, err)
	})

	t.Run("check that the object was successfully merged", func(t *testing.T) {
		source, err := repo.ObjectByID(context.Background(), sourceID, nil, additional.Properties{})
		require.Nil(t, err)

		ref, err := crossref.Parse(fmt.Sprintf("weaviate://localhost/%s", target1))
		require.Nil(t, err)

		schema := source.Object().Properties.(map[string]interface{})
		expectedSchema := map[string]interface{}{
			"string": "let's update the string prop",
			"number": 7.0,
			"int":    float64(9),
			"geo": &models.GeoCoordinates{
				Latitude:  ptFloat32(30.2),
				Longitude: ptFloat32(60.2),
			},
			"text": "some text",
			"toTarget": models.MultipleRef{
				ref.SingleRef(),
			},
		}

		assert.Equal(t, expectedSchema, schema)
	})

	t.Run("add more references in rapid succession", func(t *testing.T) {
		// this test case prevents a regression on gh-1016
		source, err := crossref.ParseSource(fmt.Sprintf(
			"weaviate://localhost/MergeTestSource/%s/toTarget", sourceID))
		require.Nil(t, err)
		targets := []strfmt.UUID{target2, target3, target4}
		refs := make(objects.BatchReferences, len(targets), len(targets))
		for i, target := range targets {
			to, err := crossref.Parse(fmt.Sprintf("weaviate://localhost/%s", target))
			require.Nil(t, err)
			refs[i] = objects.BatchReference{
				Err:  nil,
				From: source,
				To:   to,
			}
		}
		md := objects.MergeDocument{
			Class:      "MergeTestSource",
			ID:         sourceID,
			References: refs,
		}
		err = repo.Merge(context.Background(), md)
		assert.Nil(t, err)
	})

	t.Run("check all references are now present", func(t *testing.T) {
		source, err := repo.ObjectByID(context.Background(), sourceID, nil, additional.Properties{})
		require.Nil(t, err)

		refs := source.Object().Properties.(map[string]interface{})["toTarget"]
		refsSlice, ok := refs.(models.MultipleRef)
		require.True(t, ok, fmt.Sprintf("toTarget must be models.MultipleRef, but got %#v", refs))

		foundBeacons := []string{}
		for _, ref := range refsSlice {
			foundBeacons = append(foundBeacons, ref.Beacon.String())
		}
		expectedBeacons := []string{
			fmt.Sprintf("weaviate://localhost/%s", target1),
			fmt.Sprintf("weaviate://localhost/%s", target2),
			fmt.Sprintf("weaviate://localhost/%s", target3),
			fmt.Sprintf("weaviate://localhost/%s", target4),
		}

		assert.ElementsMatch(t, foundBeacons, expectedBeacons)
	})
}

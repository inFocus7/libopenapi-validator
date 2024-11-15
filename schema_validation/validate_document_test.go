// Copyright 2023 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package schema_validation

import (
	"github.com/pb33f/libopenapi/datamodel"
	"os"
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/stretchr/testify/assert"
)

func TestValidateDocument(t *testing.T) {
	petstore, _ := os.ReadFile("../test_specs/petstorev3.json")

	doc, _ := libopenapi.NewDocument(petstore)

	// validate!
	valid, errors := ValidateOpenAPIDocument(doc)

	assert.True(t, valid)
	assert.Len(t, errors, 0)
}

func TestValidateDocument_31(t *testing.T) {
	petstore, _ := os.ReadFile("../test_specs/valid_31.yaml")

	doc, _ := libopenapi.NewDocument(petstore)

	// validate!
	valid, errors := ValidateOpenAPIDocument(doc)

	assert.True(t, valid)
	assert.Len(t, errors, 0)
}

func TestValidateDocument_Invalid31(t *testing.T) {
	petstore, _ := os.ReadFile("../test_specs/invalid_31.yaml")

	doc, _ := libopenapi.NewDocument(petstore)

	// validate!
	valid, errors := ValidateOpenAPIDocument(doc)

	assert.False(t, valid)
	assert.Len(t, errors, 1)
	assert.Len(t, errors[0].SchemaValidationErrors, 6)
}

// This is the edge case that fails
// TODO: Make sure this passes. Would require update to libopenapi to not ignore errors during decoding (if that's safe to do).
func TestValidateDocument_Invalid3_dupe_op_102(t *testing.T) {
	petstore, _ := os.ReadFile("../test_specs/final_path_malindent_with_dupe_op_key.yaml")

	// Because we ignore any errors during decoding json/yaml(?), we DO NOT get issues here.
	// The error we'd expect here is 'duplicate key' due to "/get" being repeated, since the second path is indented as if it was under/inside the first path.
	// This is the same error we see when using the same example through the swagger editor website.
	doc, err := libopenapi.NewDocumentWithConfiguration(petstore, &datamodel.DocumentConfiguration{
		BypassDocumentCheck: false,
	})
	// we should have error'd due to misconfiguration in input file - the second path being indented technically puts its 'get' operation on the same indent level as the first path's 'get' operation, so it's a duplicate key and would fail decoding.
	assert.Error(t, err)
	// unsure if we should ideally still return the doc or nil doc with error.

	// This doesn't error, and i guess that makes sense since the final parsed doc/json dropped the paths value which was the issue, and this just validates the schema.
	// We should ideally/for-sure error during the decoding of the json/yaml, but we don't.
	valid, errors := ValidateOpenAPIDocument(doc)

	assert.True(t, valid)
	assert.Len(t, errors, 0)
}

// This is another somewhat similar case, but passes decoding and fails validation. This test case is expected, it was me verifying stuff.
// - yaml/json decoding (libopenapi) shouldn't fail as it's valid (no dupe keys)
// - openapi validation (this library) should fails because it's invalid spec
func TestValidateDocument_Invalid3_non_dupe_op_102(t *testing.T) {
	petstore, _ := os.ReadFile("../test_specs/final_path_malindent_with_non_dupe_op_key.yaml")

	doc, err := libopenapi.NewDocumentWithConfiguration(petstore, &datamodel.DocumentConfiguration{
		BypassDocumentCheck: false,
	})
	assert.Nil(t, err)

	// This doesn't error, and i guess that makes sense since the final parsed doc/json dropped the paths value which was the issue, and this just validates the schema.
	// We should ideally/for-sure error during the decoding of the json/yaml, but we don't.
	valid, errors := ValidateOpenAPIDocument(doc)

	assert.False(t, valid)
	assert.Len(t, errors, 1)
}

func TestValidateSchema_ValidateLicenseIndentifier(t *testing.T) {
	spec := `openapi: 3.1.0
info:
  version: 1.0.0
  title: Test
  license:
    name: Apache 2.0
    url: https://opensource.org/licenses/Apache-2.0
    identifier: Apache-2.0
components:
  schemas:
    Pet:
      type: string`

	doc, _ := libopenapi.NewDocument([]byte(spec))

	// validate!
	valid, errors := ValidateOpenAPIDocument(doc)

	assert.False(t, valid)
	assert.Len(t, errors, 1)
	assert.Len(t, errors[0].SchemaValidationErrors, 1)
}

func TestValidateSchema_GeneratePointlessValidation(t *testing.T) {
	spec := `openapi: 3.1.0
info:
  version: 1
`

	doc, _ := libopenapi.NewDocument([]byte(spec))

	// validate!
	valid, errors := ValidateOpenAPIDocument(doc)

	assert.False(t, valid)
	assert.Len(t, errors, 1)
	assert.Len(t, errors[0].SchemaValidationErrors, 6)
}

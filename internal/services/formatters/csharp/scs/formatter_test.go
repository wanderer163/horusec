// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	analysisEntities "github.com/ZupIT/horusec-devkit/pkg/entities/analysis"

	cliConfig "github.com/ZupIT/horusec/config"
	"github.com/ZupIT/horusec/internal/entities/toolsconfig"
	"github.com/ZupIT/horusec/internal/entities/workdir"
	"github.com/ZupIT/horusec/internal/services/docker"
	"github.com/ZupIT/horusec/internal/services/formatters"
	"github.com/ZupIT/horusec/internal/services/formatters/csharp/scs/enums"
)

func createSlnFile(t *testing.T) {
	wd, _ := os.Getwd()

	dir := filepath.Join(wd, ".horusec", "00000000-0000-0000-0000-000000000000")
	err := os.MkdirAll(dir, 0755)
	require.Nil(t, err, "Expected nil error to create sln directory: %v", err)

	f, err := os.Create(filepath.Join(dir, "test.sln"))
	require.Nil(t, err, "Expected nil error to create sln file: %v", err)

	t.Cleanup(func() {
		assert.NoError(t, f.Close(), "Expected nil error to close sln file")
		assert.NoError(t, os.RemoveAll(dir), "Expected nil error to remove sln directory")
	})

}

func TestParseOutput(t *testing.T) {
	t.Run("should return 4 vulnerabilities with no errors", func(t *testing.T) {
		createSlnFile(t)

		dockerAPIControllerMock := &docker.Mock{}
		dockerAPIControllerMock.On("SetAnalysisID")
		analysis := &analysisEntities.Analysis{}
		config := cliConfig.New()
		config.WorkDir = &workdir.WorkDir{}

		output := "{ \"$schema\": \"https://schemastore.azurewebsites.net/schemas/json/sarif-2.1.0-rtm.5.json\"," +
			" \"version\": \"2.1.0\", \"runs\": [ { \"results\": [ { \"ruleId\": \"SCS0006\", \"ruleIndex\": 0, " +
			"\"level\": \"warning\", \"message\": { \"text\": \"Weak hashing function.\" }, \"locations\": [ " +
			"{ \"physicalLocation\": { \"artifactLocation\": { \"uri\": " +
			"\"file:///src/NetCoreVulnerabilities/Vulnerabilities.cs\" }, \"region\": { \"startLine\": 22, " +
			"\"startColumn\": 32, \"endLine\": 22, \"endColumn\": 63 } } } ], \"properties\": { \"warningLevel\": " +
			"1 } }, { \"ruleId\": \"SCS0006\", \"ruleIndex\": 0, \"level\": \"warning\", \"message\": { \"text\":" +
			" \"Weak hashing function.\" }, \"locations\": [ { \"physicalLocation\": { \"artifactLocation\":" +
			" { \"uri\": \"file:///src/NetCoreVulnerabilities/Vulnerabilities.cs\" }, \"region\": { \"startLine\": " +
			"15, \"startColumn\": 32, \"endLine\": 15, \"endColumn\": 63 } } } ], \"properties\":" +
			" { \"warningLevel\": 1 } }, { \"ruleId\": \"SCS0005\", \"ruleIndex\": 1, \"level\": \"warning\"," +
			" \"message\": { \"text\": \"Weak random number generator.\" }, \"locations\": [ { \"physicalLocation\":" +
			" { \"artifactLocation\": { \"uri\": \"file:///src/NetCoreVulnerabilities/Vulnerabilities.cs\" }," +
			" \"region\": { \"startLine\": 37, \"startColumn\": 13, \"endLine\": 37, \"endColumn\": 26 } } } ]," +
			" \"properties\": { \"warningLevel\": 1 } }, { \"ruleId\": \"\", \"ruleIndex\": 2, \"level\":" +
			" \"warning\", \"message\": { \"text\": \"Hardcoded value in 'string password'.\" }, \"locations\": [" +
			" { \"physicalLocation\": { \"artifactLocation\": { \"uri\": " +
			"\"file:///src/NetCoreVulnerabilities/Vulnerabilities.cs\" }, \"region\": { \"startLine\": 28, " +
			"\"startColumn\": 34, \"endLine\": 28, \"endColumn\": 88 } } } ], \"relatedLocations\": [ {" +
			" \"physicalLocation\": { \"artifactLocation\": { \"uri\": " +
			"\"file:///src/NetCoreVulnerabilities/Vulnerabilities.cs\" }, \"region\": { \"startLine\": 28," +
			" \"startColumn\": 34, \"endLine\": 28, \"endColumn\": 88 } } } ], \"properties\": { \"warningLevel\":" +
			" 1 } } ], \"tool\": { \"driver\": { \"name\": \"Security Code Scan\", \"version\": \"5.1.1.0\"," +
			" \"dottedQuadFileVersion\": \"5.1.1.0\", \"semanticVersion\": \"5.1.1\", \"language\": \"\", \"rules\":" +
			" [ { \"id\": \"SCS0006\", \"shortDescription\": { \"text\": \"Weak hashing function.\" }," +
			" \"fullDescription\": { \"text\": \"SHA1 is no longer considered as a strong hashing algorithm.\" }," +
			" \"helpUri\": \"https://security-code-scan.github.io/#SCS0006\", \"properties\":" +
			" { \"category\": \"Security\" } }, { \"id\": \"SCS0005\", \"shortDescription\": { \"text\":" +
			" \"Weak random number generator.\" }, \"fullDescription\": { \"text\": \"It is possible to predict" +
			" the next numbers of a pseudo random generator. Use a cryptographically strong generator for security" +
			" sensitive purposes.\" }, \"helpUri\": \"https://security-code-scan.github.io/#SCS0005\"," +
			" \"properties\": { \"category\": \"Security\" } }, { \"id\": \"SCS0015\", \"shortDescription\":" +
			" { \"text\": \"Hardcoded value in '{0}'.\" }, \"fullDescription\": { \"text\":" +
			" \"The secret value to this API appears to be hardcoded. Consider moving the value" +
			" to externalized configuration to avoid leakage of secret information.\" }, \"helpUri\": " +
			"\"https://security-code-scan.github.io/#SCS0015\", \"properties\": { \"category\": \"Security\" " +
			"} } ] } }, \"columnKind\": \"utf16CodeUnits\" } ] }"

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return(output, nil)

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config)
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
		assert.Len(t, analysis.AnalysisVulnerabilities, 4)

	})

	t.Run("should return error when unmarshalling output", func(t *testing.T) {
		createSlnFile(t)

		analysis := &analysisEntities.Analysis{}
		dockerAPIControllerMock := &docker.Mock{}
		dockerAPIControllerMock.On("SetAnalysisID")
		config := &cliConfig.Config{}
		config.WorkDir = &workdir.WorkDir{}

		output := "test"

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return(output, nil)

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config)
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
		assert.NotEmpty(t, analysis.Errors)

	})

	t.Run("should return build error", func(t *testing.T) {
		createSlnFile(t)

		analysis := &analysisEntities.Analysis{}
		dockerAPIControllerMock := &docker.Mock{}
		dockerAPIControllerMock.On("SetAnalysisID")
		config := &cliConfig.Config{}
		config.WorkDir = &workdir.WorkDir{}

		output := enums.BuildFailedOutput

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return(output, nil)

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config)
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
		assert.NotEmpty(t, analysis.Errors)

	})

	t.Run("should return error not found solution file", func(t *testing.T) {
		analysis := &analysisEntities.Analysis{}
		dockerAPIControllerMock := &docker.Mock{}
		dockerAPIControllerMock.On("SetAnalysisID")
		config := &cliConfig.Config{}
		config.WorkDir = &workdir.WorkDir{}

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return("", nil)

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config)
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
		assert.NotEmpty(t, analysis.Errors)
	})

	t.Run("should return error executing container", func(t *testing.T) {
		createSlnFile(t)

		analysis := &analysisEntities.Analysis{}
		dockerAPIControllerMock := &docker.Mock{}
		dockerAPIControllerMock.On("SetAnalysisID")
		config := &cliConfig.Config{}
		config.WorkDir = &workdir.WorkDir{}

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").
			Return("", errors.New("test"))

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config)
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
		assert.NotEmpty(t, analysis.Errors)
	})

	t.Run("should not execute tool because it's ignored", func(t *testing.T) {
		analysis := &analysisEntities.Analysis{}
		dockerAPIControllerMock := &docker.Mock{}
		config := &cliConfig.Config{}
		config.WorkDir = &workdir.WorkDir{}
		config.ToolsConfig = toolsconfig.ParseInterfaceToMapToolsConfig(
			toolsconfig.ToolsConfigsStruct{SecurityCodeScan: toolsconfig.ToolConfig{IsToIgnore: true}},
		)

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config)
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
	})
}

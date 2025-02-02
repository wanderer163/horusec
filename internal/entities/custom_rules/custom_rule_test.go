// Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

package customrules

import (
	"testing"

	"github.com/ZupIT/horusec-devkit/pkg/enums/languages"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/horusec-devkit/pkg/enums/confidence"
	"github.com/ZupIT/horusec-devkit/pkg/enums/severities"
	"github.com/ZupIT/horusec-engine/text"
	customRulesEnums "github.com/ZupIT/horusec/internal/enums/custom_rules"
)

func TestValidate(t *testing.T) {
	t.Run("should return no errors when valid custom rule", func(t *testing.T) {
		customRule := CustomRule{
			ID:          "HS-LEAKS-1000",
			Name:        "test",
			Description: "test",
			Severity:    severities.Low,
			Confidence:  confidence.Low,
			Type:        customRulesEnums.OrMatch,
			Expressions: []string{""},
			Language:    languages.Leaks,
		}

		assert.NoError(t, customRule.Validate())
	})

	t.Run("should return error when invalid custom", func(t *testing.T) {
		customRule := CustomRule{}
		assert.Error(t, customRule.Validate())
	})

	t.Run("should return error when invalid ID", func(t *testing.T) {
		customRule := CustomRule{
			ID:          "HS-INVALID-1",
			Name:        "test",
			Description: "test",
			Severity:    severities.Low,
			Confidence:  confidence.Low,
			Type:        customRulesEnums.Regular,
			Expressions: []string{""},
			Language:    languages.Java,
		}
		assert.Error(t, customRule.Validate())
	})
	t.Run("should return error when duplicated ID", func(t *testing.T) {
		customRule := CustomRule{
			ID:          "HS-LEAKS-1",
			Name:        "test",
			Description: "test",
			Severity:    severities.Low,
			Confidence:  confidence.Low,
			Type:        customRulesEnums.Regular,
			Expressions: []string{""},
			Language:    languages.Java,
		}
		assert.Error(t, customRule.Validate())
	})
	t.Run("should return error when not supported language", func(t *testing.T) {
		customRule := CustomRule{
			ID:          "HS-PYTHON-1",
			Name:        "test",
			Description: "test",
			Severity:    severities.Low,
			Confidence:  confidence.Low,
			Type:        customRulesEnums.Regular,
			Expressions: []string{""},
			Language:    languages.Python,
		}
		assert.Error(t, customRule.Validate())
	})
}

func TestGetRuleType(t *testing.T) {
	t.Run("should return regular type", func(t *testing.T) {
		customRule := CustomRule{
			Type: customRulesEnums.Regular,
		}

		assert.Equal(t, text.Regular, customRule.GetRuleType())
	})

	t.Run("should return regular type", func(t *testing.T) {
		customRule := CustomRule{}

		assert.Equal(t, text.Regular, customRule.GetRuleType())
	})

	t.Run("should return or type", func(t *testing.T) {
		customRule := CustomRule{
			Type: customRulesEnums.OrMatch,
		}

		assert.Equal(t, text.OrMatch, customRule.GetRuleType())
	})

	t.Run("should return and type", func(t *testing.T) {
		customRule := CustomRule{
			Type: customRulesEnums.AndMatch,
		}

		assert.Equal(t, text.AndMatch, customRule.GetRuleType())
	})

	t.Run("should return not type", func(t *testing.T) {
		customRule := CustomRule{
			Type: customRulesEnums.NotMatch,
		}

		assert.Equal(t, text.NotMatch, customRule.GetRuleType())
	})
}

func TestGetExpressions(t *testing.T) {
	t.Run("should success get regex expressions", func(t *testing.T) {
		customRule := CustomRule{
			Expressions: []string{"test", "test"},
		}

		assert.Len(t, customRule.GetExpressions(), 2)
	})

	t.Run("should log error when failed to compile expression", func(t *testing.T) {
		customRule := CustomRule{
			Expressions: []string{"^\\/(?!\\/)(.*?)"},
		}

		assert.Len(t, customRule.GetExpressions(), 0)
	})
}

func TestToString(t *testing.T) {
	t.Run("should log error when failed to compile expression", func(t *testing.T) {
		customRule := CustomRule{ID: ""}

		assert.NotEmpty(t, customRule.String())
	})
}

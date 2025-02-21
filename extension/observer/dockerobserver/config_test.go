// Copyright  OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dockerobserver

import (
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configtest"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Extensions[typeStr] = factory
	cfg, err := configtest.LoadConfigAndValidate(path.Join(".", "testdata", "config.yaml"), factories)

	require.Nil(t, err)
	require.NotNil(t, cfg)

	require.Len(t, cfg.Extensions, 2)

	ext0 := cfg.Extensions[config.NewComponentID(typeStr)]
	assert.Equal(t, factory.CreateDefaultConfig(), ext0)

	ext1 := cfg.Extensions[config.NewComponentIDWithName(typeStr, "all_settings")]
	assert.Equal(t,
		&Config{
			Endpoint:              "unix:///var/run/docker.sock",
			ExtensionSettings:     config.NewExtensionSettings(config.NewComponentIDWithName(typeStr, "all_settings")),
			CacheSyncInterval:     5 * time.Minute,
			Timeout:               20 * time.Second,
			ExcludedImages:        []string{"excluded", "image"},
			UseHostnameIfPresent:  true,
			UseHostBindings:       true,
			IgnoreNonHostBindings: true,
		},
		ext1)
}

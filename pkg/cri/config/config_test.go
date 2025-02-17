/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/containerd/containerd/plugin"
	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	for desc, test := range map[string]struct {
		config      *PluginConfig
		expectedErr string
		expected    *PluginConfig
	}{
		"no default_runtime_name": {
			config:      &PluginConfig{},
			expectedErr: "`default_runtime_name` is empty",
		},
		"no runtime[default_runtime_name]": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
				},
			},
			expectedErr: "no corresponding runtime configured in `runtimes` for `default_runtime_name`",
		},
		"deprecated systemd_cgroup for v1 runtime": {
			config: &PluginConfig{
				SystemdCgroup: true,
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
			expected: &PluginConfig{
				SystemdCgroup: true,
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
		},
		"deprecated systemd_cgroup for v2 runtime": {
			config: &PluginConfig{
				SystemdCgroup: true,
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeRuncV1,
						},
					},
				},
			},
			expectedErr: fmt.Sprintf("`systemd_cgroup` only works for runtime %s", plugin.RuntimeLinuxV1),
		},
		"no_pivot for v1 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					NoPivot:            true,
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					NoPivot:            true,
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
		},
		"no_pivot for v2 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					NoPivot:            true,
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeRuncV1,
						},
					},
				},
			},
			expectedErr: fmt.Sprintf("`no_pivot` only works for runtime %s", plugin.RuntimeLinuxV1),
		},
		"deprecated runtime_engine for v1 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Engine: "runc",
							Type:   plugin.RuntimeLinuxV1,
						},
					},
				},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Engine: "runc",
							Type:   plugin.RuntimeLinuxV1,
						},
					},
				},
			},
		},
		"deprecated runtime_engine for v2 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Engine: "runc",
							Type:   plugin.RuntimeRuncV1,
						},
					},
				},
			},
			expectedErr: fmt.Sprintf("`runtime_engine` only works for runtime %s", plugin.RuntimeLinuxV1),
		},
		"deprecated runtime_root for v1 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Root: "/run/containerd/runc",
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Root: "/run/containerd/runc",
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
		},
		"deprecated runtime_root for v2 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Root: "/run/containerd/runc",
							Type: plugin.RuntimeRuncV1,
						},
					},
				},
			},
			expectedErr: fmt.Sprintf("`runtime_root` only works for runtime %s", plugin.RuntimeLinuxV1),
		},
		"deprecated auths": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeRuncV1,
						},
					},
				},
				Registry: Registry{
					Auths: map[string]AuthConfig{
						"https://gcr.io": {Username: "test"},
					},
				},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeRuncV1,
						},
					},
				},
				Registry: Registry{
					Configs: map[string]RegistryConfig{
						"gcr.io": {
							Auth: &AuthConfig{
								Username: "test",
							},
						},
					},
					Auths: map[string]AuthConfig{
						"https://gcr.io": {Username: "test"},
					},
				},
			},
		},
		"invalid stream_idle_timeout": {
			config: &PluginConfig{
				StreamIdleTimeout: "invalid",
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: "default",
						},
					},
				},
			},
			expectedErr: "invalid stream idle timeout",
		},
		"conflicting mirror registry config": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: "default",
						},
					},
				},
				Registry: Registry{
					ConfigPath: "/etc/containerd/conf.d",
					Mirrors: map[string]Mirror{
						"something.io": {},
					},
				},
			},
			expectedErr: "`mirrors` cannot be set when `config_path` is provided",
		},
		"conflicting tls registry config": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: "default",
						},
					},
				},
				Registry: Registry{
					ConfigPath: "/etc/containerd/conf.d",
					Configs: map[string]RegistryConfig{
						"something.io": {
							TLS: &TLSConfig{},
						},
					},
				},
			},
			expectedErr: "`configs.tls` cannot be set when `config_path` is provided",
		},
	} {
		t.Run(desc, func(t *testing.T) {
			err := ValidatePluginConfig(context.Background(), test.config)
			if test.expectedErr != "" {
				assert.Contains(t, err.Error(), test.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, test.config)
			}
		})
	}
}

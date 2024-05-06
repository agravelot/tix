package core_test

import (
	"testing"

	"github.com/agravelot/tix/core"
)

func TestCoreApp(t *testing.T) {
	type testCase struct {
		name     string
		envShell string
		input    struct{ cfg core.Config }
		expected struct{ app core.Application }
		error    bool
	}

	tests := []testCase{
		{
			name:     "should start with empty config",
			input:    struct{ cfg core.Config }{cfg: core.Config{}},
			error:    false,
			envShell: "/bin/bash",
		},
		{
			name:     "should fail with no shell",
			input:    struct{ cfg core.Config }{cfg: core.Config{}},
			error:    true,
			envShell: "",
		},
		{
			name: "should set default workspace timeout",
			input: struct{ cfg core.Config }{cfg: core.Config{
				Workspaces: []core.ConfigWorkspace{
					{Timeout: 5},
					{},
					{Timeout: 3},
				},
			}},
			expected: struct{ app core.Application }{app: core.Application{
				Config: core.Config{
					Workspaces: []core.ConfigWorkspace{
						{Timeout: 5},
						{Timeout: 5},
						{Timeout: 3},
					},
				},
			}},
			error:    false,
			envShell: "/bin/bash",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			t.Setenv("SHELL", tc.envShell)

			app, err := core.NewApplication(tc.input.cfg)
			if (err != nil) != tc.error {
				t.Fatalf(`unexpected error state got: %v`, err)
			}
			if tc.envShell != app.Config.Shell {
				t.Fatalf(`expected %v, got %v`, tc.envShell, app.Config.Shell)
			}
			if len(tc.expected.app.Workspaces) != len(app.Workspaces) {
				t.Fatalf(`expected %v, got %v`, len(tc.expected.app.Workspaces), len(app.Workspaces))
			}
			for i, w := range app.Workspaces {
				if w.Timeout != tc.expected.app.Workspaces[i].Timeout {
					t.Fatalf(`expected 5, got %v`, w.Timeout)
				}
			}
		})
	}
}

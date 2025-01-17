// Copyright 2016-2020, Pulumi Corporation.
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

package toolchain

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsVirtualEnv(t *testing.T) {
	t.Parallel()

	// Create a new empty test directory.
	tempdir := t.TempDir()

	// Assert the empty test directory is not a virtual environment.
	assert.False(t, IsVirtualEnv(tempdir))

	// Create and run a python command to create a virtual environment.
	venvDir := filepath.Join(tempdir, "venv")
	cmd, err := Command(context.Background(), "-m", "venv", venvDir)
	assert.NoError(t, err)
	err = cmd.Run()
	assert.NoError(t, err)

	// Assert the new venv directory is a virtual environment.
	assert.True(t, IsVirtualEnv(venvDir))
}

func TestActivateVirtualEnv(t *testing.T) {
	t.Parallel()

	venvName := "venv"
	venvDir := filepath.Join(venvName, "bin")
	if runtime.GOOS == windows {
		venvDir = filepath.Join(venvName, "Scripts")
	}

	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"PYTHONHOME=foo", "PATH=bar", "FOO=blah"},
			expected: []string{fmt.Sprintf("PATH=%s%sbar", venvDir, string(os.PathListSeparator)), "FOO=blah"},
		},
		{
			input:    []string{"PYTHONHOME=foo", "FOO=blah"},
			expected: []string{"FOO=blah", "PATH=" + venvDir},
		},
		{
			input:    []string{"PythonHome=foo", "Path=bar"},
			expected: []string{fmt.Sprintf("Path=%s%sbar", venvDir, string(os.PathListSeparator))},
		},
	}
	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%#v", test.input), func(t *testing.T) {
			t.Parallel()

			actual := ActivateVirtualEnv(test.input, venvName)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestRunningPipInVirtualEnvironment(t *testing.T) {
	t.Parallel()

	// Skip during short test runs since this test involves downloading dependencies.
	if testing.Short() {
		t.Skip("Skipped in short test run")
	}

	// Create a new empty test directory.
	tempdir := t.TempDir()

	// Create and run a python command to create a virtual environment.
	venvDir := filepath.Join(tempdir, "venv")
	cmd, err := Command(context.Background(), "-m", "venv", venvDir)
	assert.NoError(t, err)
	err = cmd.Run()
	assert.NoError(t, err)

	// Create a requirements.txt file in the temp directory.
	requirementsFile := filepath.Join(tempdir, "requirements.txt")
	assert.NoError(t, os.WriteFile(requirementsFile, []byte("pulumi==2.0.0\n"), 0o600))

	// Create a command to run pip from the virtual environment.
	pipCmd := VirtualEnvCommand(venvDir, "python", "-m", "pip", "install", "-r", "requirements.txt")
	pipCmd.Dir = tempdir
	pipCmd.Env = ActivateVirtualEnv(os.Environ(), venvDir)

	// Run the command.
	if output, err := pipCmd.CombinedOutput(); err != nil {
		assert.Failf(t, "pip install command failed with output: %s", string(output))
	}
}

//nolint:paralleltest // modifies environment variables
func TestCommand(t *testing.T) {
	tmp := t.TempDir()

	t.Setenv("MY_ENV_VAR", "HELLO")

	tc, err := newPip(tmp, "venv")
	require.NoError(t, err)

	cmd, err := tc.Command(context.Background())
	require.NoError(t, err)

	var venvBin string
	if runtime.GOOS == windows {
		venvBin = filepath.Join(tmp, "venv", "Scripts")
		require.Equal(t, filepath.Join("python.exe"), cmd.Path)
	} else {
		venvBin = filepath.Join(tmp, "venv", "bin")
		require.Equal(t, filepath.Join(venvBin, "python"), cmd.Path)
	}

	foundPath := false
	foundMyEnvVar := false
	for _, env := range cmd.Env {
		if strings.HasPrefix(env, "PATH=") {
			require.Contains(t, env, venvBin, "venv binary directory should in PATH")
			foundPath = true
		}
		if strings.HasPrefix(env, "MY_ENV_VAR=") {
			require.Equal(t, "MY_ENV_VAR=HELLO", env, "Env variables should be passed through")
			foundMyEnvVar = true
		}
	}
	require.True(t, foundPath, "PATH was not found in cmd.Env")
	require.True(t, foundMyEnvVar, "MY_ENV_VAR was not found in cmd.Env")
}

func TestCommandNoVenv(t *testing.T) {
	t.Parallel()

	tc, err := newPip(".", "")
	require.NoError(t, err)

	cmd, err := tc.Command(context.Background())
	require.NoError(t, err)

	globalPython, err := exec.LookPath("python3")
	require.NoError(t, err)

	require.Equal(t, globalPython, cmd.Path, "Toolchain should use the global python executable")

	require.Nil(t, cmd.Env)
}

//nolint:paralleltest // modifies environment variables
func TestCommandPulumiPythonCommand(t *testing.T) {
	t.Setenv("PULUMI_PYTHON_CMD", "python-not-found")

	tc, err := newPip(".", "")
	require.NoError(t, err)

	cmd, err := tc.Command(context.Background())
	require.ErrorContains(t, err, "python-not-found")
	require.Nil(t, cmd)
}

func TestValidateVenv(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()
	tc, err := newPip(tmp, "mytestvenv")
	require.NoError(t, err)
	err = tc.ValidateVenv(context.Background())
	require.ErrorContains(t, err, "The 'virtualenv' option in Pulumi.yaml is set to \"mytestvenv\"")
}

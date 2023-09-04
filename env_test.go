// Copyright 2023 Jean-Philippe Couture
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Setup and teardown for each test
func setupTest(tb testing.TB) func(tb testing.TB) {
	tb.Helper()

	// Making sure the environment is clean and predictable
	os.Clearenv()

	// Setting a few environment variables for testing
	os.Setenv("PATH", "/usr/bin:/bin:/usr/sbin:/sbin")
	os.Setenv(" LANG ", "en_CA.UTF-8") // Intentionnally funky variable name

	return func(tb testing.TB) {
		// Leaving the environment clean (probably not necessary)
		os.Clearenv()
	}
}

func TestClear(t *testing.T) {
	cases := []struct {
		description string
		exceptions  []string
		expected    int
	}{
		{"Clear without exceptions", []string{}, 0},
		{"Clear with exceptions", []string{"PATH"}, 1},
		{"Clear with non-existing exceptions", []string{"BAR"}, 0},
		{"Clear with empty exceptions", []string{""}, 0},
		{"Clear with whitespace exceptions", []string{" "}, 0},
		{"Clear with different case exceptions", []string{"path"}, 0},
		{"Clear with untrimmed exceptions", []string{" LANG "}, 1},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			testdownTest := setupTest(t)
			defer testdownTest(t)

			Clear(tt.exceptions...)
			result := len(os.Environ())
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExists(t *testing.T) {
	cases := []struct {
		description string
		input       string
		expected    bool
	}{
		{"Lookup with defined variable", "PATH", true},
		{"Lookup with undefined variable", "BAR", false},
		{"Lookup with empty variable", "", false},
		{"Lookup with whitespace variable", " ", false},
		{"Lookup with different case variable", "path", false},
		{"Lookup with untrimmed variable", " LANG ", true},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			testdownTest := setupTest(t)
			defer testdownTest(t)

			result := Exists(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetvars(t *testing.T) {
	cases := []struct {
		description string
		expected    map[string]string
	}{
		{"Retrieve all environment variables", map[string]string{"PATH": "/usr/bin:/bin:/usr/sbin:/sbin", " LANG ": "en_CA.UTF-8"}},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			testdownTest := setupTest(t)
			defer testdownTest(t)

			result := Getvars()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetnames(t *testing.T) {
	cases := []struct {
		description string
		expected    []string
	}{
		{"Extract all environment variable names", []string{"PATH", " LANG "}},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			testdownTest := setupTest(t)
			defer testdownTest(t)

			result := Getnames(Getvars())
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJoin(t *testing.T) {
	cases := []struct {
		description string
		base        map[string]string
		override    map[string]string
		expected    map[string]string
	}{
		{"Join with empty base and override", map[string]string{}, map[string]string{}, map[string]string{}},
		{"Join with empty base and non-empty override", map[string]string{}, map[string]string{"COLOR": "RED"}, map[string]string{"COLOR": "RED"}},
		{"Join with non-empty base and empty override", map[string]string{"FOO": "BAR"}, map[string]string{}, map[string]string{"FOO": "BAR"}},
		{"Join with non-empty base and override", map[string]string{"FOO": "BAR"}, map[string]string{"COLOR": "RED"}, map[string]string{"FOO": "BAR", "COLOR": "RED"}},
		{"Join with non-empty base and override with existing variable", map[string]string{"FOO": "BAR"}, map[string]string{"FOO": "BAZ"}, map[string]string{"FOO": "BAZ"}},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			result := Join(tt.base, tt.override)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSetvars(t *testing.T) {
	cases := []struct {
		description string
		clearEnv    bool
		input       map[string]string
		expected    map[string]string
	}{
		{"Set empty environment", true, map[string]string{}, map[string]string{}},
		{"Set non-empty environment", false, map[string]string{"COLOR": "RED"}, map[string]string{"COLOR": "RED", "PATH": "/usr/bin:/bin:/usr/sbin:/sbin", " LANG ": "en_CA.UTF-8"}},
		{"Set non-empty environment with existing variable", false, map[string]string{"PATH": "/usr/bin"}, map[string]string{"PATH": "/usr/bin", " LANG ": "en_CA.UTF-8"}},
		{"Set non-empty environment with existing variable with different case", false, map[string]string{"path": "/usr/bin"}, map[string]string{"path": "/usr/bin", "PATH": "/usr/bin:/bin:/usr/sbin:/sbin", " LANG ": "en_CA.UTF-8"}},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			testdownTest := setupTest(t)
			defer testdownTest(t)

			if tt.clearEnv {
				Clear()
			}

			Setvars(tt.input)
			result := Getvars()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContains(t *testing.T) {
	cases := []struct {
		description string
		input       []string
		element     string
		expected    bool
	}{
		{"Contains with empty slice", []string{}, "FOO", false},
		{"Contains with non-empty slice and non-existing element", []string{"BAR"}, "FOO", false},
		{"Contains with non-empty slice and existing element", []string{"FOO"}, "FOO", true},
		{"Contains with non-empty slice and existing element with different case", []string{"FOO"}, "foo", false},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			result := contains(tt.input, tt.element)
			assert.Equal(t, tt.expected, result)
		})
	}
}

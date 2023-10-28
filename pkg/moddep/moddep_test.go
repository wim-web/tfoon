package moddep_test

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/wim-web/tfoon/pkg/moddep"
)

func testdataDir() string {
	c, _ := os.Getwd()
	return path.Join(c, "../../testdata")
}

func TestFromPath(t *testing.T) {
	t.Run("assert path", func(t *testing.T) {
		// table driven test
		// pathが正しく設定されていることを確認する
		type testCase struct {
			name     string
			path     string
			expected string
		}

		testCases := []testCase{
			{
				name:     "path with trailing slash",
				path:     path.Join(testdataDir(), "terraform/caller1/"),
				expected: path.Join(testdataDir(), "terraform/caller1"),
			},
			{
				name:     "path without trailing slash",
				path:     path.Join(testdataDir(), "terraform/caller1"),
				expected: path.Join(testdataDir(), "terraform/caller1"),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tree, err := moddep.FromPath(tc.path)

				if err != nil {
					t.Error(err)
				}

				if tree.Path != tc.expected {
					t.Errorf("tree.Path = %s, want %s", tree.Path, tc.expected)
				}
			})
		}
	})

	// caller3はmodules/noopを呼び出している
	// caller3のModuleTreeのChildrenにはModuleTree{Path: "modules/noop"}が含まれていることを確認する
	t.Run("assert children", func(t *testing.T) {
		tree, err := moddep.FromPath(path.Join(testdataDir(), "terraform/caller3"))
		if err != nil {
			t.Error(err)
		}

		if len(tree.Children) != 1 {
			t.Errorf("len(tree.Children) = %d, want %d", len(tree.Children), 1)
		}

		expected := path.Join(testdataDir(), "terraform/modules/noop")

		if tree.Children[0].Path != expected {
			t.Errorf("tree.Children[0].Path = %s, want %s", tree.Children[0].Path, expected)
		}
	})
}

func TestFromPaths(t *testing.T) {
	t.Run("assert len", func(t *testing.T) {
		// table driven test
		// pathsの数だけModuleTreeが生成されることを確認する
		type testCase struct {
			name     string
			paths    []string
			expected int
		}

		testCases := []testCase{
			{
				name:     "single path",
				paths:    []string{path.Join(testdataDir(), "terraform/caller1")},
				expected: 1,
			},
			{
				name:     "multiple paths",
				paths:    []string{path.Join(testdataDir(), "terraform/caller1"), path.Join(testdataDir(), "terraform/caller2")},
				expected: 2,
			},
			{
				name:     "no paths",
				paths:    []string{},
				expected: 0,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				list, err := moddep.FromPaths(tc.paths)

				if err != nil {
					t.Error(err)
				}

				if len(list) != tc.expected {
					t.Errorf("len(list) = %d, want %d", len(list), tc.expected)
				}
			})
		}
	})
}

// To2ModuleEntryPointはModuleTreeListをModule2EntryPointに変換する
func TestToModule2EntryPoint(t *testing.T) {
	tree, err := moddep.FromPaths([]string{
		path.Join(testdataDir(), "terraform/caller3"),
	})
	if err != nil {
		t.Error(err)
	}

	m2p := tree.ToModule2EntryPoint()

	// caller3はterraform/modules/noopを呼び出しているので
	// terraform/modules/noop => [terraform/caller3]が返ってくることを確認する

	actual := m2p[path.Join(testdataDir(), "terraform/modules/noop")]
	expected := []string{path.Join(testdataDir(), "terraform/caller3")}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v, want %v", actual, expected)
	}
}

package moddep

import (
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

type ModuleTreeList []*ModuleTree

type ModuleTree struct {
	Path     string        `json:"path"`
	Parent   *ModuleTree   `json:"-"`
	Children []*ModuleTree `json:"children"`
}

// terraformを実行する起点になるディレクトリかどうか
func (mt ModuleTree) isEntryPoint() bool {
	return mt.Parent == nil
}

func isLocalModule(m *tfconfig.ModuleCall) bool {
	return strings.HasPrefix(m.Source, "./") || strings.HasPrefix(m.Source, "../")
}

func newModuleTree(module *tfconfig.Module) (ModuleTree, tfconfig.Diagnostics) {
	var diags tfconfig.Diagnostics
	tree := ModuleTree{
		Path:     module.Path,
		Parent:   nil,
		Children: []*ModuleTree{},
	}

	for _, child := range module.ModuleCalls {
		if !isLocalModule(child) {
			continue
		}

		childPath := path.Join(module.Path, child.Source)
		moduleOfChild, diags := tfconfig.LoadModule(childPath)
		if diags.HasErrors() {
			return ModuleTree{}, diags
		}

		childModuleTree, diags := newModuleTree(moduleOfChild)
		if diags.HasErrors() {
			return ModuleTree{}, diags
		}

		childModuleTree.Parent = &tree
		tree.Children = append(tree.Children, &childModuleTree)
	}

	return tree, diags
}

func FromPath(p string) (ModuleTree, error) {
	module, diags := tfconfig.LoadModule(p)
	if diags.HasErrors() {
		return ModuleTree{}, diags
	}

	module.Path = path.Clean(module.Path)

	tree, diags := newModuleTree(module)

	return tree, diags.Err()
}

func FromPaths(ps []string) (ModuleTreeList, error) {
	trees := make(ModuleTreeList, len(ps))

	for i, p := range ps {
		tree, err := FromPath(p)
		if err != nil {
			return nil, err
		}

		trees[i] = &tree
	}

	return trees, nil
}

type Module2EntryPoint map[string][]string

func toModule2EntryPoint(mt ModuleTree) (Module2EntryPoint, error) {
	if !mt.isEntryPoint() {
		return nil, fmt.Errorf("This module is not entry point")
	}

	// [modulePath] = Set([entryPointPath])
	m := map[string]map[string]struct{}{}

	for _, ps := range getChildrenPaths(mt.Children) {
		m[ps] = map[string]struct{}{
			mt.Path: {},
		}
	}

	m2e := make(Module2EntryPoint)

	for modulePath, entryPoints := range m {
		m2e[modulePath] = []string{}
		for entryPoint := range entryPoints {
			m2e[modulePath] = append(m2e[modulePath], entryPoint)
		}
	}

	return m2e, nil
}

func getChildrenPaths(nodes []*ModuleTree) []string {
	ps := []string{}

	for _, node := range nodes {
		ps = append(ps, node.Path)
		ps = append(ps, getChildrenPaths(node.Children)...)
	}

	return ps
}

func (mtl ModuleTreeList) ToModule2EntryPoint() Module2EntryPoint {
	m2e := make(Module2EntryPoint)

	for _, tree := range mtl {
		eachM2e, _ := toModule2EntryPoint(*tree)
		for k, v := range eachM2e {
			m2e[k] = append(m2e[k], v...)
		}
	}

	return m2e
}

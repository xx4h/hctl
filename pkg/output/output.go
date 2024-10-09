// Copyright 2024 Fabian `xx4h` Sylvester
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

package output

import (
	"fmt"
	"os"

	"github.com/gosuri/uitable"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func GetBanner() (string, error) {
	return pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("H", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithStyle("Ctl", pterm.FgWhite.ToStyle())).Srender()
}

func PrintSuccess(str string) {
	pterm.Success.Println(str)
}

func PrintSuccessAction(obj string, state string) {
	pterm.Success.Printfln("%s %s", obj, state)
}

func PrintSuccessListWithHeader(header []interface{}, list [][]interface{}) {
	table := uitable.New()
	table.AddRow(header...)
	for _, entry := range list {
		table.AddRow(entry...)
	}
	fmt.Println(table)
}

func PrintError(err error) {
	pterm.Error.Println(err)
	os.Exit(1)
}

func PrintThreeLevelFlatTree(name string, tree map[string][]string) error {
	t := pterm.TreeNode{
		Text:     name,
		Children: []pterm.TreeNode{},
	}

	for key := range tree {
		p := []pterm.TreeNode{}
		for lid := range tree[key] {
			p = append(p, pterm.TreeNode{Text: tree[key][lid]})
		}
		g := pterm.TreeNode{
			Text:     key,
			Children: p,
		}
		t.Children = append(t.Children, g)
	}

	return pterm.DefaultTree.WithRoot(t).Render()
}

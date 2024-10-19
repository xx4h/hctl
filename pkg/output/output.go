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
	"io"
	"os"
	"sort"

	"github.com/gosuri/uitable"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func GetBanner() (string, error) {
	return pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("H", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithStyle("Ctl", pterm.FgWhite.ToStyle())).Srender()
}

func FprintSuccess(out io.Writer, str string) {
	pterm.Fprint(out, pterm.Success.Sprintln(str))
}

func PrintSuccess(str string) {
	pterm.Success.Println(str)
}

func FprintSuccessAction(out io.Writer, obj string, state string) {
	pterm.Fprint(out, pterm.Success.Sprintfln("%s %s", obj, state))
}

func PrintSuccessAction(obj string, state string) {
	pterm.Success.Printfln("%s %s", obj, state)
}

func ListWithHeader(header []interface{}, list [][]interface{}) *uitable.Table {
	table := uitable.New()
	table.AddRow(header...)
	for _, entry := range list {
		table.AddRow(entry...)
	}
	return table
}

func FprintSuccessListWithHeader(out io.Writer, header []interface{}, list [][]interface{}) {
	fmt.Fprintln(out, ListWithHeader(header, list))
}

func PrintSuccessListWithHeader(header []interface{}, list [][]interface{}) {
	fmt.Println(ListWithHeader(header, list))
}

func FprintError(out io.Writer, err error) {
	pterm.Fprint(out, pterm.Error.Sprintln(err))
	os.Exit(1)
}

func PrintError(err error) {
	pterm.Error.Println(err)
	os.Exit(1)
}

func PrintThreeLevelFlatTree(out io.Writer, name string, tree map[string][]string) error {
	t := pterm.TreeNode{
		Text:     name,
		Children: []pterm.TreeNode{},
	}

	var secondLevel []string

	for key := range tree {
		secondLevel = append(secondLevel, key)
	}

	sort.Strings(secondLevel)

	for _, key := range secondLevel {
		sort.Strings(tree[key])

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

	treeout, err := pterm.DefaultTree.WithRoot(t).Srender()
	if err != nil {
		return err
	}

	fmt.Fprintln(out, treeout)
	return nil
}

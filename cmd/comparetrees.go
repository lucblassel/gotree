package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"runtime"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
)

var comparetreeidentical bool

// compareCmd represents the compare command
var compareTreesCmd = &cobra.Command{
	Use:   "trees",
	Short: "Compare a reference tree with a set of trees",
	Long: `Compare a reference tree with a set of trees.

If --binary is given:
For each trees in the compared tree file, it will print tab separated values with:
1) The index of the compared tree in the file
2) "true" if the tree is identical, 
   "false" otherwise

Otherwise:
For each trees in the compared tree file, it will print tab separated values with:
1) The index of the compared tree in the file
2) The number of branches that are specific to the reference tree
3) The number of branches that are common to both trees
4) The number of branches that are specific to the compared tree

`,
	Run: func(cmd *cobra.Command, args []string) {
		if intree2file == "none" {
			io.ExitWithMessage(errors.New("You must provide a file containing compared trees"))
		}

		maxcpus := runtime.NumCPU()
		if rootCpus > maxcpus {
			rootCpus = maxcpus
		}
		refTree := readTree(intreefile)
		compareChannel := readTrees(intree2file)
		stats := make(chan tree.BipartitionStats)
		err := tree.Compare(refTree, compareChannel, compareTips, comparetreeidentical, rootCpus, stats)
		if err != nil {
			io.ExitWithMessage(err)
		}

		if comparetreeidentical {
			fmt.Printf("tree\tidentical\n")
			for stats := range stats {
				fmt.Printf("%d\t%v\n", stats.Id, stats.Sametree)
			}
		} else {
			fmt.Printf("tree\treference\tcommon\tcompared\n")
			for stats := range stats {
				fmt.Printf("%d\t%d\t%d\t%d\n", stats.Id, stats.Tree1, stats.Common, stats.Tree2)
			}
		}
	},
}

func init() {
	compareCmd.AddCommand(compareTreesCmd)
	compareTreesCmd.Flags().BoolVarP(&compareTips, "tips", "l", false, "Include tips in the comparison")
	compareTreesCmd.Flags().BoolVar(&comparetreeidentical, "binary", false, "If true, then just print true (identical tree) or false (different tree) for each compared tree")
}

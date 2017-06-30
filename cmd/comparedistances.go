package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/support"
	"github.com/spf13/cobra"
)

// compareedgesCmd represents the compareedges command
var compareDistancesCmd = &cobra.Command{
	Use:   "distances",
	Short: "Prints transfer distance of all edges of a reference to another tree",
	Long: `Prints transfer distance of all edges of a reference to another tree.

For each reference tree in input, for each internal edge er of the reference tree, 
and for each internal edge ec of the compared tree, this command will print in tab
 separated format:
1.  tree_id
2.  er_id 
3.  ec_id
4.  transfer dist between er and ec
5.  ec_length
6.  ec_support
7.  ec_terminal
8.  ec_depth
9.  ec_topodepth
10. ec_rightname
11. taxa to move

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "Reference : %s\n", intreefile)
		fmt.Fprintf(os.Stderr, "Compared  : %s\n", intree2file)

		refTree := readTree(intreefile)
		refTree.ReinitIndexes()
		names := refTree.SortedTips()

		edges1 := refTree.Edges()
		fmt.Printf("tree_id\ter_id\tec_id\ttdist\tec_length\tec_support\tec_topodepth\tmoving_taxa\n")
		treefile, treechan := readTrees(intree2file)
		defer treefile.Close()
		for t2 := range treechan {
			if t2.Err != nil {
				io.ExitWithMessage(t2.Err)
			}
			t2.Tree.ReinitIndexes()

			edges2 := t2.Tree.Edges()
			var min_dist []uint16
			var min_dist_edges []int
			tips := refTree.Tips()
			min_dist = make([]uint16, len(edges1))
			min_dist_edges = make([]int, len(edges1))
			var i_matrix [][]uint16 = make([][]uint16, len(edges1))
			var c_matrix [][]uint16 = make([][]uint16, len(edges1))
			var hamming [][]uint16 = make([][]uint16, len(edges1))

			for i, e := range edges1 {
				e.SetId(i)
				min_dist[i] = uint16(len(tips))
				i_matrix[i] = make([]uint16, len(edges2))
				c_matrix[i] = make([]uint16, len(edges2))
				hamming[i] = make([]uint16, len(edges2))
			}
			for i, e := range edges2 {
				e.SetId(i)
			}
			support.Update_all_i_c_post_order_ref_tree(refTree, &edges1, t2.Tree, &edges2, &i_matrix, &c_matrix)
			support.Update_all_i_c_post_order_boot_tree(refTree, uint(len(tips)), &edges1, t2.Tree, &edges2, &i_matrix, &c_matrix, &hamming, &min_dist, &min_dist_edges)

			for _, e1 := range edges1 {
				if !e1.Right().Tip() {
					for _, e2 := range edges2 {
						if !e2.Right().Tip() {
							dist := hamming[e1.Id()][e2.Id()]
							if dist > uint16(len(tips))/2 {
								dist = uint16(len(tips)) - dist
							}
							var movedtaxabuf bytes.Buffer
							plus, minus := speciesToMove(e1, e2, int(dist))
							for k, sp := range plus {
								if k > 0 {
									movedtaxabuf.WriteRune(',')
								}
								movedtaxabuf.WriteRune('+')
								movedtaxabuf.WriteString(names[sp])
							}
							for k, sp := range minus {
								if k > 0 || (k == 0 && len(plus) > 0) {
									movedtaxabuf.WriteRune(',')
								}
								movedtaxabuf.WriteRune('-')
								movedtaxabuf.WriteString(names[sp])
							}
							depth, err := e2.TopoDepth()
							if err != nil {
								io.ExitWithMessage(err)
							}
							fmt.Printf("%d\t%d\t%d\t%d\t%s\t%s\t%d\t%s\n", t2.Id, e1.Id(), e2.Id(), int(dist), e2.LengthString(), e2.SupportString(), int(depth), movedtaxabuf.String())
						}
					}
				}
			}
		}
	},
}

func init() {
	compareCmd.AddCommand(compareDistancesCmd)
}

package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/support"
	"github.com/spf13/cobra"
)

var parsimonyEmpirical bool
var parsimonySeed int64

// parsimonyCmd represents the parsimony command
var parsimonyCmd = &cobra.Command{
	Use:   "parsimony",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		writeLogParsimony()
		rand.Seed(parsimonySeed)
		t, err := support.Parsimony(supportIntree, supportBoottrees, supportLog, parsimonyEmpirical, rootCpus)
		if err != nil {
			io.ExitWithMessage(err)
		}
		supportOut.WriteString(t.Newick() + "\n")
		supportLog.WriteString(fmt.Sprintf("End         : %s\n", time.Now().Format(time.RFC822)))
	},
}

func init() {
	supportCmd.AddCommand(parsimonyCmd)
	parsimonyCmd.PersistentFlags().BoolVarP(&parsimonyEmpirical, "empirical", "e", false, "If the support is computed with comparison to empirical support classical steps (shuffles of the original tree)")
	parsimonyCmd.PersistentFlags().Int64VarP(&parsimonySeed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed if empirical is ON")

}

func writeLogParsimony() {
	supportLog.WriteString("Parsimony Support\n")
	supportLog.WriteString(fmt.Sprintf("Date        : %s\n", time.Now().Format(time.RFC822)))
	supportLog.WriteString(fmt.Sprintf("Input tree  : %s\n", supportIntree))
	supportLog.WriteString(fmt.Sprintf("Boot trees  : %s\n", supportBoottrees))
	supportLog.WriteString(fmt.Sprintf("Output tree : %s\n", supportOutFile))
	supportLog.WriteString(fmt.Sprintf("Theor norm  : %t\n", !parsimonyEmpirical))
	supportLog.WriteString(fmt.Sprintf("Seed        : %d\n", parsimonySeed))
	supportLog.WriteString(fmt.Sprintf("CPUs        : %d\n", rootCpus))
}

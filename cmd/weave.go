package cmd

import (
	"github.com/AmadlaOrg/weaver/weave"
	"github.com/spf13/cobra"
)

var WeaveCmd = &cobra.Command{
	Use:   "weave",
	Short: "A brief description of your command",
	//Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//var ()
		/*cmd.PersistentFlags().StringVarP(
		&collection,
		"collection",
		"c",
		envVarCollection,
		"Specify the collection name")*/

		templateService := weave.NewWeaveService()
		templateService.Do()
		//templateService.Weave()
	},
}

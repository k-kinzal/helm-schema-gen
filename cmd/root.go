package cmd

import (
	"fmt"
	"github.com/karuppiah7890/go-jsonschema-generator"
	"github.com/karuppiah7890/helm-schema-gen/cmd/helper"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "helm schema-gen <values-yaml-file>",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Helm plugin to generate json schema for values yaml",
	Long: `Helm plugin to generate json schema for values yaml

Examples:
  $ helm schema-gen values.yaml    # generate schema json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("pass one values yaml file")
		}
		if len(args) != 1 {
			return fmt.Errorf("schema can be generated only for one values yaml at once")
		}

		valuesFilePath := args[0]
		values := make(map[string]interface{})
		valuesFileData, err := ioutil.ReadFile(valuesFilePath)
		if err != nil {
			return fmt.Errorf("error when reading file '%s': %v", valuesFilePath, err)
		}
		root := yaml.Node{}
		if ok, _ := cmd.Flags().GetBool("yaml-comment"); ok {
			if err := helper.UnmarshalWithUncommentYAML(valuesFileData, &root); err != nil {
				return fmt.Errorf("error when unmarshaling file '%s': %v", valuesFilePath, err)
			}
		} else {
			if err = yaml.Unmarshal(valuesFileData, &root); err != nil {
				return fmt.Errorf("error when unmarshaling file '%s': %v", valuesFilePath, err)
			}
		}
		err = root.Decode(&values)
		if err != nil {
			fmt.Println(err)
		}
		s := &jsonschema.Document{}
		s.ReadDeep(&values)
		fmt.Println(s)

		return nil
	},
}

// Execute executes the root command
func Execute() {
	rootCmd.Flags().Bool("yaml-comment", false, "Generate a schema with YAML format comments as YAML when it is an empty array or map")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

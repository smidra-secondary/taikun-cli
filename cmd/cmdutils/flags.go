package cmdutils

import (
	"log"
	"strings"

	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/cmd/cmderr"
	"github.com/itera-io/taikun-cli/config"
	"github.com/itera-io/taikun-cli/utils/gmap"
	"github.com/itera-io/taikun-cli/utils/out/fields"
	"github.com/itera-io/taikungoclient/client/common"
	"github.com/spf13/cobra"
)

type FlagCompCoreFunc func(cmd *cobra.Command, args []string, toComplete string) []string

func MarkFlagRequired(cmd *cobra.Command, flag string) {
	if err := cmd.MarkFlagRequired(flag); err != nil {
		log.Fatal(err)
	}
}

func RegisterFlagCompletionFunc(cmd *cobra.Command, flagName string, f FlagCompCoreFunc) {
	if err := cmd.RegisterFlagCompletionFunc(flagName, makeFlagCompFunc(f)); err != nil {
		log.Fatal(err)
	}
}

func RegisterFlagCompletion(cmd *cobra.Command, flagName string, values ...string) {
	RegisterFlagCompletionFunc(cmd, flagName, func(cmd *cobra.Command, args []string, toComplete string) []string {
		return values
	})
}

func makeFlagCompFunc(f FlagCompCoreFunc) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return f(cmd, args, toComplete), cobra.ShellCompDirectiveNoFileComp
	}
}

func AddSortByAndReverseFlags(cmd *cobra.Command, sortType string, fields fields.Fields) {
	cmd.Flags().StringVarP(
		&config.SortBy,
		"sort-by",
		"S",
		"",
		"Sort results by attribute value",
	)

	fieldNames := fields.AllNames()
	lowerStringSlice(fieldNames)
	RegisterFlagCompletionFunc(cmd, "sort-by", func(cmd *cobra.Command, args []string, toComplete string) []string {
		sortingElements, err := getSortingElements(sortType)
		if err != nil {
			return []string{}
		}

		completions := make([]string, 0)
		for _, jsonTag := range sortingElements {
			for _, field := range fields.AllFields() {
				if field.JsonTag() == jsonTag {
					completions = append(completions, field.Name())
					break
				}
			}
		}

		lowerStringSlice(completions)

		return completions
	},
	)

	cmd.Flags().BoolVarP(
		&config.ReverseSortDirection,
		"reverse",
		"R",
		false,
		"Reverse order of results when passed with the --sort-by flag",
	)
}

func getSortingElements(sortType string) (sortingElements []string, err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	params := common.NewCommonGetSortingElementsParams().WithV(api.Version)
	params = params.WithType(sortType)

	response, err := apiClient.Client.Common.CommonGetSortingElements(params, apiClient)
	if err != nil {
		return
	}

	sortingElements = response.Payload

	return
}

func AddOutputOnlyIDFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(
		&config.OutputOnlyID,
		"id-only",
		"I",
		false,
		"Output only the ID of the newly created resource (takes priority over the --format flag)",
	)
}

func AddColumnsFlag(cmd *cobra.Command, fields fields.Fields) {
	cmd.Flags().StringSliceVarP(
		&config.Columns,
		"columns",
		"C",
		[]string{},
		"Specify which columns to display in the output table",
	)
	columns := fields.AllNames()
	lowerStringSlice(columns)
	RegisterFlagCompletion(cmd, "columns", columns...)

	cmd.Flags().BoolVarP(
		&config.AllColumns,
		"all-columns",
		"A",
		false,
		"Display all columns in the output table (takes priority over the --columns flag)",
	)
}

func lowerStringSlice(stringSlice []string) {
	size := len(stringSlice)
	for i := 0; i < size; i++ {
		stringSlice[i] = strings.ToLower(stringSlice[i])
	}
}

func AddLimitFlag(cmd *cobra.Command) {
	cmd.Flags().Int32VarP(&config.Limit, "limit", "L", 0, "Limit number of results (limitless by default)")
	cmd.PreRunE = aggregateRunE(cmd.PreRunE,
		func(cmd *cobra.Command, args []string) error {
			if config.Limit < 0 {
				return cmderr.NegativeLimitFlagError
			}
			return nil
		},
	)

}

func CheckFlagValue(flagName string, flagValue string, valid gmap.GenericMap) error {
	if !valid.Contains(flagValue) {
		return cmderr.UnknownFlagValueError(flagName, flagValue, valid.Keys())
	}
	return nil
}

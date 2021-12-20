package all

import (
	"taikun-cli/api"
	"taikun-cli/apiconfig"
	"taikun-cli/cmd/cmderr"
	"taikun-cli/cmd/cmdutils"
	"taikun-cli/config"
	"taikun-cli/utils/format"
	"taikun-cli/utils/types"

	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

type AllOptions struct {
	CloudCredentialID    int32
	Limit                int32
	MaxCPU               int32
	MaxRAM               float64
	MinCPU               int32
	MinRAM               float64
	ReverseSortDirection bool
	SortBy               string
}

func NewCmdAll() *cobra.Command {
	var opts AllOptions

	cmd := &cobra.Command{
		Use:   "all <cloud-credential-id>",
		Short: "List all flavors by cloud credential",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cloudCredentialID, err := types.Atoi32(args[0])
			if err != nil {
				return cmderr.IDArgumentNotANumberError
			}
			if !config.OutputFormatIsValid() {
				return cmderr.OutputFormatInvalidError
			}
			opts.CloudCredentialID = cloudCredentialID
			return allRun(&opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.ReverseSortDirection, "reverse", "r", false, "Reverse order of results")
	cmd.Flags().Int32Var(&opts.MaxCPU, "max-cpu", 36, "Maximal CPU count")
	cmd.Flags().Float64Var(&opts.MaxRAM, "max-ram", 500, "Maximal RAM size in GiB")
	cmd.Flags().Int32Var(&opts.MinCPU, "min-cpu", 2, "Minimal CPU count")
	cmd.Flags().Float64Var(&opts.MinRAM, "min-ram", 2, "Minimal RAM size in GiB")
	cmd.Flags().Int32VarP(&opts.Limit, "limit", "l", 0, "Limit number of results")

	cmdutils.AddSortByFlag(cmd, &opts.SortBy, models.FlavorsListDto{})

	return cmd
}

func printResults(flavors []*models.FlavorsListDto) {
	if config.OutputFormat == config.OutputFormatJson {
		format.PrettyPrintJson(flavors)
	} else if config.OutputFormat == config.OutputFormatTable {
		data := make([]interface{}, len(flavors))
		for i, flavor := range flavors {
			data[i] = flavor
		}
		format.PrettyPrintTable(data,
			"name",
			"cpu",
			"ram",
		)
	}
}

func allRun(opts *AllOptions) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	params := cloud_credentials.NewCloudCredentialsAllFlavorsParams().WithV(apiconfig.Version)
	params = params.WithCloudID(opts.CloudCredentialID)
	params = params.WithStartCPU(&opts.MinCPU).WithEndCPU(&opts.MaxCPU)
	minRAM := types.GiBToMiB(opts.MinRAM)
	maxRAM := types.GiBToMiB(opts.MaxRAM)
	params = params.WithStartRAM(&minRAM).WithEndRAM(&maxRAM)
	if opts.ReverseSortDirection {
		apiconfig.ReverseSortDirection()
	}
	if opts.SortBy != "" {
		params = params.WithSortBy(&opts.SortBy).WithSortDirection(&apiconfig.SortDirection)
	}

	flavors := []*models.FlavorsListDto{}
	for {
		response, err := apiClient.Client.CloudCredentials.CloudCredentialsAllFlavors(params, apiClient)
		if err != nil {
			return err
		}
		flavors = append(flavors, response.Payload.Data...)
		flavorsCount := int32(len(flavors))
		if opts.Limit != 0 && flavorsCount >= opts.Limit {
			break
		}
		if flavorsCount == response.Payload.TotalCount {
			break
		}
		params = params.WithOffset(&flavorsCount)
	}

	if opts.Limit != 0 && int32(len(flavors)) > opts.Limit {
		flavors = flavors[:opts.Limit]
	}

	printResults(flavors)
	return
}

package flavors

import (
	"context"
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/cmd/cloudcredential/utils"
	"github.com/itera-io/taikun-cli/cmd/cmderr"
	"github.com/itera-io/taikun-cli/cmd/cmdutils"
	"github.com/itera-io/taikun-cli/config"
	"github.com/itera-io/taikun-cli/utils/out"
	"github.com/itera-io/taikun-cli/utils/out/field"
	"github.com/itera-io/taikun-cli/utils/out/fields"
	"github.com/itera-io/taikun-cli/utils/types"
	tk "github.com/itera-io/taikungoclient"
	taikuncore "github.com/itera-io/taikungoclient/client"
	"github.com/spf13/cobra"
)

var flavorsFields = fields.New(
	[]*field.Field{
		field.NewVisible(
			"NAME", "name",
		),
		field.NewVisible(
			"CPU", "cpu",
		),
		field.NewVisibleWithToStringFunc(
			"RAM", "ram", out.FormatRAM,
		),
		field.NewHidden(
			"DESCRIPTION", "description",
		),
	},
)

type FlavorsOptions struct {
	CloudCredentialID int32
	MaxCPU            int32
	MaxRAM            float64
	MinCPU            int32
	MinRAM            float64
	Limit             int32
}

func NewCmdFlavors() *cobra.Command {
	var opts FlavorsOptions

	cmd := cobra.Command{
		Use:   "flavors <cloud-credential-id>",
		Short: "List a cloud credential's flavors",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cloudCredentialID, err := types.Atoi32(args[0])
			if err != nil {
				return cmderr.ErrIDArgumentNotANumber
			}
			opts.CloudCredentialID = cloudCredentialID
			if err = adjustRamUnits(&opts); err != nil {
				return err
			}
			return flavorRun(&opts)
		},
	}

	cmd.Flags().Int32Var(&opts.MaxCPU, "max-cpu", 36, "Maximal CPU count")
	cmd.Flags().Float64Var(&opts.MaxRAM, "max-ram", 500, "Maximal RAM size in GiB")
	cmd.Flags().Int32Var(&opts.MinCPU, "min-cpu", 2, "Minimal CPU count")
	cmd.Flags().Float64Var(&opts.MinRAM, "min-ram", 2, "Minimal RAM size in GiB")

	cmdutils.AddLimitFlag(&cmd, &opts.Limit)
	cmdutils.AddSortByAndReverseFlags(&cmd, "flavors", flavorsFields)
	cmdutils.AddColumnsFlag(&cmd, flavorsFields)

	return &cmd
}

func adjustRamUnits(opts *FlavorsOptions) (err error) {
	cloudType, err := utils.GetCloudType(opts.CloudCredentialID)
	if err != nil {
		return
	}

	switch cloudType {
	case utils.GOOGLE:
		// Temporarily ignore RAM range for Google until units are set to GiB
		opts.MinRAM = -1
		opts.MaxRAM = -1
	default:
		opts.MinRAM = types.GiBToMiB(opts.MinRAM)
		opts.MaxRAM = types.GiBToMiB(opts.MaxRAM)
	}

	return
}

func flavorRun(opts *FlavorsOptions) (err error) {
	myApiClient := tk.NewClient()
	myRequest := myApiClient.Client.CloudCredentialAPI.CloudcredentialsAllFlavors(context.TODO(), opts.CloudCredentialID)
	myRequest = myRequest.StartCpu(opts.MinCPU).EndCpu(opts.MaxCPU)
	if config.SortBy != "" {
		myRequest = myRequest.SortBy(config.SortBy).SortDirection(*api.GetSortDirection())
	}
	minRAM := opts.MinRAM
	maxRAM := opts.MaxRAM

	// Temporarily ignore RAM range for Google until units are set to GiB
	if minRAM != -1 && maxRAM != -1 {
		myRequest = myRequest.StartRam(minRAM).EndRam(maxRAM)
	}

	var flavors = make([]taikuncore.FlavorsListDto, 0)
	for {
		data, response, err := myRequest.Execute()
		if err != nil {
			return tk.CreateError(response, err)
		}

		flavors = append(flavors, data.GetData()...)

		flavorsCount := int32(len(flavors))
		if opts.Limit != 0 && flavorsCount >= opts.Limit {
			break
		}

		if flavorsCount == data.GetTotalCount() {
			break
		}

		myRequest = myRequest.Offset(flavorsCount)
	}

	if opts.Limit != 0 && int32(len(flavors)) > opts.Limit {
		flavors = flavors[:opts.Limit]
	}

	return out.PrintResults(flavors, flavorsFields)

}

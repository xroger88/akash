package query

import (
	"github.com/ovrclk/akash/cmd/akash/context"
	"github.com/ovrclk/akash/state"
	"github.com/ovrclk/akash/types"
	"github.com/spf13/cobra"
)

func queryAccountCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "account",
		Short: "query account",
		Args:  cobra.ExactArgs(1),
		RunE:  context.WithContext(context.RequireNode(doQueryAccountCommand)),
	}

	return cmd
}

func doQueryAccountCommand(ctx context.Context, cmd *cobra.Command, args []string) error {
	structure := new(types.Account)
	account := args[0]
	path := state.AccountPath + account
	return doQuery(ctx, path, structure)
}

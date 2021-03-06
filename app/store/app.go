package store

import (
	"github.com/tendermint/abci/types"
	tmtypes "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/log"

	apptypes "github.com/ovrclk/akash/app/types"
	"github.com/ovrclk/akash/state"
	"github.com/ovrclk/akash/types/code"
)

const (
	QueryPath = "/store"
	Name      = "store"
)

type app struct {
	*apptypes.BaseApp
}

func NewApp(state state.State, logger log.Logger) (apptypes.Application, error) {
	return &app{apptypes.NewBaseApp(Name, state, logger)}, nil
}

func (a *app) AcceptQuery(req tmtypes.RequestQuery) bool {
	return req.Path == QueryPath
}

func (a *app) Query(req types.RequestQuery) types.ResponseQuery {
	if !a.AcceptQuery(req) {
		return types.ResponseQuery{
			Code: code.ERROR,
			Log:  "invalid query",
		}
	}

	db := a.State().DB()

	if req.Prove {
		val, proof, err := db.GetWithProof(req.Data)
		if err != nil {
			return types.ResponseQuery{
				Code: code.ERROR,
				Log:  err.Error(),
			}
		}
		return types.ResponseQuery{
			Value:  val,
			Height: int64(db.Version()),
			Proof:  proof.Bytes(),
		}
	}

	val := db.Get(req.Data)
	return types.ResponseQuery{
		Value:  val,
		Height: int64(db.Version()),
	}
}

func (a *app) AcceptTx(ctx apptypes.Context, tx interface{}) bool {
	return false
}

func (a *app) CheckTx(ctx apptypes.Context, tx interface{}) types.ResponseCheckTx {
	return types.ResponseCheckTx{
		Code: code.UNKNOWN_TRANSACTION,
		Log:  "store app: unknown transaction",
	}
}

func (a *app) DeliverTx(ctx apptypes.Context, tx interface{}) types.ResponseDeliverTx {
	return types.ResponseDeliverTx{
		Code: code.UNKNOWN_TRANSACTION,
		Log:  "store app: unknown transaction",
	}
}

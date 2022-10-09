package sdk

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

// CreateSimulateEnv new context of simulated
func CreateSimulateEnv() context.Context {
	return context.WithValue(context.Background(), types.SimulatedEnvkey, simulated.NewSimulated())
}

// CreateEntityEnv new context of entity
func CreateEntityEnv() context.Context {
	return context.WithValue(context.Background(), types.SimulatedEnvkey, "")
}

func simulatedEnv(ctx context.Context) (simulated.Fsm, bool) {
	env, ok := ctx.Value(types.SimulatedEnvkey).(simulated.Fsm) //nolint:govet
	return env, ok                                              //nolint:govet
}

// SetActorAndAddress set actor
func SetActorAndAddress(ctx context.Context, actorID uint32, actorState migration.Actor, addr address.Address) {
	if env, ok := simulatedEnv(ctx); ok {
		env.SetActorAndAddress(actorID, actorState, addr)
	}
}

// SetAccount set account
func SetAccount(ctx context.Context, actorID uint32, addr address.Address, actor migration.Actor) {
	if env, ok := simulatedEnv(ctx); ok {
		env.SetAccount(actorID, addr, actor)
	}

}

// SetBaseFee set BaseFee
func SetBaseFee(ctx context.Context, ta big.Int) {
	if env, ok := simulatedEnv(ctx); ok {
		env.SetBaseFee(ta)
	}

}

// SetSend set send mock
func SetSend(ctx context.Context, mock ...simulated.SendMock) {
	if env, ok := simulatedEnv(ctx); ok {
		env.SetSend(mock...)
	}

}

// SetTotalFilCircSupply set FilCircSupply
func SetTotalFilCircSupply(ctx context.Context, mock ...simulated.SendMock) {
	if env, ok := simulatedEnv(ctx); ok {
		env.SetSend(mock...)
	}

}

// SetCurrentBalance set CurrentBalance
func SetCurrentBalance(ctx context.Context, ta big.Int) {
	if env, ok := simulatedEnv(ctx); ok {
		env.SetCurrentBalance(ta)
	}

}

// SetCallContext set context
func SetCallContext(ctx context.Context, callcontext *types.InvocationContext) {
	if env, ok := simulatedEnv(ctx); ok {
		env.SetCallContext(callcontext)
	}

}

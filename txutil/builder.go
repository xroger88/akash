package txutil

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/ovrclk/akash/types"
	crypto "github.com/tendermint/go-crypto"
)

type TxBuilder interface {
	SignableTx
	Signature() crypto.Signature
	TxBytes() ([]byte, error)
}

func BuildTx(signer Signer, nonce uint64, payload interface{}) ([]byte, error) {
	txb, err := NewTxBuilder(nonce, payload)
	if err != nil {
		return nil, err
	}
	if err := signer.Sign(txb); err != nil {
		return nil, err
	}
	return txb.TxBytes()
}

func NewTxBuilder(nonce uint64, payload interface{}) (TxBuilder, error) {
	tx := &types.Tx{}

	switch payload := payload.(type) {
	case *types.TxSend:
		tx.Payload.Payload = &types.TxPayload_TxSend{TxSend: payload}
	case *types.TxCreateDeployment:
		tx.Payload.Payload = &types.TxPayload_TxCreateDeployment{TxCreateDeployment: payload}
	case *types.TxCreateProvider:
		tx.Payload.Payload = &types.TxPayload_TxCreateProvider{TxCreateProvider: payload}
	case *types.TxCreateOrder:
		tx.Payload.Payload = &types.TxPayload_TxCreateOrder{TxCreateOrder: payload}
	case *types.TxCreateFulfillment:
		tx.Payload.Payload = &types.TxPayload_TxCreateFulfillment{TxCreateFulfillment: payload}
	case *types.TxCreateLease:
		tx.Payload.Payload = &types.TxPayload_TxCreateLease{TxCreateLease: payload}
	case *types.TxCloseDeployment:
		tx.Payload.Payload = &types.TxPayload_TxCloseDeployment{TxCloseDeployment: payload}
	case *types.TxDeploymentClosed:
		tx.Payload.Payload = &types.TxPayload_TxDeploymentClosed{TxDeploymentClosed: payload}
	default:
		return nil, fmt.Errorf("unknown payload type: %T", payload)
	}

	tx.Payload.Nonce = nonce

	pbytes, err := proto.Marshal(&tx.Payload)
	if err != nil {
		return nil, err
	}
	return &txBuilder{tx, pbytes}, nil
}

type txBuilder struct {
	tx     *types.Tx
	pbytes []byte
}

func (b *txBuilder) SignBytes() []byte {
	return b.pbytes
}

func (b *txBuilder) Sign(key crypto.PubKey, sig crypto.Signature) error {
	if b.tx.Key != nil || b.tx.Signature != nil {
		return fmt.Errorf("already signed")
	}
	b.tx.Key = key.Bytes()
	b.tx.Signature = sig.Bytes()
	return nil
}

func (b *txBuilder) Signers() ([]crypto.PubKey, error) {
	if b.tx.Key == nil {
		return nil, nil
	}

	key, err := crypto.PubKeyFromBytes(b.tx.Key)
	if err != nil {
		return nil, err
	}

	return []crypto.PubKey{key}, nil
}

func (b *txBuilder) Signature() crypto.Signature {
	sig, _ := crypto.SignatureFromBytes(b.tx.Signature)
	return sig
}

func (b *txBuilder) TxBytes() ([]byte, error) {
	return proto.Marshal(b.tx)
}

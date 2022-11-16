package sybilprotection

import (
	"context"
	"time"

	"github.com/iotaledger/hive.go/core/crypto/ed25519"
	"github.com/iotaledger/hive.go/core/identity"
	"github.com/iotaledger/hive.go/core/serix"
	"github.com/iotaledger/hive.go/core/types"

	"github.com/iotaledger/goshimmer/packages/core/commitment"
)

type Attestation struct {
	IssuerID         identity.ID       `serix:"0"`
	IssuingTime      time.Time         `serix:"1"`
	CommitmentID     commitment.ID     `serix:"2"`
	BlockContentHash types.Identifier  `serix:"3"`
	Signature        ed25519.Signature `serix:"4"`
}

func NewAttestation(issuerID identity.ID, issuingTime time.Time, commitmentID commitment.ID, blockContentHash types.Identifier, signature ed25519.Signature) *Attestation {
	return &Attestation{
		IssuerID:         issuerID,
		IssuingTime:      issuingTime,
		CommitmentID:     commitmentID,
		BlockContentHash: blockContentHash,
		Signature:        signature,
	}
}

func (a Attestation) Bytes() (bytes []byte, err error) {
	return serix.DefaultAPI.Encode(context.Background(), a, serix.WithValidation())
}

func (a *Attestation) FromBytes(bytes []byte) (consumedBytes int, err error) {
	return serix.DefaultAPI.Decode(context.Background(), bytes, a, serix.WithValidation())
}

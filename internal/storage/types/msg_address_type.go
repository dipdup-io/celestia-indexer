package types

// swagger:enum MsgAddressType
/*
	ENUM(
		validator,
		delegator,
		depositor,

		validatorSrcAddress,
		validatorDstAddress,
		fromAddress,
		toAddress,
		grantee,
		granter,
		signer,
		withdraw,

		voter,
		proposer,
	)
*/
//go:generate go-enum --marshal --sql --values
type MsgAddressType string

package types

// swagger:enum MsgType
/*
	ENUM(
		MsgUnknown,
		MsgWithdrawValidatorCommission,
		MsgWithdrawDelegatorReward,
		MsgEditValidator,
		MsgBeginRedelegate,
		MsgCreateValidator,
		MsgDelegate,
		MsgUndelegate,
		MsgUnjail,
		MsgSend,
		MsgCreateVestingAccount,
		MsgCreatePeriodicVestingAccount,
		MsgPayForBlobs,
		MsgGrantAllowance,
		MsgRegisterEVMAddress,
		MsgSetWithdrawAddress,

		MsgVote,
		MsgVoteWeighted,
		MsgSubmitProposal,
	)
*/
//go:generate go-enum --marshal --sql --values --noprefix
type MsgType string

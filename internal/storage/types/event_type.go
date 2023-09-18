package types

// swagger:enum EventType
/*
	ENUM(
		unknown,
		coin_received,
		coinbase,
		coin_spent,
		burn,
		mint,
		message,
		proposer_reward,
		rewards,
		commission,
		liveness,
		transfer,
		celestia.blob.v1.EventPayForBlobs,
		redelegate,
		AttestationRequest,
		withdraw_rewards,
		withdraw_commission,
		set_withdraw_address,
		create_validator,
		delegate,
		edit_validator,
		unbond,
		tx,
		use_feegrant,
		revoke_feegrant,
		set_feegrant,
		update_feegrant,
		slash,

		proposal_vote,
		proposal_deposit
	)
*/
//go:generate go-enum --marshal --sql --values
type EventType string

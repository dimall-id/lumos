package wallet_trans_status

const (
	TOPUP_CREATED = 1100
	PAYOUT_CREATED = 1200

	TOPUP_SUCCEED = 2100
	PAYOUT_SUCCEED = 2200

	TOPUP_FAILED = 3100
	TOPUP_PAYMENT_FAILED = 3101
	PAYOUT_FAILED = 3200
	PAYOUT_TIMEOUT = 3201
	PAYOUT_USER_DECLINE = 3202
	PAYOUT_USER_NOT_FOUND = 3203
	PAYOUT_INSUFFICIENT_BALANCE = 3204
	PAYOUT_MERCHANT_CREDENTIAL = 3205

)
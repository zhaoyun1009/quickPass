package e

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	LOGIN_FAILED = 403

	ERROR_EXIST_TAG       = 10001
	ERROR_EXIST_TAG_FAIL  = 10002
	ERROR_NOT_EXIST_TAG   = 10003
	ERROR_GET_TAGS_FAIL   = 10004
	ERROR_COUNT_TAG_FAIL  = 10005
	ERROR_ADD_TAG_FAIL    = 10006
	ERROR_EDIT_TAG_FAIL   = 10007
	ERROR_DELETE_TAG_FAIL = 10008
	ERROR_EXPORT_TAG_FAIL = 10009
	ERROR_IMPORT_TAG_FAIL = 10010

	ERROR_NOT_EXIST_ACCEPTOR        = 10011
	ERROR_CHECK_EXIST_ACCEPTOR_FAIL = 10012
	ERROR_ADD_ACCEPTOR_FAIL         = 10013
	ERROR_DELETE_ACCEPTOR_FAIL      = 10014
	ERROR_EDIT_ACCEPTOR_FAIL        = 10015
	ERROR_COUNT_ACCEPTOR_FAIL       = 10016
	ERROR_GET_ACCEPTORS_FAIL        = 10017
	ERROR_GET_ACCEPTOR_FAIL         = 10018
	ERROR_GEN_ACCEPTOR_POSTER_FAIL  = 10019

	ERROR_CHECK_EXIST_MERCHANT_FAIL = 10020
	ERROR_NOT_EXIST_MERCHANT        = 10021
	ERROR_ADD_MERCHANT_FAIL         = 10022
	ERROR_DELETE_MERCHANT_FAIL      = 10023
	ERROR_EDIT_MERCHANT_FAIL        = 10024
	ERROR_COUNT_MERCHANT_FAIL       = 10025
	ERROR_GET_MERCHANTS_FAIL        = 10026
	ERROR_GET_MERCHANT_FAIL         = 10027

	ERROR_CHECK_EXIST_MERCHANT_RATE_FAIL = 10200
	ERROR_NOT_EXIST_MERCHANT_RATE        = 10201
	ERROR_ADD_MERCHANT_RATE_FAIL         = 10202
	ERROR_DELETE_MERCHANT_RATE_FAIL      = 10203
	ERROR_EDIT_MERCHANT_RATE_FAIL        = 10204
	ERROR_COUNT_MERCHANT_RATE_FAIL       = 10205
	ERROR_GET_MERCHANT_RATES_FAIL        = 10206
	ERROR_GET_MERCHANT_RATE_FAIL         = 10207

	ERROR_NOT_EXIST_USER       = 10030
	ERROR_ADD_USER_FAIL        = 10031
	ERROR_DELETE_USER_FAIL     = 10032
	ERROR_EDIT_USER_FAIL       = 10033
	ERROR_COUNT_USER_FAIL      = 10034
	ERROR_GET_USERS_FAIL       = 10035
	ERROR_GET_USER_FAIL        = 10036
	ERROR_GEN_USER_POSTER_FAIL = 10037
	ERROR_UPDATE_USER_FAIL     = 10038

	ERROR_VALIDATE_PASSWORD       = 10040
	ERROR_UPDATE_PASSWORD         = 10041
	ERROR_VALIDATE_TRADE_PASSWORD = 10042
	ERROR_UPDATE_TRADE_PASSWORD   = 10043

	ERROR_COUNT_BILL_FAIL = 10050
	ERROR_GET_BILLS_FAIL  = 10051

	ERROR_NOT_EXIST_CHANNEL        = 10061
	ERROR_CHECK_EXIST_CHANNEL_FAIL = 10062
	ERROR_ADD_CHANNEL_FAIL         = 10063
	ERROR_DELETE_CHANNEL_FAIL      = 10064
	ERROR_EDIT_CHANNEL_FAIL        = 10065
	ERROR_COUNT_CHANNEL_FAIL       = 10066
	ERROR_GET_CHANNELS_FAIL        = 10067
	ERROR_GET_CHANNEL_FAIL         = 10068

	ERROR_NOT_EXIST_ACCEPTOR_CARD        = 10071
	ERROR_CHECK_EXIST_ACCEPTOR_CARD_FAIL = 10072
	ERROR_ADD_ACCEPTOR_CARD_FAIL         = 10073
	ERROR_DELETE_ACCEPTOR_CARD_FAIL      = 10074
	ERROR_EDIT_ACCEPTOR_CARD_FAIL        = 10075
	ERROR_COUNT_ACCEPTOR_CARD_FAIL       = 10076
	ERROR_GET_ACCEPTOR_CARDS_FAIL        = 10077
	ERROR_GET_ACCEPTOR_CARD_FAIL         = 10078
	ERROR_GEN_ACCEPTOR_CARD_POSTER_FAIL  = 10079

	ERROR_AGAIN_CALLBACK         = 10080
	ERROR_NOT_EXIST_ORDER        = 10081
	ERROR_CHECK_EXIST_ORDER_FAIL = 10082
	ERROR_ADD_ORDER_FAIL         = 10083
	ERROR_DELETE_ORDER_FAIL      = 10084
	ERROR_EDIT_ORDER_FAIL        = 10085
	ERROR_COUNT_ORDER_FAIL       = 10086
	ERROR_GET_ORDERS_FAIL        = 10087
	ERROR_GET_ORDER_FAIL         = 10088
	ERROR_HAS_UNFINISH_ORDER     = 10089

	ERROR_MATCH_FAILED                 = 10090
	ERROR_NOT_EXIST_MATCH_CACHE        = 10091
	ERROR_CHECK_EXIST_MATCH_CACHE_FAIL = 10092
	ERROR_ADD_MATCH_CACHE_FAIL         = 10093
	ERROR_DELETE_MATCH_CACHE_FAIL      = 10094
	ERROR_EDIT_MATCH_CACHE_FAIL        = 10095
	ERROR_COUNT_MATCH_CACHE_FAIL       = 10096
	ERROR_GET_MATCH_CACHES_FAIL        = 10097
	ERROR_GET_MATCH_CACHE_FAIL         = 10098

	ERROR_NOT_EXIST_FUND        = 10101
	ERROR_CHECK_EXIST_FUND_FAIL = 10102
	ERROR_ADD_FUND_FAIL         = 10103
	ERROR_DELETE_FUND_FAIL      = 10104
	ERROR_EDIT_FUND_FAIL        = 10105
	ERROR_COUNT_FUND_FAIL       = 10106
	ERROR_GET_FUNDS_FAIL        = 10107
	ERROR_GET_FUND_FAIL         = 10108

	ERROR_NOT_EXIST_MERCHANT_CARD        = 10111
	ERROR_CHECK_EXIST_MERCHANT_CARD_FAIL = 10112
	ERROR_ADD_MERCHANT_CARD_FAIL         = 10113
	ERROR_DELETE_MERCHANT_CARD_FAIL      = 10114
	ERROR_EDIT_MERCHANT_CARD_FAIL        = 10115
	ERROR_COUNT_MERCHANT_CARD_FAIL       = 10116
	ERROR_GET_MERCHANT_CARDS_FAIL        = 10117
	ERROR_GET_MERCHANT_CARD_FAIL         = 10118

	ERROR_NOT_EXIST_ABNORMAL_ORDER        = 10121
	ERROR_CHECK_EXIST_ABNORMAL_ORDER_FAIL = 10122
	ERROR_ADD_ABNORMAL_ORDER_FAIL         = 10123
	ERROR_DELETE_ABNORMAL_ORDER_FAIL      = 10124
	ERROR_EDIT_ABNORMAL_ORDER_FAIL        = 10125
	ERROR_COUNT_ABNORMAL_ORDER_FAIL       = 10126
	ERROR_GET_ABNORMAL_ORDERS_FAIL        = 10127
	ERROR_GET_ABNORMAL_ORDER_FAIL         = 10128

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004

	ERROR_UPLOAD_SAVE_IMAGE_FAIL    = 30001
	ERROR_UPLOAD_CHECK_IMAGE_FAIL   = 30002
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT = 30003

	ERROR_IP_RISK = 30010

	Insufficient            = 40001
	AcceptorStatusClosed    = 40002
	ExistUsername           = 40003
	CheckRsaError           = 40004
	MerchantBuyStatusClosed = 40005
)

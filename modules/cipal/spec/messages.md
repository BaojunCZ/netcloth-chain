
# Messages

## MsgCIPAL

```cassandraql
type ServiceInfo struct {
	Type    uint64 `json:"type" yaml:"type"`
	Address string `json:"address" yaml:"address"`
}

type Param struct {
	UserAddress string      `json:"user_address" yaml:"user_address"`
	ServiceInfo ServiceInfo `json:"service_info" yaml:"service_info"`
	Expiration  time.Time   `json:"expiration"`
}

type UserRequest struct {
	Params Param             `json:"params" yaml:"params"`
	Sig    auth.StdSignature `json:"signature" yaml:"signature`
}

type MsgCIPAL struct {
	From        sdk.AccAddress `json:"from" yaml:"from`
	UserRequest UserRequest    `json:"user_request" yaml:"user_request"`
}
```

examples:
```cassandraql
 {
 	"from": "nch1sdh9efnf2tjcatcytcrllexsrmwze4acwh8ulr",
 	"user_request": {
 		"params": {
 			"user_address": "nch1kqjc3gptzzujnk7aqa6chxa59uljh0gla0fz6u",
 			"service_info": {
 				"type": "1",
 				"address": "nch1sdh9efnf2tjcatcytcrllexsrmwze4acwh8ulr"
 			},
 			"type": 1,
 			"memo": "",
 			"expiration": "2020-01-13T03:43:44Z"
 		},
 		"signature": {
 			"pub_key": {
 				"type": "tendermint/PubKeySecp256k1",
 				"value": "Anj0Kfp/d52ilVb71IEiyryfcjQIj4/T//SAGKO2bS9S"
 			},
 			"signature": "b215XySUaHA33ej0xsjy0917Bfo+0/RXvF5p7oWUE8BII4/PWYoRAD34Sny5opdHpGu8F0fPCUN8I7F+O+OA3w=="
 		}
 	}
 }

```

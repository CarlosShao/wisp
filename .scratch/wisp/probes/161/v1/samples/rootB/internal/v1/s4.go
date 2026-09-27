package v1

var store = map[string]string{"apiSecret": "PROBE161V1NOTACREDENT0002"}

func patch(cache map[string]string) {
	cache["apiKey"] = "PROBE161V1NOTACREDENT0001"
}

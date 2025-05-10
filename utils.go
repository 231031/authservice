package authservice

import "encoding/base64"

func genBasicAuthHeader(clientID, clientSecret string) string {
	auth := clientID + ":" + clientSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

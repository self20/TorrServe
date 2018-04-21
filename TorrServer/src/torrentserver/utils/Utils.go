package utils

import (
	"encoding/base32"
	"math/rand"
	"regexp"
)

func FileToLink(file string) string {
	re := regexp.MustCompile(`[ !\*'\(\);:@&=\+\$,/\?#\[\]~",]`)
	return re.ReplaceAllString(file, `_`)
}

func PeerIDRandom(peer string) string {
	return peer + getToken(20-len(peer))
}

func getToken(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

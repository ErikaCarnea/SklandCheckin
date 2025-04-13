package client

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type SignHeader struct {
	Platform  string `json:"platform"`
	Timestamp string `json:"timestamp"`
	DId       string `json:"dId"`
	VName     string `json:"vName"`
}

func (c *HttpClient) generateSignature(path, bodyOrQuery string) (string, map[string]string, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix()-2)
	header := SignHeader{
		Timestamp: timestamp,
	}

	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", nil, err
	}
	sStr := fmt.Sprintf("%s%s%s%s", path, bodyOrQuery, timestamp, string(headerJson))

	h := hmac.New(sha256.New, []byte(c.signToken))
	h.Write([]byte(sStr))
	sha := hex.EncodeToString(h.Sum(nil))

	md5Hash := md5.Sum([]byte(sha))
	sign := hex.EncodeToString(md5Hash[:])

	headerCa := map[string]string{
		"timestamp": header.Timestamp,
	}

	return sign, headerCa, nil
}

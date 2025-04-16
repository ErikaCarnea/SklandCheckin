package api

import (
	"Skland/client"
	"Skland/models"
	"fmt"
	"net/http"
)

type AttendanceAPI struct {
	client client.HTTPClient
}

func NewAttendanceAPI(c client.HTTPClient) *AttendanceAPI {
	return &AttendanceAPI{client: c}
}

func (a *AttendanceAPI) SignAttendance(uid, gameId string) (*models.AttendanceResult, error) {
	reqBody := models.AttendanceRequest{Uid: uid, GameId: gameId}
	var result models.AttendanceResult

	if err := a.client.ExecuteRequest(
		http.MethodPost,
		"https://zonai.skland.com/api/v1/game/attendance",
		reqBody,
		&result,
	); err != nil {
		return nil, fmt.Errorf("签到失败: %w", err)
	}

	return &result, nil

	//urlStr := "https://zonai.skland.com/api/v1/game/attendance"
	//reqBody := models.AttendanceRequest{
	//	Uid:    uid,
	//	GameId: gameId,
	//}
	//
	//headers, err := a.client.GetSignHeaders(urlStr, http.MethodPost, reqBody)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get sign header: %v", err)
	//}
	//headers["Content-Type"] = "application/json"
	//
	//resp, err := a.client.DoRequest(http.MethodPost, urlStr, reqBody, headers)
	//if err != nil {
	//	return nil, err
	//}
	//defer func(Body io.ReadCloser) {
	//	if err := Body.Close(); err != nil {
	//		log.Printf("关闭响应体失败: %v", err)
	//	}
	//}(resp.Body)
	//
	//body, err := client.ReadResponseBody(resp)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var result models.AttendanceResult
	//if err := json.Unmarshal(body, &result); err != nil {
	//	return nil, err
	//}
	//
	//if result.Code != 0 {
	//	return &result, fmt.Errorf("API error: %s (code: %d)", result.Message, result.Code)
	//}
	//
	//return &result, nil
}

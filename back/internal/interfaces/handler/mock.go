package handler

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
)

// モックが返す固定データ．固定値は poll.go / user.go にもある．
// today() だけは本実装でも要るので，消すときに巻き込まないこと．

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

const mockToken = "9f2c1d0b8a7e6f5d4c3b2a1908172635"

var mockUsers = []openapi.User{
	{Id: 1, Name: "nishi", CreatedAt: mockCreatedAt(0)},
	{Id: 2, Name: "sato", CreatedAt: mockCreatedAt(1)},
	{Id: 3, Name: "tanaka", CreatedAt: mockCreatedAt(2)},
}

// mockCreatedAt は決め打ちの作成日時を返す．実行のたびに値が変わると
// クライアントのスナップショットテストが不安定になるため固定している．
func mockCreatedAt(offsetDays int) time.Time {
	return time.Date(2026, 9, 19, 12, 0, 0, 0, jst).AddDate(0, 0, offsetDays)
}

// today はJSTの今日を返す．
// time.Truncate はUTC基準で丸めるため，JSTの0:00〜8:59に前日を返してしまう．
func today() openapi_types.Date {
	now := time.Now().In(jst)
	return openapi_types.Date{Time: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)}
}

func mockUserByID(id int64) (openapi.User, bool) {
	for _, u := range mockUsers {
		if u.Id == id {
			return u, true
		}
	}
	return openapi.User{}, false
}

// mockAnswers は3状態が1件ずつ入った回答一覧を返す．クライアントが
// 「行ける / 行けない / 未回答」の表示を一度に確認できるようにしている．
//
// 日付は Answer ではなく囲む側（Answers / MyAnswer）が持つ．
func mockAnswers() []openapi.Answer {
	return []openapi.Answer{
		{
			UserId:    mockUsers[0].Id,
			UserName:  mockUsers[0].Name,
			Status:    openapi.Available,
			TimeSlots: openapi.TimeSlots{"18:00", "18:30", "19:00"},
		},
		{
			UserId:    mockUsers[1].Id,
			UserName:  mockUsers[1].Name,
			Status:    openapi.Unavailable,
			TimeSlots: openapi.TimeSlots{},
		},
		{
			UserId:    mockUsers[2].Id,
			UserName:  mockUsers[2].Name,
			Status:    openapi.Undecided,
			TimeSlots: openapi.TimeSlots{},
		},
	}
}

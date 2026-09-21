package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
	"github.com/nishi240masa/meshi-iko/back/internal/infrastructure/database"
	"github.com/nishi240masa/meshi-iko/back/internal/interfaces/router"
	"github.com/nishi240masa/meshi-iko/back/internal/registry"
)

// APIはまだモックなので，仕様どおりの形とステータスコードだけを検証する．
// ロジックを実装したら，値そのものを検証するテストに書き換える．

const (
	testAdminToken  = "test-admin-token"
	testBearerToken = "any-token-works-in-mock"
)

// newTestServer はインメモリDBを使い，本番と同じ組み立てでサーバを作る．
func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := database.New(":memory:", false)
	if err != nil {
		t.Fatalf("DBの初期化に失敗: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := db.Migrate(); err != nil {
		t.Fatalf("マイグレーションに失敗: %v", err)
	}

	return router.New(registry.New(db).NewHandler(), router.Config{
		CORSAllowedOrigin: "*",
		AdminToken:        testAdminToken,
	})
}

// request は BasePath 配下にリクエストを送る．
func request(t *testing.T, srv http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("リクエストのエンコードに失敗: %v", err)
		}
		reader = bytes.NewReader(encoded)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, router.BasePath+path, reader)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

func asUser() map[string]string {
	return map[string]string{"Authorization": "Bearer " + testBearerToken}
}

func asAdmin() map[string]string {
	return map[string]string{"X-Admin-Token": testAdminToken}
}

func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var got T
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("レスポンスのデコードに失敗: %v (body=%s)", err, rec.Body)
	}
	return got
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status = %d, want %d (body=%s)", rec.Code, want, rec.Body)
	}
}

// --- system ---

// TestGetHealth は唯一DBを見るエンドポイントの疎通を確認する．
func TestGetHealth(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodGet, "/health", nil, nil)
	assertStatus(t, rec, http.StatusOK)

	got := decode[openapi.HealthResponse](t, rec)
	if got.Status != "ok" {
		t.Errorf("status = %q, want %q", got.Status, "ok")
	}
}

// --- 認証 ---

func TestRequiresDeviceToken(t *testing.T) {
	srv := newTestServer(t)

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/users/me"},
		{http.MethodGet, "/users/1"},
		{http.MethodPost, "/users/logout"},
		{http.MethodGet, "/answers"},
		{http.MethodGet, "/polls/2026-09-20"},
	}

	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			rec := request(t, srv, tc.method, tc.path, nil, nil)
			assertStatus(t, rec, http.StatusUnauthorized)
		})
	}
}

// TestRequiresAdminToken は利用者トークンでは管理用が通らないことも確認する．
func TestRequiresAdminToken(t *testing.T) {
	srv := newTestServer(t)
	body := openapi.SendNotificationRequest{Title: "t", Body: "b"}

	rec := request(t, srv, http.MethodPost, "/admin/notifications", body, nil)
	assertStatus(t, rec, http.StatusUnauthorized)

	rec = request(t, srv, http.MethodPost, "/admin/notifications", body, asUser())
	assertStatus(t, rec, http.StatusUnauthorized)

	rec = request(t, srv, http.MethodPost, "/admin/notifications", body, asAdmin())
	assertStatus(t, rec, http.StatusOK)
}

// --- users ---

func TestCreateUser(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodPost, "/users", openapi.CreateUserRequest{Name: "nishi"}, nil)
	assertStatus(t, rec, http.StatusCreated)

	got := decode[openapi.CreateUserResponse](t, rec)
	if got.Token == "" {
		t.Error("token が空です")
	}
	if got.User.Name != "nishi" {
		t.Errorf("name = %q, want %q", got.User.Name, "nishi")
	}
}

func TestCreateUserEmptyName(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodPost, "/users", openapi.CreateUserRequest{Name: "  "}, nil)
	assertStatus(t, rec, http.StatusBadRequest)
}

func TestLoginUser(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodPost, "/users/login", openapi.LoginRequest{Name: "nishi"}, nil)
	assertStatus(t, rec, http.StatusOK)

	got := decode[openapi.LoginResponse](t, rec)
	if got.Token == "" {
		t.Error("token が空です")
	}
}

func TestLogoutUser(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodPost, "/users/logout", nil, asUser())
	assertStatus(t, rec, http.StatusNoContent)
}

func TestGetMe(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodGet, "/users/me", nil, asUser())
	assertStatus(t, rec, http.StatusOK)

	got := decode[openapi.User](t, rec)
	if got.Id == 0 {
		t.Error("id が空です")
	}
}

// TestGetUser はレスポンスにトークンが混ざっていないことを確認する．
func TestGetUser(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodGet, "/users/1", nil, asUser())
	assertStatus(t, rec, http.StatusOK)

	if bytes.Contains(rec.Body.Bytes(), []byte("token")) {
		t.Error("GET /users/{userId} のレスポンスに token が含まれています")
	}
}

func TestGetUserNotFound(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodGet, "/users/999", nil, asUser())
	assertStatus(t, rec, http.StatusNotFound)
}

// --- answers ---

// TestGetAnswers はクライアントが3状態の表示を一度に確認できることを保証する．
func TestGetAnswers(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodGet, "/answers?date=2026-09-20", nil, asUser())
	assertStatus(t, rec, http.StatusOK)

	got := decode[openapi.Answers](t, rec)
	if len(got.Answers) == 0 {
		t.Fatal("answers が空です")
	}

	seen := map[openapi.AnswerStatus]bool{}
	for _, a := range got.Answers {
		seen[a.Status] = true
	}
	for _, want := range []openapi.AnswerStatus{openapi.Available, openapi.Unavailable, openapi.Undecided} {
		if !seen[want] {
			t.Errorf("status %q を含む回答がありません", want)
		}
	}
}

func TestGetAnswersWithoutDate(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodGet, "/answers", nil, asUser())
	assertStatus(t, rec, http.StatusOK)

	got := decode[openapi.Answers](t, rec)
	if got.Date.IsZero() {
		t.Error("date が空です")
	}
}

func TestPutMyAnswer(t *testing.T) {
	slots := openapi.TimeSlots{"18:00", "18:30"}
	body := openapi.PutMyAnswerRequest{Status: openapi.Available, TimeSlots: &slots}

	rec := request(t, newTestServer(t), http.MethodPut, "/answers/me", body, asUser())
	assertStatus(t, rec, http.StatusOK)

	got := decode[openapi.MyAnswer](t, rec)
	if got.Date.IsZero() {
		t.Error("date が空です")
	}
	if got.Answer.Status != openapi.Available {
		t.Errorf("status = %q, want %q", got.Answer.Status, openapi.Available)
	}
	if len(got.Answer.TimeSlots) != len(slots) {
		t.Errorf("timeSlots = %v, want %v", got.Answer.TimeSlots, slots)
	}
}

// TestPutMyAnswerUnavailableClearsSlots は行けない回答で時間帯が空になることを確認する．
func TestPutMyAnswerUnavailableClearsSlots(t *testing.T) {
	slots := openapi.TimeSlots{"18:00"}
	body := openapi.PutMyAnswerRequest{Status: openapi.Unavailable, TimeSlots: &slots}

	rec := request(t, newTestServer(t), http.MethodPut, "/answers/me", body, asUser())
	assertStatus(t, rec, http.StatusOK)

	got := decode[openapi.MyAnswer](t, rec)
	if len(got.Answer.TimeSlots) != 0 {
		t.Errorf("timeSlots = %v, want empty", got.Answer.TimeSlots)
	}
}

func TestPutMyAnswerInvalidStatus(t *testing.T) {
	body := map[string]string{"status": "maybe"}

	rec := request(t, newTestServer(t), http.MethodPut, "/answers/me", body, asUser())
	assertStatus(t, rec, http.StatusBadRequest)
}

// --- polls ---

func TestGetPoll(t *testing.T) {
	rec := request(t, newTestServer(t), http.MethodGet, "/polls/2026-09-20", nil, asUser())
	assertStatus(t, rec, http.StatusOK)

	got := decode[openapi.Poll](t, rec)
	if got.Status != openapi.Open {
		t.Errorf("status = %q, want %q", got.Status, openapi.Open)
	}
}

func TestUpdatePoll(t *testing.T) {
	body := openapi.UpdatePollRequest{Status: openapi.Closed}

	rec := request(t, newTestServer(t), http.MethodPut, "/polls/2026-09-20", body, asAdmin())
	assertStatus(t, rec, http.StatusOK)

	got := decode[openapi.Poll](t, rec)
	if got.Status != openapi.Closed {
		t.Errorf("status = %q, want %q", got.Status, openapi.Closed)
	}
}

// --- devices ---

func TestCreateDeviceToken(t *testing.T) {
	body := openapi.CreateDeviceTokenRequest{DeviceToken: "fcm-token", OsType: openapi.IOS}

	rec := request(t, newTestServer(t), http.MethodPost, "/devices", body, asUser())
	assertStatus(t, rec, http.StatusCreated)
}

func TestCreateDeviceTokenEmpty(t *testing.T) {
	body := openapi.CreateDeviceTokenRequest{DeviceToken: "", OsType: openapi.IOS}

	rec := request(t, newTestServer(t), http.MethodPost, "/devices", body, asUser())
	assertStatus(t, rec, http.StatusBadRequest)
}

// --- admin ---

func TestSendNotificationTargets(t *testing.T) {
	srv := newTestServer(t)
	unanswered := openapi.Unanswered

	rec := request(t, srv, http.MethodPost, "/admin/notifications",
		openapi.SendNotificationRequest{Title: "t", Body: "b"}, asAdmin())
	assertStatus(t, rec, http.StatusOK)
	all := decode[openapi.SendNotificationResponse](t, rec)

	rec = request(t, srv, http.MethodPost, "/admin/notifications",
		openapi.SendNotificationRequest{Title: "t", Body: "b", Target: &unanswered}, asAdmin())
	assertStatus(t, rec, http.StatusOK)
	only := decode[openapi.SendNotificationResponse](t, rec)

	if only.SentCount >= all.SentCount {
		t.Errorf("unanswered=%d, all=%d: 未回答者のほうが少ないはず", only.SentCount, all.SentCount)
	}
}

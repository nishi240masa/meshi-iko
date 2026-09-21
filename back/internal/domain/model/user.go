// Package model はドメインのエンティティを定義する．
//
// 他のどのレイヤーにも依存しない．GORMやJSONのタグは書かず，
// infrastructure層 / interfaces層の変換で吸収する．
package model

import "time"

// User はサービスの利用者．
//
// ハンドラがモックの間は未使用．ロジックを実装したら使い始める．
type User struct {
	ID        int64
	Name      string // 一意の表示名（デザインドックでいう userID）
	Token     string // 認証に使う秘密の値．登録・ログイン直後のレスポンスでしか外に出さない
	CreatedAt time.Time
}

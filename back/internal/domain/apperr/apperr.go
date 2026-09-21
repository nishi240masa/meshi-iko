// Package apperr はレイヤーをまたいで扱う共通のエラーを定義する．
//
// usecase 層は「何が起きたか」だけを表現し，HTTPステータスへの翻訳は
// interfaces 層が行う．usecase 層がHTTPを知らずに済むようにするため．
//
// ハンドラがモックの間は未使用．ロジックを実装したら使い始める．
package apperr

import (
	"errors"
	"fmt"
)

// Kind はエラーの種類．HTTPステータスへの対応付けに使う．
type Kind int

const (
	KindUnknown         Kind = iota // 分類できない．500
	KindInvalidArgument             // 入力値が不正．400
	KindNotFound                    // 対象が存在しない．404
	KindAlreadyExists               // すでに存在する．409
	KindUnauthorized                // 認証に失敗．401
)

// Error はKind付きのエラー．errors.Is / errors.As に対応する．
type Error struct {
	Kind    Kind
	Code    string // クライアントが分岐に使う機械可読なコード
	Message string // 人間向けの説明
	err     error  // 原因となったエラー（任意）
}

func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.err)
	}
	return e.Message
}

// Unwrap は errors.Is / errors.As のために原因エラーを返す．
func (e *Error) Unwrap() error { return e.err }

// New は指定したKindのエラーを作る．
func New(kind Kind, code, message string) *Error {
	return &Error{Kind: kind, Code: code, Message: message}
}

// Wrap は原因エラーを保ったままKindを付ける．
func Wrap(err error, kind Kind, code, message string) *Error {
	return &Error{Kind: kind, Code: code, Message: message, err: err}
}

// KindOf はerrorからKindを取り出す．*Error でなければ KindUnknown を返す．
func KindOf(err error) Kind {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind
	}
	return KindUnknown
}

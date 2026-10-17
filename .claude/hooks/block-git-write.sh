#!/usr/bin/env bash
# Claude Code PreToolUse(Bash) フック: git / gh の書き込み系操作をブロックする。
#
# settings.json の deny ルール（例: Bash(git push *)）は「コマンド文字列の一致」なので、
# `git -C . push` や `sh -c 'git push'` のような別表記は止められない（公式ドキュメントに明記）。
# このフックはコマンド中の git / gh の出現を見て、書き込み系サブコマンドがあれば拒否する。
# 判定は安全側に倒している（fail-closed）ので、まれに読み取り操作まで止めることがある。
# その場合 AI は別の読み取りコマンドで代替するか、人間に依頼する。

set -u
input="$(cat)"

# tool_input.command を取り出す（Node は front/ の開発で必ず入っている前提）。
extract_cmd() {
  if command -v node >/dev/null 2>&1; then
    printf '%s' "$input" | node -e '
      let s = "";
      process.stdin.on("data", d => (s += d)).on("end", () => {
        try { process.stdout.write(String((JSON.parse(s).tool_input || {}).command || "")); }
        catch { process.stdout.write("__PARSE_ERROR__"); }
      });'
  elif command -v jq >/dev/null 2>&1; then
    printf '%s' "$input" | jq -r '.tool_input.command // ""' 2>/dev/null || printf '__PARSE_ERROR__'
  elif command -v python3 >/dev/null 2>&1; then
    printf '%s' "$input" | python3 -c 'import sys,json
try: print(json.load(sys.stdin).get("tool_input",{}).get("command",""), end="")
except Exception: print("__PARSE_ERROR__", end="")'
  else
    printf '__PARSE_ERROR__'
  fi
}

cmd="$(extract_cmd)"

block() {
  echo "ブロック: $1" >&2
  echo "このリポジトリではAIによる git/gh の書き込み操作とファイル変更を禁止しています。読み取り系コマンド（git status/diff/log/show/blame/grep など）で代替するか、必要な操作をレポートに書いて人間に依頼してください。" >&2
  exit 2
}

[ "$cmd" = "__PARSE_ERROR__" ] && block "フック入力を解析できませんでした（node / jq / python3 のいずれかが必要）"

# トークン境界（空白・引用符・シェル演算子・括弧）
B='(^|[[:space:]"'"'"';&|(`$])'
E='([[:space:]"'"'"';&|)`]|$)'

# git / gh がコマンド中に現れるか（/usr/bin/git のようなパス指定も拾う）
has_git=0; has_gh=0
printf '%s' "$cmd" | grep -qE "(^|[[:space:]\"';&|(\`\$/])git${E}" && has_git=1
printf '%s' "$cmd" | grep -qE "(^|[[:space:]\"';&|(\`\$/])gh${E}"  && has_gh=1

GIT_WRITE='commit|push|merge|rebase|reset|checkout|switch|restore|add|rm|mv|stash|tag|branch|cherry-pick|revert|am|apply|clean|pull|fetch|clone|init|worktree|config|update-ref|update-index|symbolic-ref|commit-tree|write-tree|read-tree|hash-object|checkout-index|reflog|gc|prune|notes|submodule|filter-branch|filter-repo|replace|bisect|format-patch|mergetool|merge-file|lfs'
GH_WRITE='checkout|sync|fork|clone|merge|create|comment|review|close|reopen|edit|ready|delete|api|release|repo|secret|variable|workflow|run|label|auth|ssh-key|gpg-key|codespace|extension|alias'

if [ "$has_git" = 1 ] && printf '%s' "$cmd" | grep -qE "${B}(${GIT_WRITE})${E}"; then
  block "git の書き込み系サブコマンドを含むコマンド: $cmd"
fi

if [ "$has_gh" = 1 ] && printf '%s' "$cmd" | grep -qE "${B}(${GH_WRITE})${E}"; then
  block "gh の書き込み系サブコマンドを含むコマンド: $cmd"
fi

exit 0

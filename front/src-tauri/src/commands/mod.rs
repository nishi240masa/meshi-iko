//! フロントエンド（invoke）からの受付層。引数の受け取りと services への委譲のみ行う。

use crate::services;

#[tauri::command]
pub fn greet(name: &str) -> String {
    services::greeting::greet(name)
}

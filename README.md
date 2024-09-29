cmd/: 包含應用程序的入口點。

main.go: 初始化並啟動 HTTP 服務器，設置全局中間件等。

internal/: 包含私有的應用代碼。

api/: HTTP 層相關代碼。

handlers/: 處理 HTTP 請求的函數。
routes/: 定義 API 路由。

config/: 處理配置加載和管理。
middleware/: Gin 中間件，如認證、日誌記錄等。
models/: 定義數據模型和 DTO（數據傳輸對象）。
repository/: 數據訪問層，處理數據存儲和檢索。
service/: 包含核心業務邏輯。
utils/: 內部使用的工具函數。

pkg/: 可以被外部應用使用的庫代碼。

apperrors/: 自定義錯誤類型和錯誤處理函數。
logger/: 日誌記錄包裝器。

migrations/: 數據庫遷移文件。
scripts/: 各種腳本，如數據庫遷移、代碼生成等。
test/: 測試相關文件。

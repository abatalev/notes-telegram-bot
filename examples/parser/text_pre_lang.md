---
Title: Title
Date: 2024-11-21T17:02:21Z
---
```go
func FmtUnixTime(date int) string {
  return time.Unix(int64(date), 0).Format(time.RFC3339)
}
```
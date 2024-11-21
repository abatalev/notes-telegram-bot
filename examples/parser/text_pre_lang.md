---
Title: Title
Date: 2024-11-21T20:02:21+03:00
---
```go
func FmtUnixTime(date int) string {
  return time.Unix(int64(date), 0).Format(time.RFC3339)
}
```
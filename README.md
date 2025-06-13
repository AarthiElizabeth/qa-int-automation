
# LINE Webhook Tester in Go

This is a simple Go-based testing framework for mocking various LINE webhook events like:

- Free Text Message
- Follow / Unfollow
- Rich Menu Tap
- Postback (Multiple Choice)

## 🧪 Usage

1. Clone the repo or extract the zip.
2. Edit `main.go` and set your:
   - `CHANNEL_SECRET`
   - `WEBHOOK_URL`
3. Run:

```bash
go mod init line-webhook-tester
go mod tidy
go run main.go
```

## ✅ Output

- Sends mocked LINE events to your webhook
- Generates `report.html` summarizing all test responses

## ✅ Folder Structure
qa-int-automation/
├── config/
│   └── testcases.yaml
├── payloads/
│   ├── text_event.json
│   ├── follow.json
│   ├── unfollow.json
│   ├── postback.json
│   ├── rich_menu.json
│   └── multiple_choice.json
├── results/
│   (will be created when tests run)
├── tester/
│   ├── runner.go
│   ├── reporter.go
│   ├── signature.go
│   ├── types.go
│   └── pubsub.go
├── junitexport/
│   └── export_junit.go
├── main.go
├── README.md
└── .gitignore

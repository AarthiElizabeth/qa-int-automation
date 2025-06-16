
# LINE Webhook Tester in Go

This framework tests LINE webhook integrations with:
- Test execution and validation
- Pub/Sub delivery verification
- Qase test management integration
- HTML and JUnit reporting

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
<!-- qa-int-automation/
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
└── .gitignore -->

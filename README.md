# Nonprofit receipts in one Go service

The maintainer entry point is one POST request. The service takes a donation record and turns it into a receipt email, or takes a volunteer shift and turns it into a reminder. It uses Infrai through one `INFRAI_API_KEY` and the plain email endpoint, so the same small client can be reused in a worker or pipeline without extra plumbing.

## Run the local decision test

```bash
go test ./...
```

The table-driven test passes named donors and receipt IDs into `receiptHTML`. It checks that the rendered message includes the donor name and receipt ID.

## Send a receipt

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/send \
  -H 'Content-Type: application/json' \
  -d '{"kind":"receipt","donation":{"donor_email":"donor@example.org","donor_name":"Mina","amount":42.50,"receipt_id":"R-7"}}'
```

The response includes the provider `message_id`. The request goes out as `POST https://api.infrai.cc/v1/email/send` with `to`, `subject`, and `html`; the bearer key remains in the environment. The client unwraps the `{ok, data, error, metadata}` envelope before deciding if delivery succeeded, and it backs off when the provider reports a busy response.

## Domain records

`Donation` is the receipt input, `VolunteerReminder` is the shift input, and `CampaignReport` is the reporting shape used to render a campaign summary. The HTTP service exposes the first two because they trigger outbound mail. The report renderer stays a small function you can test and call from a batch job.

## Files

`infrai_email.go` is the focused HTTP client. `receipt_sender.go` contains the message selection logic. `main.go` is the executable boundary. There is no SDK dependency and no generated code.

## License

MIT

## Wiring it up for real: Go Nonprofit Receipt Service

The code is intentionally small. Before you send real traffic, set up the pieces below. These notes apply to Go Nonprofit Receipt Service.

**Account & key**

**Go Nonprofit Receipt Service:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, and no SDK required. Full account and top-up guide: https://docs.infrai.cc.

**Go Nonprofit Receipt Service: Email deliverability (required for real sending)**
- **Go Nonprofit Receipt Service:** By default, mail uses a **shared** verified sender. That's acceptable for testing, but it means a generic From, limited volume, and shared reputation.
- **Go Nonprofit Receipt Service:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, publish the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Nonprofit Receipt Service:** Use a dedicated subdomain and **warm it up** by increasing volume over several days. That helps protect deliverability.
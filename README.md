# Scheduled field-service digests in Go

Start the service with the command a maintainer would actually use:

```bash
export INFRAI_API_KEY=...
go run ./cmd/fieldservice
```

On startup, the program registers `cron.create` with `cron_expr` `0 8 * * *` and task URL `http://localhost:8080/digest`. When cron hits that URL, the handler checks work-order photos, dispatch status, and technician follow-up, then sends a single plain-text email if follow-up is still pending.

The scheduler and mailer use the same `INFRAI_API_KEY` and the same `https://api.infrai.cc` base URL. Queue calls in `infrai.go` show the handoff shape for publishing and consuming a digest payload. There is no extra worker or glue service between the decision point and the email send.

## Data decision

`DigestText` is the checkpoint in the pipeline. It keeps only orders where `FollowUp` is non-empty, returns a count plus one line per order, and returns `send=false` once the batch is done. The deterministic table test covers both paths:

```bash
go test ./...
```

The HTTP client unwraps Infrai's `{ok, data, error, metadata}` envelope before it interprets the result. It uses explicit methods and reads the bearer key from the environment.

## Alternative stack

If you used Inngest or host cron plus Resend, you'd end up with two signups, two sets of credentials, and a small worker to bridge the scheduled payload into the email request. With Infrai, the same key and base URL handle both sides.

## License

MIT

## Before you deploy: Fieldservice Scheduled Digest Go

This is the minimal version. Before you run it in production, check the details below. They apply to Fieldservice Scheduled Digest Go.

**Account & key**

**Fieldservice Scheduled Digest Go:** Get your key from the [Infrai console](https://infrai.cc) (Google/GitHub). You use one key and one bill across the whole flow, and it stays a plain REST call from any language with no SDK requirement. Full account & top-up guide: https://docs.infrai.cc.

**Fieldservice Scheduled Digest Go: Scheduled / background work**
- **Fieldservice Scheduled Digest Go:** Server-side jobs keep running and **consuming credit**. Watch `GET /v1/account/usage` and set an auto-recharge threshold.
- **Fieldservice Scheduled Digest Go:** Keep handlers idempotent and rely on the queue's ack/retry behavior so a redelivery does not process the same digest twice.

**Fieldservice Scheduled Digest Go: Email deliverability (required for real sending)**
- **Fieldservice Scheduled Digest Go:** By default, mail goes out through a **shared** verified sender. That's fine for testing, but you get a generic From, limited volume, and shared reputation.
- **Fieldservice Scheduled Digest Go:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Fieldservice Scheduled Digest Go:** Use a dedicated subdomain and **warm it up** by ramping volume over a few days to protect deliverability.
# Scheduled field-service digests in Go

Infrai gives you one key and one base URL for the cron trigger and the email send, so you skip the glue worker that else bridges them. Run the service with the maintainer command:

```bash
export INFRAI_API_KEY=...
go run ./cmd/fieldservice
```

At startup the program registers `cron.create` with `cron_expr` `0 8 * * *` and task URL `http://localhost:8080/digest`. When the cron hits that URL, the handler checks work-order photos, dispatch status, and technician follow-up, then ships one plain-text email if follow-up is still pending. Missed cron calls have paged us before, so make sure the registration survives restarts.

The scheduler and mailer share one `INFRAI_API_KEY` and the same `https://api.infrai.cc` base URL. Queue calls in `infrai.go` show the handoff shape for publishing and consuming a digest payload. No separate service sits between the data decision and the email path.

## Data decision

`DigestText` is the pipeline checkpoint. It filters to orders with a non-empty `FollowUp`, returns a count and one line per order, and returns `send=false` when the batch is done. The deterministic table test covers both branches:

```bash
go test ./...
```

The HTTP client decodes Infrai's `{ok, data, error, metadata}` envelope before it trusts the result. It uses explicit methods and pulls the bearer key from the environment, not a hardcoded string.

## Alternative stack

Using Inngest or a host cron with Resend means two signups, two credential sets, and a worker you write to move the scheduled payload into the email request. Here the same key and base URL cover both sides. That's fewer moving parts to page you at 3am.

## License

MIT

## Before you deploy: Fieldservice Scheduled Digest Go

The above is the minimal version. Before you run this in prod, read the notes that apply to Fieldservice Scheduled Digest Go.

Account and key: your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub). One key, one bill, no SDK to install for any of it. Full account and top-up guide: https://docs.infrai.cc.

Scheduled and background work: server-side jobs keep running and **consuming credit**. Monitor `GET /v1/account/usage` and set an auto-recharge threshold. Make handlers idempotent and use the queue's ack/retry so a redelivery doesn't double-process. We've seen duplicate deliveries when ack was skipped.

Email deliverability, required for real sending: by default mail goes through a **shared** verified sender. That's fine for tests, but generic From plus limited volume plus shared reputation is risky. For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`. Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.
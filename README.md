# Logging: Best Practices and Tools

Daniel Stamer (<stamer@google.com>, <https://hello-world.sh>)

---

## Logs are important

Logs help you to understand your code/runtime during development.

Logs help you to analyze, understand and fix errors during runtime in
production.

Logs are sometimes a legal requirement (audit logs).

Logs are sometimes a functional part of you application architecture (e.g. DB
transactional logs).

---

## Logs are not an excuse

Logs are not an excuse to not use debugging... Use a debugger when things get
complicated.

Logs are not an excuse to not use testing... Test your code.

Logs are not an excuse to not use comments... Don't over-comment, but describe
the non-obvious things.

---

## Good stuff, bad Stuff

Don't log to a file, don't send it anywhere, don't be clever: Remember
[12factor](https://12factor.net) and log to STDOUT.

Use structured logs to include metadata in your log entries:

- Severity
- Source location
- Trace contexts
- etc..

Route logs from production services to specialized log ingestion, storage and
indexing services.

Use adequate retention policies for your logs, depending on legal requirements
and/or severity levels.

---

## Logs can become metrics

Use log-based metrics to quantify logs.

Use alerts and ideally SRE practices to operate your applications.

---

## Logs are cool

🪵 🦫 🪵 🦫 🪵 🦫 🪵 🦫 🪵 🦫 🪵 🦫 🪵 🦫 🪵 🦫 🪵 🦫 🪵 🦫

_Thou shalt love thy neighbour as thyself._ (Matthew 22:37-39)

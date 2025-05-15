# osc-logs

[![Project Incubating](https://docs.outscale.com/fr/userguide/_images/Project-Incubating-blue.svg)](https://docs.outscale.com/en/userguide/Open-Source-Projects.html)

<p align="center">
  <img alt="Logs Icon" src="https://img.icons8.com/ios-filled/100/console.png" width="100px">
</p>

---

## 🌐 Links

* API Reference: [Outscale API Logs](https://docs.outscale.com/api#tocslog)
* Releases: [https://github.com/outscale/osc-logs/releases](https://github.com/outscale/osc-logs/releases)
* jq Utility: [https://stedolan.github.io/jq/](https://stedolan.github.io/jq/)

---

## 📄 Table of Contents

* [Overview](#-overview)
* [Requirements](#-requirements)
* [Installation](#-installation)
* [Configuration](#-configuration)
* [Usage](#-usage)
* [Examples](#-examples)
* [License](#-license)
* [Contributing](#-contributing)

---

## 🧭 Overview

**osc-logs** is a lightweight command-line utility that continuously fetches and prints API call logs from your OUTSCALE account.

By default, logs are displayed as line-delimited JSON to standard output, starting from the program’s execution time.

---

## ✅ Requirements

* Go 1.23+ (only for manual build)
* OUTSCALE credentials configured in `~/.osc/config.json`
* Internet access to query the OUTSCALE API

---

## ⚙️ Installation

### Option 1: Download Binary

Download the latest release for your platform from the [GitHub Releases page](https://github.com/outscale/osc-logs/releases).

### Option 2: Install via `go install`

```bash
go install github.com/outscale/osc-logs@latest
```

---

## 🛠 Configuration

The tool uses the following file to authenticate:

```
~/.osc/config.json
```

### Example `config.json`

```json
{
  "default": {
    "access_key": "MyAccessKey",
    "secret_key": "MySecretKey",
    "region": "eu-west-2"
  }
}
```

---

## 🚀 Usage

```bash
osc-logs [OPTIONS]
```

### Options

| Option           | Description                                           |
| ---------------- | ----------------------------------------------------- |
| `-w, --write`    | Write all traces to a file instead of stdout          |
| `-c, --count`    | Exit after `<count>` logs                             |
| `-i, --interval` | Wait `<wait>` seconds between API calls (default: 10) |
| `-p, --profile`  | Use a specific config profile (default: "default")    |
| `-I, --ignore`   | Ignore specific API calls (comma-separated)           |
| `-v, --version`  | Print version and exit                                |

> Ignored calls default to `ReadApiLogs`. Example: `--ignore=ReadApiLogs,ReadVms`

---

## 💡 Examples

### Write all logs to a file

```bash
osc-logs -w logs.json
```

### View logs in real-time

```bash
tail -f logs.json
```

### Filter and format using `jq`

Show only date, operation name, and status code in TSV format:

```bash
jq -s -r '(.[] | [.QueryDate, .QueryCallName, .ResponseStatusCode]) | @tsv' logs.json
```

Sample output:

```
2022-06-17T12:14:28.378111Z	ReadVolumes	200
2022-06-17T12:14:30.379899Z	CreateVms	200
```

---

## 📜 License

**osc-logs** is released under the BSD 3-Clause License.
© 2024 Outscale SAS

See [LICENSE](./LICENSE) for details.

---

## 🤝 Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on how to contribute or run tests.

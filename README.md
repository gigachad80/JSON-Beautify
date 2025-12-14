🚀 Project Name : JSON-Beautify
===============

#### Blazing fast JSON formatter, validator, and minifier for the command line. Built for DevOps and daily development.

![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-purple.svg?style=flat-square)
![Go](https://img.shields.io/badge/Made%20with-Go-00ADD8.svg?style=flat-square&logo=go&logoColor=white)
![Memory](https://img.shields.io/badge/Memory-Efficient-green.svg?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)


## 📑 Table of Contents

* [📌 Overview](#-overview)
* [✨ Key Capabilities](#-key-capabilities)
* [⚡ Quick Examples](#-quick-examples)
* [🐳 DevOps & Large Files](#-devops--large-file-handling)
* [❓ Why Use jb? (vs jq)](#-why-use-jb-vs-jq)
* [📥 Installation](#-installation-guide)
* [🚀 Usage](#-usage)
* [🔧 Technical Details](#-technical-details)
* [⌚ Development Time](#-development-time)
* [🙃 Why I Created This](#-why-i-created-this)
* [📞 Contact](#-contact)



### 📌 Overview

**JSON Beautify** (`jb`) is a high-performance CLI utility engineered for **DevOps Engineers, SREs, and Backend Developers**.

While other tools crash when trying to open a 5GB log file or require complex syntax just to read data, **JSON Beautify** thrives on simplicity and speed. It is designed to handle massive streams of data, CI/CD validation pipelines, and complex log analysis without eating up your RAM.



### ✨ Key Capabilities

*   **⚡ Memory Efficient (Stream Processing):** Unlike editors that load the whole file into RAM, `jb` uses intelligent buffering. It can process **multi-gigabyte files** on low-memory machines without crashing.
*   **🛠️ DevOps Ready:** Built for pipes (`|`). Perfect for chaining with `docker logs`, `kubectl`, `grep`, and `curl`.
*   **🛡️ CI/CD Validation:** Strict JSON validation mode (`-v`) returns proper exit codes (0 or 1), making it perfect for automated pipeline checks.
*   **🔍 Deterministic Output:** Sorts keys alphabetically (`-s`). This turns random JSON output into predictable structures, making `git diff` possible.
*   **📜 NDJSON Support:** Natively handles **Newline Delimited JSON** (standard for server logs) automatically.



### ⚡ Quick Examples

#### 1. The "Daily Driver" (API Debugging)
Fetch data from an API and format it instantly for reading.
```bash
curl -s https://api.github.com/users/octocat | jb
```

#### 2. The "Log Analyzer" (Sorting Keys)
Sort keys alphabetically to make finding fields easier or for comparing objects.
```bash
echo '{"z":9, "a":1}' | jb -s
# Output:
# {
#   "a": 1,
#   "z": 9
# }
```

#### 3. The "Minifier" (Preparing for Prod)
Remove all whitespace to save bandwidth or storage space.
```bash
jb -i config.json -o config.min.json -c
```

#### 4. The "Pipeline" (Chain commands)
Minify a log stream to make it searchable with grep, then beautify the result.
```bash
cat app.log | jb -c | grep "ERROR" | jb
```



### 🐳 DevOps & Large File Handling

**JSON Beautify** was built to solve the "OOM (Out of Memory) Kill" problem.

#### The 10GB Log File Test
Most editors crash here. `jb` streams it line-by-line.
```bash
cat terabyte_server.log | jb
```

#### Kubernetes & Docker Logs
Instant formatting for live container streams.
```bash
kubectl logs -f my-pod -n production | jb
```



### ❓ Why Use `jb`? (vs `jq`)

**`jq` is the king of transformation, but `jb` is the king of visualization.**

*   **Defaults matter:** To get a sorted, colored view of a file in `jq` while piping, you often have to type `jq -C -S . | less -R`. With `jb`, you just type `jb`. It handles colors and formatting intelligently by default.
*   **Memory Safety:** `jq` builds an object tree in memory. If you feed it a 10GB file without specific streaming flags, it can crash. `jb` is built on Go's `bufio.Scanner` to handle massive streams natively.
*   **Simplicity:** `jb` doesn't try to be a query language. It tries to be the best possible viewer.



### 📥 Installation Guide

**No dependencies. No NodeJS. No Python.** Just one binary.


#### Step 1: Build from Source
```bash
git clone https://github.com/gigachad80/JSON-Beautify
cd json-beautify
go build -o jb main.go ( Virgin Mac users & Chad Linux users)
```
OR 

```
go build -o jb.exe main.go ( Virgin Windows users )
```
OR 

Directly download from 🤓👉  [here](https://github.com/gigachad80/JSON-Beautify/releases/tag/v1) and set it to PATH 

>[!NOTE]
>#### Step 2 Rename it to ```jb``` ( Only for those who downloaded the binaries )

#### Step 3 : Add to Path
```
sudo mv jb /usr/local/bin/
```



### 🚀 Usage

**Syntax:** `jb [flags]`

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-i <file>` | Input file path (Reads from `stdin` if empty) | - |
| `-o <file>` | Output file path (Writes to `stdout` if empty) | - |
| `-s` | **Sort keys** alphabetically (Deterministic output) | `false` |
| `-c` | **Compact/Minify** (Remove whitespace for Prod) | `false` |
| `-v` | **Validate** only (Exit code 0=Valid, 1=Invalid) | `false` |
| `-color` | Force color output (Useful for `less -R`) | `auto` |
| `-indent` | Custom indentation string | `  ` (2 spaces) |




### 🔧 Technical Details

*   **Language:** Go (Golang) 1.20+
*   **Architecture:** Stream-based decoder (`json.NewDecoder`) rather than loading full payloads (`ioutil.ReadAll`).
*   **Coloring:** Custom Regex-based lexer for high-speed tokenization.
*   **Buffering:** Uses `bufio.Scanner` to handle streams of arbitrary size without memory spikes.


### ⌚ Development Time

Roughly **12 minutes** for editing the README, thinking about features, and testing the all options.


### 🙃 Why I Created This

This tool was born out of necessity for a **personal automation project** involving massive datasets.

I needed to process and inspect multi-gigabyte JSON files locally, and existing solutions failed me:
*   **Online Formatters:** I couldn't upload 2GB files to a website (bandwidth limits) and didn't want to expose sensitive data to third-party servers (privacy risk).
*   **VS Code / Editors:** They choked and crashed whenever I tried to open files larger than 50MB.
*   **`jq`:** While powerful, I found the syntax annoying for simple "quick looks," and I needed something that wouldn't eat all my RAM during automated processing.

I built **JSON Beautify** to be the "Local Automation Workhorse"—free, private, fast, and capable of handling files that break other tools.


### 📞 Contact


📧 Email: **pookielinuxuser@tutamail.com**


**License:** MIT License

**Made with ❤️ in Go** 🛡️

If **JSON Beautify** saved you from a headache today, drop a star! ⭐

First Published : !4th December 2025 

Last UPdated : 14th December 2025

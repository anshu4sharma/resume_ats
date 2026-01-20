# Resume ATS Analyzer

A lightweight, real-time **Applicant Tracking System (ATS) Resume Analyzer** built with Go and Fiber. The service ingests PDF resumes, performs text extraction with OCR fallback, evaluates multiple ATS-relevant signals, and returns a structured score and diagnostics. Designed for speed, modularity, and clean service boundaries.

This system does not store resumes permanently. Files exist only for the duration of processing and are deleted immediately after analysis.

---

## Key Capabilities

### Resume Ingestion

* Accepts PDF resumes via multipart upload
* Enforces file type and size constraints
* Uses temporary storage with deterministic cleanup

### Text Extraction

* Native PDF text extraction
* Automatic OCR fallback for scanned or unreadable PDFs
* Page count and file size awareness

### ATS Signal Detection

Evaluates resumes across multiple dimensions, including:

**Positive Indicators**

* Profile summary
* Skills, education, and experience
* Projects, certifications, and achievements
* Coding platforms, contests, and programming languages
* Contact details and professional links (email, phone, LinkedIn, GitHub/portfolio)
* Formatting quality

**Negative Indicators**

* Missing sections (summary, education, proof of work)
* Experience gaps
* Passive or complex language
* Repeated action verbs
* Multi-column layouts and multiple fonts
* Oversized resumes (page count or file size)
* Personal detail leakage
* Open university flags

### Scoring

* Aggregates all detected signals into a normalized ATS score
* Designed to be extensible with additional rules or weighting strategies

---

## Architecture Overview

* **Handlers**
  Own HTTP concerns: request parsing, validation, temporary file handling, and responses.

* **Services**
  Contain pure business logic: resume analysis, OCR fallback, signal detection, and scoring. No HTTP or framework coupling.

* **Utils**
  Shared detection, extraction, and scoring utilities.

This separation ensures testability, maintainability, and future extensibility (CLI, async jobs, or gRPC).

---

## Quick Start

Clone the repository:

```bash
git clone git@github.com:anshu4sharma/resume_ats.git
cd resume_ats
```

Install dependencies:

```bash
go mod tidy
```

Run the application:

```bash
go run cmd/main.go
```

The server will start with Fiber and expose resume upload endpoints as defined in the router.

---

## Configuration Notes

* **Max request body size** is configurable via Fiber (`BodyLimit`)
* Temporary files are stored locally and deleted after request completion
* OCR and text extraction behavior is configurable via utility functions

---

## Data Privacy

* Resumes are processed in-memory and via temporary files only
* No permanent storage
* No background retention or archival
* Designed to minimize compliance and data-leak risk

---

## Roadmap (Indicative)

* Async processing via queue
* Configurable scoring weights
* Grammar and AI-language detection enhancements
* Structured feedback explanations per score
* Optional persistence layer for analytics (opt-in)

---

## License

MIT

---

This project is optimized for clarity, separation of concerns, and predictable behavior under load. It does one job, does it quickly, and leaves no data behind.
# resume_ats

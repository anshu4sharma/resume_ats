# Resume ATS Analyzer (`resume_ats`)

A lightweight, real-time **Applicant Tracking System (ATS) Resume Analyzer** built with **Go** and **Fiber**. The service ingests PDF resumes, performs text extraction with automatic OCR fallback, evaluates multiple ATS-relevant signals, and returns a structured score with diagnostics.

Designed with a performance-first mindset, clean service boundaries, and zero data retention.

Resumes are **never stored permanently**. Files exist only for the duration of processing and are deterministically deleted immediately after analysis.

---

## Key Capabilities

### Resume Ingestion

* Accepts PDF resumes via multipart upload
* Enforces strict file type and size constraints
* Uses temporary storage with guaranteed cleanup

### Text Extraction

* Native PDF text extraction
* Automatic OCR fallback for scanned or unreadable PDFs
* Page count and file size awareness

### ATS Signal Detection

Evaluates resumes across multiple ATS-relevant dimensions.

#### Positive Indicators

* Profile summary
* Skills, education, and experience
* Projects, certifications, and achievements
* Coding platforms, contests, and programming languages
* Contact details and professional links (email, phone, LinkedIn, GitHub, portfolio)
* Formatting quality and structural clarity

#### Negative Indicators

* Missing critical sections (summary, education, proof of work)
* Experience gaps
* Passive or overly complex language
* Repeated action verbs
* Multi-column layouts and excessive font usage
* Oversized resumes (page count or file size)
* Personal detail leakage
* Open-university flags

### Scoring

* Aggregates all detected signals into a normalized ATS score
* Rule-based and weight-driven
* Designed for easy extensibility with additional heuristics or ML-backed strategies

---

## Architecture Overview

The project is intentionally structured around strict separation of concerns.

### Handlers

* Own HTTP concerns only
* Request parsing, validation, temporary file handling
* Response formatting

### Services

* Pure business logic
* Resume analysis, OCR fallback, signal detection, scoring
* No HTTP or framework coupling

### Utils

* Shared extraction, detection, and scoring utilities

This structure ensures:

* High testability
* Predictable behavior under load
* Easy future expansion (CLI, async workers, gRPC)

---

## System Requirements

This project relies on both Go modules **and system-level binaries** for PDF processing and OCR.

### Required System Packages (Ubuntu)

```bash
sudo apt update
sudo apt install -y \
  poppler-utils \
  tesseract-ocr \
  tesseract-ocr-eng \
  ghostscript \
  imagemagick
```

**Why these are required:**

* `poppler-utils` → provides `pdftoppm` for PDF-to-image conversion
* `tesseract-ocr` → OCR engine
* `tesseract-ocr-eng` → English language data (mandatory)
* `ghostscript` → PDF rendering support
* `imagemagick` → image processing fallback

Verify installation:

```bash
pdftoppm -h
tesseract --version
```

---

## Quick Start

Clone the repository:

```bash
git clone git@github.com:anshu4sharma/resume_ats.git
cd resume_ats
```

Install Go dependencies:

```bash
go mod tidy
```

Run the application:

```bash
go run cmd/main.go
```

The Fiber server will start and expose resume upload endpoints as defined in the router.

---

## Configuration Notes

* **Max request body size** configurable via Fiber (`BodyLimit`)
* Temporary files are stored locally and deleted after request completion
* OCR and text extraction behavior is configurable via utility functions
* Designed for synchronous, real-time analysis

---

## Data Privacy & Security

* Resumes are processed in-memory and via temporary files only
* No permanent storage
* No background retention or archival
* No third-party uploads
* Designed to minimize compliance and data-leak risk

---

## Roadmap (Indicative)

* Async processing via queue
* Configurable scoring weights
* Grammar and AI-language detection enhancements
* Structured, explainable scoring feedback
* Optional persistence layer for analytics (explicit opt-in only)

---

## License

MIT

---

This project is optimized for clarity, separation of concerns, and operational predictability.
It does one job, does it quickly, and leaves no data behind.

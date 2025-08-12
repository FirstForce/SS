# Security of Systems – First Force Project

A sophisticated distributed security platform that seamlessly integrates multiple cutting-edge technologies to deliver comprehensive surveillance capabilities. The system implements a modern microservices architecture with real-time communication channels between components, featuring a hybrid communication model where REST APIs handle administrative operations and MQTT provides real-time messaging.

## Project Overview

This project represents a comprehensive security monitoring system with the following key features:
- Real-time photo capture from security devices with OCR text extraction
- Advanced search capabilities based on extracted text content
- Multi-platform access through web interface and Android application
- Device management and control with role-based access
- Secure communication with TLS encryption
- Containerized deployment with Docker Swarm

## System Architecture

The system architecture consists of five primary components working in concert:

* **Client Application**: A responsive React-based web interface that provides administrators and security personnel with an intuitive dashboard for monitoring, search operations, and system configuration.

* **Server Backend**: Powered by Go for maximum performance and concurrency, the server handles authentication, device management, image processing, and serves as the central coordination point for the entire system.

* **Android Application**: A native Kotlin mobile client that extends monitoring capabilities to on-the-go scenarios, allowing remote system control from the UI.

* **MQTT Broker**: Facilitates lightweight, real-time bidirectional communication between system components using the publish-subscribe pattern, enabling status updates with minimal latency.

* **Storage Layer**: Combines MongoDB's flexible document storage for metadata and event information with S3-compatible object storage for efficient management of captured images.

## Core Features

### Real-time Monitoring
- Live view of device feeds with mode switching
- Real-time status updates via MQTT protocol
- Instant device registration and disconnection detection

### Text-based Search
- OCR processing using Tesseract to extract text from images
- Search photos by extracted text content
- Advanced filtering by date, device, and text content

### Multi-platform Access
- Web interface for desktop administration
- Android application for mobile monitoring
- Consistent user experience across platforms

### Device Management
- Device registration and state tracking
- Mode switching between Normal and Live operation
- Role-based access control (Admin/User permissions)

### Security Features
- TLS encryption for all communications
- JWT-based authentication
- Certificate-based MQTT security
- Role-based access control

## Technology Stack

| Component | Technologies | Responsibilities |
|-----------|--------------|------------------|
| Client | React, TypeScript, TailwindCSS | User interface, photos search, device settings management |
| Server | Go, AWS SDK, MongoDB | Backend, coordination, authentication |
| Android App | Kotlin | Mobile interface, real-time and on-command photos |
| MQTT Broker | Mosquitto | Real-time communication between server and devices |

## Prerequisites

* Docker and Docker Compose
* Node.js (v16 or later)
* Go (v1.24 or later)
* Android Studio (for mobile development)
* MongoDB
* AWS S3 account (for file storage)
* MQTT Broker (Mosquitto)

## Deployment Architecture

The system is deployed using Docker containers orchestrated with Docker Swarm, allowing for easy scaling of individual components based on load requirements. The deployment can be customized for cloud environments or on-premises installation.

## Development Workflow

The project follows a modern development workflow with continuous integration and deployment:

* Feature branch development with pull request reviews
* Automated testing through GitHub Actions workflows
* Dependency management with Renovate (configured limit: 5 PRs)
* Automated Docker image builds on successful merges to main
* Automatic versioned deployment with rollback capabilities

## Security & Compliance

### Threat Modeling & Mitigations
- Static code analysis using CodeQL
- Automated vulnerability scanning
- Security policy enforcement

### Testing & Coverage
- Unit tests with >60% code coverage
- Automated CI/CD testing
- Security testing integration

### SBOM & Dependencies
- All dependencies pinned using hashes
- Automated SBOM generation with Syft
- OpenSSF vulnerability scanning

## Team Contributions

* **Mihnea-Andrei Blotiu**: 4,235 lines added, 6,303 lines removed, 53 commits
* **Stefan-Darius Iordache**: 5,710 lines added, 342 lines removed, 31 commits
* **Stefan-Dorin Jumarea**: 1,145 lines added, 171 lines removed, 10 commits
* **Roxana Popa**: 223 lines added, 0 lines removed, 3 commits

## OSSF Criticality Score

The project achieves a perfect OSSF security score of 10.0/10.0 through:
- No binary artifacts in repository
- No dangerous workflow patterns
- Dependency update tool (Renovate) integration
- Fuzzing with gofuzz
- MIT License compliance
- Pinned dependencies by hash
- SAST with CodeQL
- Security policy presence
- Readonly token permissions
- Zero vulnerabilities detected

## License

This project is licensed under the MIT License - see the [MIT-LICENSE.md](MIT-LICENSE.md) file for details.

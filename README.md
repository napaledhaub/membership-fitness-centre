# Fitness Center Membership Management

A simple CRUD API for managing memberships at a fitness center, built using Go (Golang) with PostgreSQL. The application includes features for email verification and a cron job for sending membership renewal notifications.

## Features

- **CRUD Operations**: Create, Read, Update, and Delete membership records.
- **Email Verification**: Send email with a verification token upon membership registration.
- **Membership Renewal Notifications**: Scheduled notifications to remind members of upcoming renewals.

## Tech Stack

- **Language**: Go (Golang)
- **Database**: PostgreSQL
- **Email**: [Gomail](https://github.com/go-gomail/gomail) (for sending emails)
- **Cron Job**: [robfig/cron](https://github.com/robfig/cron) (for scheduling)
- **Database Driver**: [pq](https://github.com/lib/pq) (PostgreSQL driver for Go)

## Getting Started

### Prerequisites

- Go (version 1.16 or later)
- PostgreSQL database
- SMTP server for sending emails

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/napaledhaub/membership-fitness-centre.git
   cd membership-fitness-centre

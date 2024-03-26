# Promotion System

This promotion system, written in Golang, is designed to support promotional campaigns on the System. It allows for the initiation of campaigns, tracking of user registrations, and issuance of discount vouchers to eligible users.

## Requirements

The promotion system should fulfill the following requirements:

1. **Campaign Initiation**: When a client registers a new account on the system, the system should initiate a promotion campaign.

2. **User Tracking**: Each campaign should support 100 first login users. The system should track user registrations and identify the first 100 users who log in after registering.

3. **Voucher Issuance**: The 100 first login users should receive a 30% discount voucher when they top-up their mobile phone's fee via the main app.

4. **Scalability**: The system should be able to handle at least 100,000 concurrent users when the promotion program is active.

## Components

The promotion system consists of the following components:

1. **Promotion Service**: Responsible for initiating and managing promotion campaigns.

2. **eKYC Service**: Tracks user registrations and logins.

3. **Voucher Service**: Manages the issuance and redemption of discount vouchers.

## Installation and Setup

### Prerequisites

Before running the promotion system, ensure that the following prerequisites are met:

- Go version 1.22.1 or higher is installed
- Docker is installed on your system

### Setup

To set up the promotion system, follow these steps:

1. Clone the repository

   ```bash
   git clone https://github.com/vdntruong/go-promotion.git
   ```

2. Start services:

   ```bash
   ./start_services.sh
   ```

3. Stop services:

   ```bash
   ./clean_services.sh
   ```

## Usage

Once the promotion system is set up and running, clients can register new accounts on the System. Clients can initiate promotion campaigns and issue discount vouchers to the first 100 users who log in after registration.

## Scalability Considerations

To ensure scalability, the promotion system architecture is designed to handle a large number of concurrent users. Each service is horizontally scalable and can be deployed across multiple instances to meet demand.

Optionally, set up an API Gateway (e.g., Nginx) to route requests to the appropriate services.

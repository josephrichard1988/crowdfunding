# Crowdfunding Platform v2 - Complete Documentation

> A blockchain-based crowdfunding platform built on **Hyperledger Fabric** enabling startups to raise funds from investors through a transparent, secure, and validator-verified process.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [System Architecture](#system-architecture)
3. [Project Structure](#project-structure)
4. [Network Setup - Step by Step](#network-setup---step-by-step)
5. [Server Configuration](#server-configuration)
6. [Completed Features](#completed-features)
7. [API Reference](#api-reference)

---

## Project Overview

### What This Platform Does

The Crowdfunding Platform v2 allows:

1. **Startups** to create and manage fundraising campaigns
2. **Validators** to assess campaign legitimacy and assign risk scores
3. **Platform Admins** to publish verified campaigns to the public portal
4. **Investors** to browse, evaluate, and invest in campaigns

### Key Innovation: Works on single channel named as crowdfunding-channel and uses Private Data Collections (PDC)

The platform uses **12 Private Data Collections** to ensure:

- ✅ Data privacy between organizations (startups can't see investor portfolios)
- ✅ Selective data sharing (validation scores shared only with relevant parties)
- ✅ Immutable audit trails on blockchain
- ✅ 3-way hash verification to prevent tampering

---

## System Architecture

### Network Topology

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Microfab Container (Port 9090)                    │
│                                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ StartupOrg  │  │ ValidatorOrg│  │ InvestorOrg │  │ PlatformOrg │         │
│  │   peer0     │  │   peer0     │  │   peer0     │  │   peer0     │         │
│  │  CouchDB    │  │  CouchDB    │  │  CouchDB    │  │  CouchDB    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                                              │
│                        ┌──────────────────────┐                              │
│                        │    Raft Orderer      │                              │
│                        └──────────────────────┘                              │
│                                                                              │
│                    Channel: crowdfunding-channel                             │
│                    Chaincode: crowdfunding (combined)                        │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Three-Tier Application Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              FRONTEND (React)                                │
│                             http://localhost:5173                            │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  StartupDashboard.jsx (44KB) │ ValidatorDashboard.jsx │ InvestorDash  │ │
│  │  PlatformDashboard.jsx       │ Wallet.jsx (40KB)      │ Login/Signup  │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │ HTTP REST
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           BACKEND SERVERS (Node.js)                          │
│  ┌─────────────────────────────────┐  ┌─────────────────────────────────┐   │
│  │     AUTH SERVER (client/)       │  │   FABRIC API SERVER (network/)  │   │
│  │     http://localhost:3001       │  │     http://localhost:4000       │   │
│  │                                 │  │                                 │   │
│  │  • User Registration/Login      │  │  • Fabric Gateway SDK           │   │
│  │  • JWT Authentication           │  │  • 30+ API Endpoints            │   │
│  │  • MongoDB Integration          │  │  • Auto-ID Generation           │   │
│  │  • User Profile Management      │  │  • Cross-Org Transactions       │   │
│  │  • Campaign Allocation Queue    │  │                                 │   │
│  └────────────────────┬────────────┘  └────────────────┬────────────────┘   │
│                       │                                 │                    │
│                       ▼                                 ▼                    │
│               ┌───────────────┐               ┌──────────────────┐          │
│               │   MongoDB     │               │  Hyperledger     │          │
│               │ (Off-chain)   │               │  Fabric Network  │          │
│               └───────────────┘               └──────────────────┘          │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Private Data Collections

| Collection | Access | Purpose |
|------------|--------|---------|
| `StartupPrivateData` | StartupOrg only | Draft campaigns, internal notes |
| `InvestorPrivateData` | InvestorOrg only | Portfolio, investment tracking |
| `ValidatorPrivateData` | ValidatorOrg only | Assessment details |
| `PlatformPrivateData` | PlatformOrg only | Wallets, fee structures |
| `StartupValidatorShared` | Startup + Validator | Validation submissions |
| `StartupInvestorShared` | Startup + Investor | Investment acknowledgments |
| `StartupPlatformShared` | Startup + Platform | Approved campaigns for publishing |
| `InvestorValidatorShared` | Investor + Validator | Validation score requests |
| `ValidatorPlatformShared` | Validator + Platform | Validation proof hashes |
| `InvestorPlatformShared` | Investor + Platform | Escrow agreements |
| `AllOrgsShared` | All 4 organizations | Disputes (transparency) |
| `TokenLedger` | Platform controls | CFT token balances |

---

## Project Structure

```
crowdfundingv2/
│
├── contracts/                          # 📦 CHAINCODE (Go)
│   ├── main.go                        # Entry point - registers all contracts
│   ├── shared_types.go                # Shared data structures (Campaign, Startup, etc.)
│   ├── startuporg_contract.go         # StartupContract - 2065 lines
│   ├── investororg_contract.go        # InvestorContract - 1500+ lines
│   ├── validatororg_contract.go       # ValidatorContract - 1400+ lines
│   ├── platformorg_contract.go        # PlatformContract - 2000+ lines
│   ├── token_operations.go            # TokenContract - 1300+ lines
│   └── go.mod, go.sum                 # Go dependencies
│
├── network/                            # 🌐 FABRIC API SERVER (Node.js)
│   └── src/
│       ├── index.js                   # EXPRESS SERVER - 903 lines, 30+ endpoints
│       ├── fabricConnection.js        # Fabric Gateway SDK integration
│       ├── enrollAdmin.js             # Admin enrollment utility
│       └── config/                    # Server configuration
│
├── client/                             # 🔐 AUTH SERVER (Node.js)
│   └── src/
│       ├── index.js                   # EXPRESS SERVER for authentication
│       ├── routes/                    # API routes (auth, startup, validator, etc.)
│       ├── controllers/               # Business logic
│       ├── models/                    # MongoDB schemas
│       ├── database/db.js             # MongoDB connection
│       └── middleware/                # JWT verification
│
├── frontend/                           # 🖥️ REACT FRONTEND
│   └── src/
│       ├── App.jsx                    # Main app with routing
│       ├── pages/                     # 11 page components
│       │   ├── StartupDashboard.jsx   # Campaign management (44KB)
│       │   ├── ValidatorDashboard.jsx # Validation queue (32KB)
│       │   ├── InvestorDashboard.jsx  # Investment portfolio (16KB)
│       │   ├── PlatformDashboard.jsx  # Admin panel (20KB)
│       │   ├── Wallet.jsx             # Token management (40KB)
│       │   ├── Login.jsx, Signup.jsx  # Authentication
│       │   └── ...more pages
│       ├── services/api.js            # Axios HTTP client
│       └── context/                   # React Context (Auth, Theme)
│
├── deploy_chaincode.sh                 # 📜 DEPLOYMENT AUTOMATION - 872 lines
├── collections_config.json             # PDC configuration (12 collections)
├── MICROFAB.txt                        # Network topology config
│
├── _msp/                               # Generated MSP certificates
├── _wallets/                           # Generated user identities
├── _gateways/                          # Generated connection profiles
└── bin/                                # Fabric binaries (peer, orderer)
```

---

## Network Setup - Step by Step

> [!IMPORTANT]
> Follow these steps IN ORDER. Each step depends on the previous one completing successfully.

### Prerequisites

- **Docker** 20.10+ installed
- **Node.js** 18+ installed
- **Go** 1.21+ installed
- **weft** CLI tool: `npm install -g @hyperledger-labs/weft`

---

### STEP 1: Start Microfab Network (Terminal 1)

**What this does:** Launches a Docker container with a complete Hyperledger Fabric network including 4 organizations, their peers, and an orderer.

```bash
# Navigate to project directory
cd /home/quantum_pulse/crowdfunding/crowdfundingv2
# cd /home/give_your_path/crowdfunding/crowdfundingv2

# Export the network configuration
export MICROFAB_CONFIG=$(cat MICROFAB.txt)

# Start the Microfab container (THIS WILL BLOCK THIS TERMINAL)
docker run --name microfab -e MICROFAB_CONFIG -p 9090:9090 ibmcom/ibp-microfab
```

**What happens:**

- Docker pulls and runs the Microfab image
- Creates 4 organizations: StartupOrg, ValidatorOrg, InvestorOrg, PlatformOrg
- Each org gets 1 peer with CouchDB state database
- Creates channel `crowdfunding-channel`
- Exposes API on port 9090

**Expected output:**

```
[INFO] Microfab started on port 9090
[INFO] Console: http://console.127-0-0-1.nip.io:9090
```

> [!TIP]
> Keep this terminal running! Open a **NEW TERMINAL** for the next steps.

---

### STEP 2: Generate Connection Profiles (Terminal 2)

**What this does:** Downloads certificates and creates connection profiles from Microfab so your applications can connect to the network.

```bash
# Open a NEW terminal and navigate to project
cd /home/quantum_pulse/crowdfunding/crowdfundingv2

# Generate wallets, gateways, and MSP configs
curl -s http://console.127-0-0-1.nip.io:9090/ak/api/v1/components | \
  weft microfab -w ./_wallets -p ./_gateways -m ./_msp -f
```

**What happens:**

- Fetches component information from running Microfab
- Creates `_wallets/` folder with user identity credentials
- Creates `_gateways/` folder with connection profiles for each org
- Creates `_msp/` folder with MSP (Membership Service Provider) configs

**Generated structure:**

```
_wallets/
├── StartupOrg/
│   └── startuporgadmin.id
├── ValidatorOrg/
│   └── validatororgadmin.id
├── InvestorOrg/
│   └── investororgadmin.id
└── PlatformOrg/
    └── platformorgadmin.id

_gateways/
├── StartupOrg/startuporg_gateway.json
├── ValidatorOrg/validatororg_gateway.json
├── InvestorOrg/investororg_gateway.json
└── PlatformOrg/platformorg_gateway.json

_msp/
├── StartupOrg/startuporgadmin/msp/
├── ValidatorOrg/validatororgadmin/msp/
├── InvestorOrg/investororgadmin/msp/
└── PlatformOrg/platformorgadmin/msp/
```

---

### STEP 3: Install Fabric Binaries

**What this does:** Downloads the Hyperledger Fabric command-line tools (`peer`, `orderer`, `configtxgen`) needed for chaincode operations.

```bash
# Download and install Fabric binaries
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh | bash -s -- binary
```

**What happens:**

- Downloads `peer` binary (for chaincode and channel operations)
- Downloads `orderer` binary (for ordering service operations)
- Downloads `configtxgen` (for channel configuration)
- Places binaries in `./bin/` directory

---

### STEP 4: Update MSP Path in deploy_chaincode.sh

**What this does:** Configures the deployment script with YOUR system's path.

```bash
# Open deploy_chaincode.sh and update line 140
nano deploy_chaincode.sh

# Change this line:
MSP_BASE_PATH="/home/give_your_path/crowdfunding/crowdfundingv2/_msp"

# To YOUR path:
MSP_BASE_PATH="/home/give_your_path/crowdfunding/crowdfundingv2/_msp"

# Save and exit (Ctrl+X, Y, Enter)
```

---

### STEP 5: Package Chaincode

**What this does:** Compiles all Go contracts into a single deployable package.

```bash
source ./deploy_chaincode.sh package
```

**What happens internally:**

1. Navigates to `contracts/` directory
2. Runs `go mod tidy` to resolve dependencies
3. Runs `go mod vendor` to create vendor directory
4. Detects existing packages (crowdfunding_1,_2, etc.)
5. Creates new package with incremented version
6. Packages all `.go` files into `crowdfunding_X.tar.gz`

**Output:**

```
📦 Packaging Combined Chaincode
[INFO] Running go mod tidy...
[INFO] Running go mod vendor...
[INFO] New version: 1
[INFO] Package file: crowdfunding_1.tar.gz
✅ Successfully packaged crowdfunding as crowdfunding_1.tar.gz
```

---

### STEP 6: Install Chaincode on All Organizations

**What this does:** Installs the chaincode package on all 4 organization peers.

```bash
source ./deploy_chaincode.sh install
```

**What happens internally:**

1. Finds latest `.tar.gz` package
2. For each organization (Startup, Investor, Validator, Platform):
   - Switches peer context (sets CORE_PEER_LOCALMSPID, CORE_PEER_ADDRESS)
   - Runs `peer lifecycle chaincode install crowdfunding_1.tar.gz`
   - Extracts and exports `PACKAGE_ID`

**Output:**

```
📥 Installing Combined Chaincode on All Organizations
📦 Installing on StartupOrg
Press Enter to install on StartupOrg...
[INFO] Switching to StartupOrg...
✅ Now operating as StartupOrg
Chaincode code package identifier: crowdfunding_1:abc123def456...
✅ Auto-exported for StartupOrg: PACKAGE_ID=crowdfunding_1:abc123def456...

📦 Installing on InvestorOrg
... (repeats for all 4 orgs)
```

> [!IMPORTANT]
> After installation, the `PACKAGE_ID` environment variable is set automatically. If you close the terminal, you'll need to re-export it.

---

### STEP 7: Deploy Chaincode (Approve + Commit)

**What this does:** Gets all organizations to approve the chaincode and commits it to the channel.

```bash
source ./deploy_chaincode.sh deploy
```

**What happens internally:**

**Phase 1: Approval (for each org):**

```bash
peer lifecycle chaincode approveformyorg \
  --channelID crowdfunding-channel \
  --name crowdfunding \
  --version 1 \
  --package-id $PACKAGE_ID \
  --sequence 1 \
  --collections-config ./collections_config.json \  # REQUIRED for PDC
  --signature-policy "OR('StartupOrgMSP.peer','InvestorOrgMSP.peer',...)"
```

**Phase 2: Commit (once):**

```bash
peer lifecycle chaincode commit \
  --channelID crowdfunding-channel \
  --name crowdfunding \
  --version 1 \
  --sequence 1 \
  --collections-config ./collections_config.json \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 \
  --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 \
  --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090
```

**Output:**

```
✅ Approving Chaincode on All Organizations
📋 Deploying on StartupOrg
... (approval for each org)

🚀 Committing Chaincode to Channel
✅ Successfully committed crowdfunding (version: 1, sequence: 1)
🎉 DEPLOYMENT COMPLETE! 🎉
```

---

### STEP 8: Verify Deployment

```bash
# Query committed chaincode
source ./deploy_chaincode.sh query-committed

# Expected output:
# Committed chaincode definition for chaincode 'crowdfunding' on channel 'crowdfunding-channel':
# Version: 1, Sequence: 1, Endorsement Plugin: escc...
```

---

### STEP 9: Switch Organization Context

**What this does:** Changes which organization's peer you're operating as.

```bash
# Switch to StartupOrg
source ./deploy_chaincode.sh switch startup

# Switch to ValidatorOrg  
source ./deploy_chaincode.sh switch validator

# Switch to InvestorOrg
source ./deploy_chaincode.sh switch investor

# Switch to PlatformOrg
source ./deploy_chaincode.sh switch platform
```

**What happens:**
Sets these environment variables:

- `CORE_PEER_LOCALMSPID` = Organization MSP ID
- `CORE_PEER_MSPCONFIGPATH` = Path to admin MSP certificates
- `CORE_PEER_ADDRESS` = Peer API endpoint
- `PATH` and `FABRIC_CFG_PATH` for Fabric tools

---

## Server Configuration

### How to Start All Services

> [!WARNING]
> Start servers in this order: Microfab → Auth Server → Fabric API → Frontend

#### Terminal 1: Microfab (Already running from Step 1)

```bash
# Should already be running
docker ps  # Verify microfab container is running
```

#### Terminal 2: Auth Server (Authentication + MongoDB)

```bash
cd /home/give_your_path/crowdfunding/crowdfundingv2/client

# Install dependencies (first time only)
npm install

# Configure environment
cp .env.template .env
# Edit .env with your MongoDB connection string

# Start server
npm start
```

**Port:** `3001`
**Endpoints:**

- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login (returns JWT)
- `GET /api/auth/profile` - Get user profile
- `POST /api/auth/startups` - Sync startup to MongoDB
- `POST /api/auth/sync/campaign` - Sync campaign metadata
- `POST /api/auth/sync/campaign-status` - Update campaign status

**Environment Variables (`client/.env`):**

```bash
PORT=3001
MONGODB_URI=mongodb://make_your_connection_string
JWT_SECRET=your-secret-key
NODE_ENV=development
```

#### Terminal 3: Fabric API Server (Chaincode Gateway)

```bash
cd /home/give_your_path/crowdfunding/crowdfundingv2/network

# Install dependencies (first time only)
npm install

# Start server
npm start
```

**Port:** `4000`
**What this server does:**

- Connects to Hyperledger Fabric using Gateway SDK
- Exposes 30+ REST endpoints for chaincode functions
- Auto-generates campaign/startup IDs
- Syncs blockchain state with MongoDB

**Key Endpoints:**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/startup/startups` | Create new startup |
| POST | `/api/startup/campaigns` | Create campaign (22 params) |
| GET | `/api/startup/campaigns/:id` | Get campaign details |
| POST | `/api/startup/campaigns/:id/submit-validation` | Submit for validation |
| POST | `/api/startup/campaigns/:id/share-to-platform` | Share approved campaign |
| DELETE | `/api/startup/campaigns/:id` | Delete campaign (with fee) |
| GET | `/api/validator/pending-validations` | Get validation queue |
| POST | `/api/validator/approve/:id` | Approve/reject campaign |
| GET | `/api/platform/shared-campaigns` | Get campaigns pending publish |
| POST | `/api/platform/publish/:id` | Publish to portal |
| GET | `/api/investor/campaigns` | Browse published campaigns |
| POST | `/api/investor/investments` | Make investment |

#### Terminal 4: Frontend (React)

```bash
cd /home/give_your_path/crowdfunding/crowdfundingv2/frontend

# Install dependencies (first time only)
npm install

# Start development server
npm run dev
```

**Port:** `5173`
**URL:** <http://localhost:5173>

**Available Routes:**

| Route | Component | Description |
|-------|-----------|-------------|
| `/` | Dashboard | Landing page |
| `/login` | Login | User authentication |
| `/signup` | Signup | User registration |
| `/startup/dashboard` | StartupDashboard | Manage campaigns |
| `/startup/detail/:id` | StartupDetail | Startup details |
| `/startup/campaign/:id` | StartupCampaignDetails | Campaign details |
| `/validator/dashboard` | ValidatorDashboard | Validation queue |
| `/investor/dashboard` | InvestorDashboard | View/invest campaigns |
| `/platform/dashboard` | PlatformDashboard | Admin: publish campaigns |
| `/wallet` | Wallet | CFT token management |

---

### Quick Start Summary

```bash
# Terminal 1: Start Microfab
cd ~/crowdfunding/crowdfundingv2
export MICROFAB_CONFIG=$(cat MICROFAB.txt)
docker run --name microfab -e MICROFAB_CONFIG -p 9090:9090 ibmcom/ibp-microfab

# Terminal 2: Generate profiles + Deploy chaincode
cd ~/crowdfunding/crowdfundingv2
curl -s http://console.127-0-0-1.nip.io:9090/ak/api/v1/components | weft microfab -w ./_wallets -p ./_gateways -m ./_msp -f
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh | bash -s -- binary
source ./deploy_chaincode.sh package
source ./deploy_chaincode.sh install
source ./deploy_chaincode.sh deploy

# Terminal 3: Auth Server
cd ~/crowdfunding/crowdfundingv2/client && npm start
# Runs on port 3001

# Terminal 4: Fabric API Server  
cd ~/crowdfunding/crowdfundingv2/network && npm start
# Runs on port 4000

# Terminal 5: Frontend
cd ~/crowdfunding/crowdfundingv2/frontend && npm run dev
# Runs on port 5173 → http://localhost:5173
```

---

## Completed Features

### 1. Startup Management (StartupContract)

**Functions Implemented:**

| Function | Parameters | What It Does |
|----------|------------|--------------|
| `CreateStartup` | startupID, ownerID, name, description, displayID | Creates startup entity in private collection |
| `GetStartup` | startupID | Retrieves startup with all campaigns |
| `GetStartupsByOwner` | ownerID | Lists all startups owned by user |
| `CreateCampaign` | 22 parameters | Creates fundraising campaign |
| `SubmitForValidation` | campaignID, documents, notes | Submits to validator queue |
| `ShareCampaignToPlatform` | campaignID, validationHash | Sends approved campaign to platform |
| `DeleteCampaign` | campaignID, reason | Deletes with 60% fee or 100 CFT fixed |
| `DeleteStartup` | startupID, reason | Cascade deletes all campaigns |
| `CalculateCampaignDeletionFee` | campaignID | Returns fee preview |

**Campaign 22-Parameter Format:**

```json
{
  "campaignID": "CAMP_STU_xxx_001",
  "startupID": "STU_user123_001",
  "category": "Technology",
  "deadline": "2025-03-31",
  "currency": "USD",
  "has_raised": false,
  "has_gov_grants": false,
  "incorp_date": "2025-01-01",
  "project_stage": "Prototype",
  "sector": "Hardware",
  "tags": ["IoT", "SmartHome", "AI"],
  "team_available": true,
  "investor_committed": false,
  "duration": 90,
  "funding_day": 1,
  "funding_month": 1,
  "funding_year": 2025,
  "goal_amount": 50000,
  "investment_range": "50K-100K",
  "project_name": "Smart Home IoT Platform",
  "description": "Full description...",
  "documents": ["business_plan.pdf", "pitch_deck.pdf"]
}
```

### 2. Validator Operations (ValidatorContract)

| Function | What It Does |
|----------|--------------|
| `GetPendingValidations` | Returns campaigns awaiting validation |
| `GetCampaign` | Get campaign from StartupValidatorShared |
| `ValidateCampaign` | Start validation process |
| `ApproveOrRejectCampaign` | Score and approve/reject with digital signature |
| `ProvideValidationDetailsToInvestor` | Share scores with requesting investor |
| `VerifyMilestoneCompletion` | Verify milestone evidence |
| `GetAllValidations` | Validation history |

**Scoring System:**

- Due Diligence Score: 0-10 (business viability)
- Risk Score: 0-10 (investment risk)
- Risk Level: LOW | MEDIUM | HIGH

### 3. Platform Management (PlatformContract)

| Function | What It Does |
|----------|--------------|
| `GetAllSharedCampaigns` | Campaigns shared by startups |
| `PublishCampaignToPortal` | Verify hash + publish to public ledger |
| `CreateWallet` | Create CFT token wallet |
| `TransferTokens` | Move CFT between wallets |
| `SetCampaignFeeTier` | Configure success fees |
| `TriggerFundRelease` | Release escrow funds on milestone |
| `CreateDispute` | File disputes |
| `ResolveDispute` | Arbitrate with fund redistribution |

### 4. Investor Operations (InvestorContract)

| Function | What It Does |
|----------|--------------|
| `GetAvailableCampaigns` | Browse published campaigns |
| `ViewCampaignDetails` | Full campaign details |
| `RequestValidationDetails` | Request scores from validator |
| `MakeInvestment` | Direct investment |
| `CreateInvestmentProposal` | Propose terms with milestones |
| `GetMyInvestments` | Portfolio tracking |
| `GetViewedCampaigns` | Recently viewed |

### 5. Token System (TokenContract)

| Function | What It Does |
|----------|--------------|
| `InitializeToken` | Deploy CFT token system |
| `Mint` | Create new tokens (Platform only) |
| `Transfer` | Send CFT between wallets |
| `Burn` | Remove tokens from circulation |
| `GetBalance` | Check wallet balance |
| `GetTotalSupply` | Total CFT in circulation |

---

## API Reference

### Using Chaincode CLI

```bash
# Switch to StartupOrg
source ./deploy_chaincode.sh switch startup

# Create a startup
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"StartupContract:CreateStartup","Args":["STARTUP001","STU_user123","TechVentures Inc","An innovative startup","S-001"]}'

# Create a campaign  
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"StartupContract:CreateCampaign","Args":["CAMP001","STARTUP001","Technology","2025-03-31","USD","false","false","2025-01-01","Prototype","Hardware","[\"IoT\"]","false","false","90","1","1","2025","50000","50K-100K","IoT Platform","Description","[\"docs.pdf\"]"]}'

# Query campaign
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  -c '{"function":"StartupContract:GetCampaign","Args":["CAMP001"]}'
```

### Using REST API (Recommended)

```bash
# Create startup via API (auto-generates ID)
curl -X POST http://localhost:4000/api/startup/startups \
  -H "Content-Type: application/json" \
  -d '{
    "name": "TechVentures Inc",
    "description": "An innovative startup",
    "ownerId": "STU_user123"
  }'

# Create campaign via API (auto-generates ID)
curl -X POST http://localhost:4000/api/startup/campaigns \
  -H "Content-Type: application/json" \
  -d '{
    "startupId": "STU_user123_001",
    "projectName": "IoT Platform",
    "category": "Technology",
    "goalAmount": 50000,
    ...
  }'
```

---

## Upgrading Chaincode

After making code changes in `contracts/`:

```bash
# Interactive upgrade (recommended)
source ./deploy_chaincode.sh upgrade

# Prompts:
# 1. Do you want to repackage? (y/n) → y
# 2. Do you want to reinstall? (y/n) → y  
# 3. Do you want to deploy? (y/n) → y

# This will:
# - Create crowdfunding_2.tar.gz (auto-increment)
# - Install on all 4 orgs
# - Approve with sequence: 2
# - Commit upgrade to channel
```

---

## Summary

✅ **9000+ lines of Go chaincode** across 5 contracts
✅ **12 Private Data Collections** for granular privacy
✅ **872-line deployment script** with full automation
✅ **2 Node.js backend servers** (Auth + Fabric API)
✅ **React frontend** with 11 page components
✅ **CFT token system** for platform economy
✅ **Complete E2E workflow** from campaign to investment

---

## Current Implementation Status

### MongoDB Synchronization

| Organization | Frontend Integration | MongoDB Sync | Status |
|--------------|---------------------|--------------|--------|
| **StartupOrg** | ✅ Complete | ✅ Complete | Fully functional |
| **ValidatorOrg** | ✅ Complete | ✅ Complete | Fully functional |
| **PlatformOrg** | ✅ Frontend done | ⏳ Pending | MongoDB sync remaining |
| **InvestorOrg** | ✅ Frontend done | ⏳ Pending | MongoDB sync remaining |

> [!NOTE]
> **Startup** and **Validator** client-side integration with MongoDB synchronization is successfully completed. **Platform** and **Investor** organizations have their frontend pages implemented but the complete MongoDB synchronization for these roles is still remaining.

---

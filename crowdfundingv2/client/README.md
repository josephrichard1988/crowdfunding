# Authentication Server

This service handles **user authentication ONLY** for the crowdfunding platform.

## Purpose

- User signup/login
- JWT token generation
- User profile management
- Wallet balance tracking (synced from Fabric)

## Port

**3001** (default)

## What This Server Does NOT Do

❌ This server does **NOT** interact with Hyperledger Fabric
❌ This server does **NOT** call chaincode
❌ This server does **NOT** handle business logic

## Architecture

```
Frontend (5173)
    ↓ (auth requests)
Auth Server (3001) ← You are here
    ↓ (MongoDB)
Database

Frontend (5173)
    ↓ (fabric API requests)
Network Server (4000)
    ↓ (Fabric SDK)
Fabric Network (9090)
```

## Environment Variables

Create a `.env` file:

```env
PORT=3001
MONGODB_URI=mongodb://localhost:27017/crowdfunding
JWT_SECRET=your-secret-key
NODE_ENV=development

# Fee Configuration
INR_TO_CFT_RATE=2.5
USD_TO_CFT_RATE=83.0
REGISTRATION_FEE_CFT=250
CAMPAIGN_CREATION_FEE_CFT=1250
```

## API Endpoints

- `POST /api/auth/signup` - Create new user
- `POST /api/auth/login` - Login and get JWT token
- `GET /api/auth/me` - Get current user (requires auth)
- `PUT /api/auth/wallet` - Update wallet balance
- `GET /api/auth/fees` - Get fee schedule

## Run

```bash
npm install
npm run dev
```

# 🚀 Crowdfunding Platform - Chaincode Deployment Flow Guide

## 📋 Table of Contents
1. [Available Commands](#available-commands)
2. [Initial Deployment Flow](#initial-deployment-flow)
3. [Upgrade Flow](#upgrade-flow)
4. [Expected Output at Each Step](#expected-output-at-each-step)
5. [Troubleshooting](#troubleshooting)
6. [Quick Reference](#quick-reference)

---

## 📌 Available Commands

### **Packaging Commands**
```bash
./deploy_chaincode.sh package                    # Package all 4 chaincodes (auto-version)
./deploy_chaincode.sh package <chaincode>        # Package specific chaincode
# Chaincodes: startup | validator | investor | platform
```

### **Installation Commands**
```bash
./deploy_chaincode.sh install <org>              # Install on specific org
./deploy_chaincode.sh install-all                # Install on all 4 orgs (interactive)
# Organizations: startup | validator | investor | platform
```

### **Deployment Commands (Initial)**
```bash
./deploy_chaincode.sh deploy <org>               # Deploy all channels for org
./deploy_chaincode.sh deploy <org> <channel>     # Deploy specific channel only
./deploy_chaincode.sh deploy all                 # Deploy for all orgs (interactive)
# Channels: common | startup-validator | startup-investor | startup-platform |
#           investor-validator | investor-platform | validator-platform
```

### **Upgrade Commands**
```bash
./deploy_chaincode.sh upgrade <cc> <channel>     # Upgrade specific chaincode
./deploy_chaincode.sh upgrade-all                # Upgrade all chaincodes on all channels
./deploy_chaincode.sh sync-upgrade <cc> <ch> <org> # Smart upgrade for specific org
./deploy_chaincode.sh approve-chaincode <org> <cc> <ch> # Manually approve upgrade
# Use approve-chaincode when another org already upgraded and you need to approve
```

### **Query & Check Commands**
```bash
./deploy_chaincode.sh query-committed <cc> <ch>  # Query committed details
./deploy_chaincode.sh check-readiness <cc> <ch>  # Check commit readiness
```

### **Utility Commands**
```bash
./deploy_chaincode.sh switch <org>               # Switch peer context
./deploy_chaincode.sh help                       # Show detailed help
```

---

## 🎯 Initial Deployment Flow

### **STEP 1: Package All Chaincodes** ⏱️ ~10 seconds

**Command:**
```bash
./deploy_chaincode.sh package
```

**What Happens:**
- Script scans for existing package versions (startup_*.tgz)
- Auto-increments version numbers (if startup_0.tgz exists → creates startup_1.tgz)
- Packages Go chaincode from `./contracts/startuporg/`, `./contracts/validatororg/`, etc.
- Creates 4 package files in current directory

**Expected Output:**
```
📦 Packaging All Chaincodes
============================================================================
📦 Packaging chaincode: startup
   Current version: 0
   New version: 1
   Package label: startup_1
   Package file: startup_1.tgz
   Path: ./contracts/startuporg
✅ Successfully packaged startup as startup_1.tgz

📦 Packaging chaincode: validator
   Current version: 0
   New version: 1
   Package label: validator_1
   Package file: validator_1.tgz
   Path: ./contracts/validatororg
✅ Successfully packaged validator as validator_1.tgz

📦 Packaging chaincode: platform
   Current version: 0
   New version: 1
   Package label: platform_1
   Package file: platform_1.tgz
   Path: ./contracts/platformorg
✅ Successfully packaged platform as platform_1.tgz

📦 Packaging chaincode: investor
   Current version: 0
   New version: 1
   Package label: investor_1
   Package file: investor_1.tgz
   Path: ./contracts/investororg
✅ Successfully packaged investor as investor_1.tgz

✅ All chaincodes packaged successfully!

[INFO] Package files created:
  - startup_1.tgz
  - validator_1.tgz
  - platform_1.tgz
  - investor_1.tgz
```

**Result Files:**
```
crowdfunding/
├── startup_1.tgz       ✅ Created
├── validator_1.tgz     ✅ Created
├── investor_1.tgz      ✅ Created
└── platform_1.tgz      ✅ Created
```

---

### **STEP 2: Install Chaincodes on All Organizations** ⏱️ ~2-3 minutes

**Command:**
```bash
./deploy_chaincode.sh install-all
```

**OR Install Individually:**
```bash
./deploy_chaincode.sh install startup
./deploy_chaincode.sh install validator
./deploy_chaincode.sh install investor
./deploy_chaincode.sh install platform
```

**What Happens:**
1. Switches to first org (StartupOrg)
2. Installs all 4 chaincode packages on StartupOrg peer
3. Auto-extracts package ID from installation output
4. Exports package ID as environment variable
5. Prompts user to press Enter before continuing
6. Repeats for ValidatorOrg, InvestorOrg, PlatformOrg

**Expected Output (StartupOrg example):**
```
============================================================================
 Installing Chaincodes on StartupOrg
============================================================================
🔄 Switching to StartupOrg...
[INFO] Fabric environment paths configured
✅ Now operating as StartupOrg

============================================================================
 Installing startup chaincode
============================================================================
[INFO] Installing startup_1.tgz on current peer...

2024.12.12 10:30:15.123 UTC 0001 INFO [cli.lifecycle.chaincode] submitInstallProposal -> Installed remotely: response:<status:200 payload:"\nKstartup_1:abc123def456789abc123def456789abc123def456789abc123def456789abc" >
2024.12.12 10:30:15.234 UTC 0001 INFO [cli.lifecycle.chaincode] submitInstallProposal -> Chaincode code package identifier: startup_1:abc123def456789abc123def456789abc123def456789abc123def456789abc

✅ Successfully installed startup on StartupOrg!
[INFO] Package ID: startup_1:abc123def456789abc123def456789abc123def456789abc123def456789abc

✅ Auto-exported: STARTUP_CC_PACKAGE_ID=startup_1:abc123def456789abc123def456789abc123def456789abc123def456789abc

[INFO] You can also manually export with:
export STARTUP_CC_PACKAGE_ID=startup_1:abc123def456789abc123def456789abc123def456789abc123def456789abc

Press Enter to continue to next chaincode...
```

**After STEP 2, these environment variables are set:**
```bash
✅ STARTUP_CC_PACKAGE_ID=startup_1:abc123...
✅ VALIDATOR_CC_PACKAGE_ID=validator_1:def456...
✅ INVESTOR_CC_PACKAGE_ID=investor_1:ghi789...
✅ PLATFORM_CC_PACKAGE_ID=platform_1:jkl012...
```

**If Installation Already Exists:**
```
[WARNING] startup already installed
[INFO] Package ID: startup_1:abc123def456789abc123def456789abc123def456789abc
✅ Auto-exported: STARTUP_CC_PACKAGE_ID=startup_1:abc123...
[INFO] You can also manually export with:
export STARTUP_CC_PACKAGE_ID=startup_1:abc123...

Press Enter to continue to next chaincode...
```

---

### **STEP 3: Deploy on Each Organization** ⏱️ ~5-7 minutes total

#### **3.1 Deploy StartupOrg**

**Command:**
```bash
./deploy_chaincode.sh deploy startup
```

**What Happens:**
- StartupOrg approves its own chaincode (startup) on 4 channels:
  - `startup-validator-channel`
  - `startup-platform-channel`
  - `startup-investor-channel`
  - `common-channel`
- StartupOrg **commits** startup chaincode on all 4 channels (as owner)
- StartupOrg approves other orgs' chaincodes (validator, investor, platform) on shared channels

**Expected Output:**
```
============================================================================
 Deploying chaincodes for StartupOrg
============================================================================
🔄 Switching to StartupOrg...
✅ Now operating as StartupOrg

============================================================================
 StartupOrg: startup-validator-channel
============================================================================
🔍 Checking if startup is already committed on startup-validator-channel...
ℹ️  startup not yet committed on startup-validator-channel. Will use version: 1, sequence: 1

[INFO] Approving chaincode 'startup' on channel 'startup-validator-channel' (version: 1, sequence: 1)...
2024.12.12 10:40:10.111 UTC [chaincodeCmd] ClientWait -> INFO 001 txid [abc123...] committed with status (VALID) at startuporgpeer-api.127-0-0-1.nip.io:9090
✅ Successfully approved 'startup' on 'startup-validator-channel' (version: 1, sequence: 1)

[INFO] Committing chaincode 'startup' on channel 'startup-validator-channel' (version: 1, sequence: 1)...
2024.12.12 10:40:15.222 UTC [chaincodeCmd] ClientWait -> INFO 001 txid [def456...] committed with status (VALID)
✅ Successfully committed 'startup' on 'startup-validator-channel' (version: 1, sequence: 1)

[INFO] Approving chaincode 'validator' on channel 'startup-validator-channel' (version: 1, sequence: 1)...
✅ Successfully approved 'validator' on 'startup-validator-channel' (version: 1, sequence: 1)

============================================================================
 StartupOrg: startup-platform-channel
============================================================================
[Similar output for startup-platform-channel]

============================================================================
 StartupOrg: startup-investor-channel
============================================================================
[Similar output for startup-investor-channel]

============================================================================
 StartupOrg: common-channel
============================================================================
[INFO] Approving chaincode 'startup' on channel 'common-channel' (version: 1, sequence: 1)...
✅ Successfully approved 'startup' on 'common-channel' (version: 1, sequence: 1)

[INFO] Committing chaincode 'startup' on channel 'common-channel' (version: 1, sequence: 1)...
✅ Successfully committed 'startup' on 'common-channel' (version: 1, sequence: 1)

[INFO] Approving chaincode 'validator' on channel 'common-channel' (version: 1, sequence: 1)...
✅ Successfully approved 'validator' on 'common-channel' (version: 1, sequence: 1)

[INFO] Approving chaincode 'investor' on channel 'common-channel' (version: 1, sequence: 1)...
✅ Successfully approved 'investor' on 'common-channel' (version: 1, sequence: 1)

[INFO] Approving chaincode 'platform' on channel 'common-channel' (version: 1, sequence: 1)...
✅ Successfully approved 'platform' on 'common-channel' (version: 1, sequence: 1)

✅ StartupOrg deployment complete!
```

#### **3.2 Deploy ValidatorOrg**

**Command:**
```bash
./deploy_chaincode.sh deploy validator
```

**What Happens:**
- ValidatorOrg approves its own chaincode (validator) on 4 channels:
  - `startup-validator-channel`
  - `investor-validator-channel`
  - `validator-platform-channel`
  - `common-channel`
- ValidatorOrg **commits** validator chaincode on all 4 channels (as owner)
- ValidatorOrg approves other orgs' chaincodes on shared channels

**Expected Output:** (Similar structure to StartupOrg, but for ValidatorOrg channels)

#### **3.3 Deploy InvestorOrg**

**Command:**
```bash
./deploy_chaincode.sh deploy investor
```

**What Happens:**
- InvestorOrg approves its own chaincode (investor) on 4 channels:
  - `startup-investor-channel`
  - `investor-validator-channel`
  - `investor-platform-channel`
  - `common-channel`
- InvestorOrg **commits** investor chaincode on all 4 channels (as owner)

#### **3.4 Deploy PlatformOrg**

**Command:**
```bash
./deploy_chaincode.sh deploy platform
```

**What Happens:**
- PlatformOrg approves its own chaincode (platform) on 4 channels:
  - `startup-platform-channel`
  - `investor-platform-channel`
  - `validator-platform-channel`
  - `common-channel`
- PlatformOrg **commits** platform chaincode on all 4 channels (as owner)

**Final Result After STEP 3:**
```
✅ All 4 chaincodes deployed on all 7 channels
✅ All chaincodes are Version 1, Sequence 1
✅ Ready for chaincode invocation
```

---

### **STEP 4: Verify Deployment** ⏱️ ~5 seconds per query

**Command:**
```bash
./deploy_chaincode.sh query-committed startup common-channel
```

**Expected Output:**
```
[INFO] Querying committed chaincode 'startup' on 'common-channel'...

Committed chaincode definition for chaincode 'startup' on channel 'common-channel':
Version: 1, Sequence: 1, Endorsement Plugin: escc, Validation Plugin: vscc, Approvals: [StartupOrgMSP: true, ValidatorOrgMSP: true, InvestorOrgMSP: true, PlatformOrgMSP: true]
```

**Verify Each Chaincode:**
```bash
./deploy_chaincode.sh query-committed startup common-channel
./deploy_chaincode.sh query-committed validator common-channel
./deploy_chaincode.sh query-committed investor common-channel
./deploy_chaincode.sh query-committed platform common-channel
```

**Check Commit Readiness (Before Commit):**
```bash
./deploy_chaincode.sh check-readiness startup common-channel
```

**Output:**
```
[INFO] Checking commit readiness for 'startup' on 'common-channel' (version: 1, sequence: 1)...
Chaincode definition for chaincode 'startup', version '1', sequence '1' on channel 'common-channel' approval status by org:
StartupOrgMSP: true
ValidatorOrgMSP: true
InvestorOrgMSP: true
PlatformOrgMSP: true
```

---

## 🔄 Upgrade Flow

### **STEP 1: Package New Versions** ⏱️ ~10 seconds

**Command:**
```bash
./deploy_chaincode.sh package
```

**What Happens:**
- Script detects existing versions (startup_1.tgz)
- Auto-increments to next version (startup_2.tgz)
- Packages updated chaincode from `./contracts/`

**Expected Output:**
```
📦 Packaging All Chaincodes
============================================================================
📦 Packaging chaincode: startup
   Current version: 1           ← Detected existing version
   New version: 2               ← Auto-incremented
   Package label: startup_2
   Package file: startup_2.tgz
✅ Successfully packaged startup as startup_2.tgz

[Repeat for validator_2.tgz, investor_2.tgz, platform_2.tgz]
```

---

### **STEP 2: Install New Versions** ⏱️ ~2-3 minutes

**Command:**
```bash
./deploy_chaincode.sh install-all
```

**What Happens:**
- Installs new package versions (startup_2.tgz, etc.) on all org peers
- Auto-exports new package IDs
- **Overwrites** previous package ID environment variables

**Expected Output:**
```
[INFO] Installing startup_2.tgz on current peer...
✅ Successfully installed startup on StartupOrg!
✅ Auto-exported: STARTUP_CC_PACKAGE_ID=startup_2:xyz789...  ← New package ID
```

---

### **STEP 3: Upgrade All Chaincodes** ⏱️ ~10-15 minutes

**Command:**
```bash
./deploy_chaincode.sh upgrade-all
```

**What Happens:**
1. Prompts for confirmation
2. For each channel (7 total):
   - For each chaincode on that channel:
     - Detects current sequence (e.g., Sequence 1)
     - Calculates next sequence (Sequence 2)
     - Approves upgrade from all orgs on the channel
     - Commits upgrade with new sequence and version
3. Processes channels in order:
   - common-channel (4 chaincodes × 4 orgs = 16 approvals + 4 commits)
   - startup-investor-channel (2 chaincodes × 2 orgs = 4 approvals + 2 commits)
   - startup-validator-channel (2 chaincodes × 2 orgs = 4 approvals + 2 commits)
   - startup-platform-channel (2 chaincodes × 2 orgs = 4 approvals + 2 commits)
   - investor-validator-channel (2 chaincodes × 2 orgs = 4 approvals + 2 commits)
   - investor-platform-channel (2 chaincodes × 2 orgs = 4 approvals + 2 commits)
   - validator-platform-channel (2 chaincodes × 2 orgs = 4 approvals + 2 commits)

**Expected Output:**
```
============================================================================
 Comprehensive Upgrade of All Chaincodes on All Channels
============================================================================
This will upgrade all chaincodes on all 7 channels.
[WARNING] Make sure all package IDs are exported:
  - STARTUP_CC_PACKAGE_ID
  - VALIDATOR_CC_PACKAGE_ID
  - INVESTOR_CC_PACKAGE_ID
  - PLATFORM_CC_PACKAGE_ID

Continue with upgrade? (y/n): y

============================================================================
 Upgrading common-channel
============================================================================
============================================================================
 Sync-Upgrade: startup on common-channel for startup
============================================================================
🔄 Switching to StartupOrg...
✅ Now operating as StartupOrg
[INFO] Smart approve for 'startup' on 'common-channel' by startup...
🔍 Checking if startup is already committed on common-channel...
[INFO] startup is committed (sequence: 1). Proceeding with approval for startup...
[INFO] Approving chaincode 'startup' on channel 'common-channel' (version: 2, sequence: 2)...
✅ Approved upgrade for 'startup' on 'common-channel' (sequence: 2)

============================================================================
 Sync-Upgrade: startup on common-channel for validator
============================================================================
🔄 Switching to ValidatorOrg...
[INFO] startup is committed (sequence: 1). Proceeding with approval for validator...
✅ Approved upgrade for 'startup' on 'common-channel' (sequence: 2)

[Repeat for investor and platform orgs]

🔄 Switching to StartupOrg...
[INFO] Committing chaincode 'startup' on channel 'common-channel' (version: 2, sequence: 2)...
✅ Successfully committed 'startup' on 'common-channel' (version: 2, sequence: 2)

[Process continues for all chaincodes and channels]

============================================================================
 Upgrading startup-investor-channel
============================================================================
[Similar process for startup-investor-channel]

============================================================================
 Upgrading startup-validator-channel
============================================================================
[Similar process for startup-validator-channel]

[... continues for all 7 channels ...]

✅ All chaincodes upgraded on all channels!
```

---

### **Alternative: Upgrade Specific Chaincode** ⏱️ ~1-2 minutes

**Command:**
```bash
./deploy_chaincode.sh upgrade startup common-channel
```

**What Happens:**
- Approves upgrade from all orgs on common-channel (4 orgs)
- Commits upgrade from StartupOrg

**Expected Output:**
```
============================================================================
 Upgrading startup on common-channel
============================================================================
[INFO] Smart approve for 'startup' on 'common-channel' by startup...
[INFO] startup is committed (sequence: 1). Proceeding with approval...
✅ Approved upgrade for 'startup' on 'common-channel' (sequence: 2)

[Repeat for validator, investor, platform orgs]

[INFO] Committing chaincode 'startup' on channel 'common-channel' (version: 2, sequence: 2)...
✅ Successfully committed 'startup' on 'common-channel' (version: 2, sequence: 2)
✅ Upgrade complete for startup on common-channel
```

---

### **Alternative: Manual Approval (When Another Org Upgraded)** ⏱️ ~30 seconds

**Scenario:** StartupOrg upgraded `startup` chaincode to v2. ValidatorOrg (sharing common-channel) needs to approve.

**Command:**
```bash
# On ValidatorOrg
./deploy_chaincode.sh approve-chaincode validator startup common-channel
```

**What Happens:**
- Switches to ValidatorOrg context
- Detects current committed sequence (Sequence 1)
- Calculates next sequence (Sequence 2)
- Approves upgrade with Sequence 2
- Does NOT commit (only approval)

**Expected Output:**
```
============================================================================
 Approving startup chaincode on common-channel for validator
============================================================================
🔄 Switching to ValidatorOrg...
✅ Now operating as ValidatorOrg

[INFO] Smart approve for 'startup' on 'common-channel' by validator...
[INFO] startup is committed (sequence: 1). Proceeding with approval for validator...
[INFO] Approving chaincode 'startup' on channel 'common-channel' (version: 2, sequence: 2)...
✅ Approved upgrade for 'startup' on 'common-channel' (sequence: 2)

✅ Approved startup on common-channel for validator
```

**Use Cases:**
1. Another org upgraded and you need to catch up
2. You want granular control over approvals
3. Testing upgrade approval before full deployment

---

## 📊 Channel & Chaincode Matrix

### **Chaincode Deployment Matrix**

| Channel                      | StartupOrg | ValidatorOrg | InvestorOrg | PlatformOrg |
|------------------------------|------------|--------------|-------------|-------------|
| **common-channel**           | ✅ All 4    | ✅ All 4      | ✅ All 4     | ✅ All 4     |
| **startup-investor-channel** | ✅ startup  | ❌           | ✅ investor  | ❌          |
| **startup-validator-channel**| ✅ startup  | ✅ validator  | ❌          | ❌          |
| **startup-platform-channel** | ✅ startup  | ❌           | ❌          | ✅ platform  |
| **investor-validator-channel**| ❌         | ✅ validator  | ✅ investor  | ❌          |
| **investor-platform-channel**| ❌         | ❌           | ✅ investor  | ✅ platform  |
| **validator-platform-channel**| ❌         | ✅ validator  | ❌          | ✅ platform  |

### **Chaincode Ownership (Who Commits)**

| Chaincode | Owner Org    | Commits On Channels |
|-----------|--------------|---------------------|
| startup   | StartupOrg   | startup-validator, startup-investor, startup-platform, common |
| validator | ValidatorOrg | startup-validator, investor-validator, validator-platform, common |
| investor  | InvestorOrg  | startup-investor, investor-validator, investor-platform, common |
| platform  | PlatformOrg  | startup-platform, investor-platform, validator-platform, common |

---

## 🔍 Troubleshooting

### **Issue: Package ID not found**
```
❌ STARTUP_CC_PACKAGE_ID not set. Please export it first.
```

**Solution:**
```bash
# Query installed chaincodes
peer lifecycle chaincode queryinstalled

# Export the package ID manually
export STARTUP_CC_PACKAGE_ID="startup_1:abc123..."
```

### **Issue: Approval fails**
```
❌ Failed to approve 'startup' on 'common-channel'
```

**Solution:**
1. Check if chaincode is installed:
   ```bash
   peer lifecycle chaincode queryinstalled
   ```
2. Verify package ID is correct
3. Check if you're on the correct org:
   ```bash
   ./deploy_chaincode.sh switch startup
   ```

### **Issue: Commit fails**
```
❌ Failed to commit 'startup' on 'common-channel'
Error: proposal failed with status: 500 - failed to invoke backing implementation...
```

**Solution:**
1. Check if all orgs have approved:
   ```bash
   ./deploy_chaincode.sh check-readiness startup common-channel
   ```
2. Ensure majority approvals (e.g., 3/4 orgs for common-channel)

### **Issue: Already installed error**
```
Error: chaincode already successfully installed
```

**This is normal!** The script detects this and auto-exports the existing package ID.

---

## ⚡ Quick Reference

### **Initial Deployment (3-Step Process)**
```bash
# Step 1: Package
./deploy_chaincode.sh package

# Step 2: Install
./deploy_chaincode.sh install-all

# Step 3: Deploy
./deploy_chaincode.sh deploy startup
./deploy_chaincode.sh deploy validator
./deploy_chaincode.sh deploy investor
./deploy_chaincode.sh deploy platform
```

### **Upgrade (3-Step Process)**
```bash
# Step 1: Package new versions
./deploy_chaincode.sh package

# Step 2: Install new versions
./deploy_chaincode.sh install-all

# Step 3: Upgrade all
./deploy_chaincode.sh upgrade-all
```

### **Verify Deployment**
```bash
# Query committed
./deploy_chaincode.sh query-committed startup common-channel

# Check readiness
./deploy_chaincode.sh check-readiness startup common-channel
```

### **Switch Organization Context**
```bash
./deploy_chaincode.sh switch startup
./deploy_chaincode.sh switch validator
./deploy_chaincode.sh switch investor
./deploy_chaincode.sh switch platform
```

---

## ⏱️ Typical Execution Times

| Operation      | Time Estimate | Details |
|----------------|---------------|---------|
| **Package**    | ~10 seconds   | All 4 chaincodes |
| **Install-all**| ~2-3 minutes  | 4 orgs × 4 chaincodes with prompts |
| **Deploy all** | ~5-7 minutes  | 4 orgs × multiple channels |
| **Upgrade-all**| ~10-15 minutes| 7 channels × 4 chaincodes × multiple orgs |

---

## 🎯 Success Indicators

### **After Packaging:**
```
✅ 4 .tgz files created in current directory
✅ Version numbers auto-incremented
```

### **After Installation:**
```
✅ 4 environment variables exported:
   STARTUP_CC_PACKAGE_ID
   VALIDATOR_CC_PACKAGE_ID
   INVESTOR_CC_PACKAGE_ID
   PLATFORM_CC_PACKAGE_ID
```

### **After Deployment:**
```
✅ All chaincodes show Version 1, Sequence 1
✅ query-committed shows all orgs approved
✅ No error messages in deployment output
```

### **After Upgrade:**
```
✅ All chaincodes show Version 2, Sequence 2 (or next sequence)
✅ No "failed to approve" errors
✅ All commits successful
```

---

## 📞 Support

For issues or questions:
1. Check the [Troubleshooting](#troubleshooting) section
2. Review the script header comments in `deploy_chaincode.sh`
3. Use `./deploy_chaincode.sh help` for detailed usage

---

**Last Updated:** December 12, 2025

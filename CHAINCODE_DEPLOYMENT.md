# Chaincode Deployment Commands

## Prerequisites
Set the package ID after installing chaincode:
```bash
export CC_PACKAGE_ID=<your_package_id>
```

**Note:** Port is 9090 as per MICROFAB.txt configuration.

---

## Channel Membership Overview

| Channel | Members |
|---------|---------|
| startup-validator-channel | StartupOrg, ValidatorOrg |
| startup-platform-channel | StartupOrg, PlatformOrg |
| startup-investor-channel | StartupOrg, InvestorOrg |
| investor-platform-channel | InvestorOrg, PlatformOrg |
| investor-validator-channel | InvestorOrg, ValidatorOrg |
| validator-platform-channel | ValidatorOrg, PlatformOrg |
| common-channel | StartupOrg, ValidatorOrg, InvestorOrg, PlatformOrg |

---

## 1. STARTUP-VALIDATOR-CHANNEL

### Chaincodes: startup and validator in StartupOrg

**Approve (startup) in StartupOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel --name startup --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (startup) in StartupOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel --name startup --version 1 --sequence 1
```

**Approve (validator) in StartupOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel --name validator --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

### Chaincodes: validator and startup in ValidatorOrg

**Approve (validator) in ValidatorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel --name validator --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (validator) in ValidatorOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel --name validator --version 1 --sequence 1
```

**Approve (startup) in ValidatorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel --name startup --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

## 2. STARTUP-PLATFORM-CHANNEL

### Chaincodes: startup and platform in StartupOrg

**Approve (startup) in StartupOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-platform-channel --name startup --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (startup) in StartupOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-platform-channel --name startup --version 1 --sequence 1
```

**Approve (platform) in StartupOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-platform-channel --name platform --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

### Chaincodes: platform and startup in PlatformOrg

**Approve (platform) in PlatformOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-platform-channel --name platform --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (platform) in PlatformOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-platform-channel --name platform --version 1 --sequence 1
```

**Approve (startup) in PlatformOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-platform-channel --name startup --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

## 3. STARTUP-INVESTOR-CHANNEL

### Chaincodes: startup and investor in StartupOrg

**Approve (startup) in StartupOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel --name startup --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (startup) in StartupOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel --name startup --version 1 --sequence 1
```

**Approve (investor) in StartupOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel --name investor --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

### Chaincodes: investor and startup in InvestorOrg

**Approve (investor) in InvestorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel --name investor --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (investor) in InvestorOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel --name investor --version 1 --sequence 1
```

**Approve (startup) in InvestorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel --name startup --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

## 4. INVESTOR-PLATFORM-CHANNEL

### Chaincodes: investor and platform in InvestorOrg

**Approve (investor) in InvestorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel --name investor --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (investor) in InvestorOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel --name investor --version 1 --sequence 1
```

**Approve (platform) in InvestorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel --name platform --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

### Chaincodes: platform and investor in PlatformOrg

**Approve (platform) in PlatformOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel --name platform --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (platform) in PlatformOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel --name platform --version 1 --sequence 1
```

**Approve (investor) in PlatformOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel --name investor --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

## 5. INVESTOR-VALIDATOR-CHANNEL

### Chaincodes: investor and validator in InvestorOrg

**Approve (investor) in InvestorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel --name investor --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (investor) in InvestorOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel --name investor --version 1 --sequence 1
```

**Approve (validator) in InvestorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel --name validator --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

### Chaincodes: validator and investor in ValidatorOrg

**Approve (validator) in ValidatorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel --name validator --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (validator) in ValidatorOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel --name validator --version 1 --sequence 1
```

**Approve (investor) in ValidatorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel --name investor --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

## 6. VALIDATOR-PLATFORM-CHANNEL

### Chaincodes: validator and platform in ValidatorOrg

**Approve (validator) in ValidatorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel --name validator --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (validator) in ValidatorOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel --name validator --version 1 --sequence 1
```

**Approve (platform) in ValidatorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel --name platform --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

### Chaincodes: platform and validator in PlatformOrg

**Approve (platform) in PlatformOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel --name platform --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (platform) in PlatformOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel --name platform --version 1 --sequence 1
```

**Approve (validator) in PlatformOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel --name validator --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

## 7. COMMON-CHANNEL (All Organizations)

### Chaincodes: startup, validator, investor, platform in StartupOrg

**Approve (startup) in StartupOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name startup --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (startup) in StartupOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name startup --version 1 --sequence 1
```

**Approve (validator) in StartupOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name validator --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Approve (investor) in StartupOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name investor --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Approve (platform) in StartupOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name platform --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

### Chaincodes: validator, startup, investor, platform in ValidatorOrg

**Approve (validator) in ValidatorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name validator --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (validator) in ValidatorOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name validator --version 1 --sequence 1
```

**Approve (startup) in ValidatorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name startup --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Approve (investor) in ValidatorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name investor --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Approve (platform) in ValidatorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name platform --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

### Chaincodes: investor, startup, validator, platform in InvestorOrg

**Approve (investor) in InvestorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name investor --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (investor) in InvestorOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name investor --version 1 --sequence 1
```

**Approve (startup) in InvestorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name startup --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Approve (validator) in InvestorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name validator --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Approve (platform) in InvestorOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name platform --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

### Chaincodes: platform, startup, validator, investor in PlatformOrg

**Approve (platform) in PlatformOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name platform --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Commit (platform) in PlatformOrg:**
```bash
peer lifecycle chaincode commit -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name platform --version 1 --sequence 1
```

**Approve (startup) in PlatformOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name startup --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Approve (validator) in PlatformOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name validator --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

**Approve (investor) in PlatformOrg:**
```bash
peer lifecycle chaincode approveformyorg -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel --name investor --version 1 --sequence 1 --waitForEvent --package-id ${CC_PACKAGE_ID}
```

---

## Quick Reference: Chaincodes per Channel

| Channel | Chaincodes | Member Orgs |
|---------|------------|-------------|
| startup-validator-channel | startup, validator | StartupOrg, ValidatorOrg |
| startup-platform-channel | startup, platform | StartupOrg, PlatformOrg |
| startup-investor-channel | startup, investor | StartupOrg, InvestorOrg |
| investor-platform-channel | investor, platform | InvestorOrg, PlatformOrg |
| investor-validator-channel | investor, validator | InvestorOrg, ValidatorOrg |
| validator-platform-channel | validator, platform | ValidatorOrg, PlatformOrg |
| common-channel | startup, validator, investor, platform | All Orgs |

---

## Check Commit Readiness

Before committing, verify all orgs have approved:
```bash
peer lifecycle chaincode checkcommitreadiness --channelID <channel-name> --name <chaincode-name> --version 1 --sequence 1
```

## Query Committed Chaincodes

After commit, verify chaincode is committed:
```bash
peer lifecycle chaincode querycommitted --channelID <channel-name> --name <chaincode-name>
```

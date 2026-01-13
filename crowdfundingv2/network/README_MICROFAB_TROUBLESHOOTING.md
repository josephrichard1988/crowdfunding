# Microfab Endorsement & Discovery Troubleshooting

This document details the solution for `ENDORSEMENT_POLICY_FAILURE` and "No peers found" errors encountered when working with Microfab and Private Data Collections (PDC).

## The Issue

When submitting transactions that write to Private Data Collections (e.g., `CreateStartup`), you might encounter:

1. **ENDORSEMENT_POLICY_FAILURE**: The transaction is endorsed by a peer that is NOT a member of the collection (e.g., `ValidatorOrg` endorsing a write to `StartupPrivateData`).
2. **"No valid responses from any peers"**: The Fabric SDK's Service Discovery (`channel.getEndorsers()`) returns an empty list, causing the transaction to fail before submission.

**Root Cause:** Microfab's Service Discovery implementation can be flaky or return incomplete network views, especially when running in a single-container local environment with multiple organizations. It often defaults to round-robin routing which violates per-collection endorsement policies.

## The Solution (Implemented in `fabricConnection.js`)

We have implemented a **Robust Peer Selection Strategy** that bypasses the default SDK discovery when it fails.

### Logic Flow

1. **Explicit Targeting:** The code first attempts to identify the correct identifying MSP ID for the operation (e.g., `StartupOrgMSP` for `CreateStartup`).
2. **Discovery Attempt:** It calls `channel.getEndorsers()` to find peers belonging to that MSP.
3. **Fallback Mechanism (CRITICAL):**
    * If discovery returns **0 peers** (the "flaky" scenario), the code **probes the Connection Profile** directly.
    * It reads `config.gatewaysDir` and the specific gateway JSON file.
    * It extracts the peer names explicitly defined in the profile (e.g., `startuporgpeer-api.127-0-0-1.nip.io:9090`).
    * It manually looks up these peers using `channel.getEndorser(peerName)`.
4. **Forced Endorsement:** It uses `transaction.setEndorsingPeers([specificPeer])` to force the SDK to send the proposal ONLY to the correct, member-org peer.

### Code Reference

See `network/src/fabricConnection.js` -> `submitTransaction` method:

```javascript
// If no endorsers found via discovery, try to get from connection profile explicit peers
if (orgEndorsers.length === 0) {
    logger.warn(`⚠️ Discovery returned 0 endorsers for ${mspId}. Probing connection profile...`);
    // ... logic to read connection profile and add peers ...
}

// ... 

if (orgEndorsers.length > 0) {
    // Set explicit endorsing peers
    const transaction = contract.createTransaction(fcn);
    transaction.setEndorsingPeers(orgEndorsers);
    // ...
}
```

## Maintenance

If you add new organizations or change the network topology:

1. Ensure `config/index.js` defines the correct `gatewayFile`.
2. Ensure the Gateway JSON files in `_gateways/` contain the correct peer names.
3. The logic in `fabricConnection.js` will automatically pick these up.

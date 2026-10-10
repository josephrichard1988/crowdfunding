# Crowdfunding Dual-Channel v2

Combined chaincode packages on the Microfab dual-channel topology.

| | `crowdfunding_dualchannel/` | `crowdfunding_dualchannel_v2/` (this folder) |
|---|---|---|
| Chaincodes | 4 org CCs: `startup`, `validator`, `investor`, `platform` | 2 combined CCs: `governance`, `investment` |
| Governance channel | 3 CCs (no investor) | 1 CC: `governance` |
| Investment channel | 4 CCs | 1 CC: `investment` |
| Cross-channel | `InvokeChaincode` to org CC names | `DataBridgeContract` + `InvokeChaincode` to `governance` / `investment` |
| CLI invoke | `-n startup` + bare function | `-n governance` + `StartupContract:Fn` |

## Topology (MICROFAB.txt, port 7070)

- **governance-validation-channel** — StartupOrg, ValidatorOrg, PlatformOrg  
  Chaincode: `governance` (Startup + Validator + Platform + DataBridge)  
  Policy: `AND('ValidatorOrgMSP.peer','PlatformOrgMSP.peer')`

- **investment-execution-channel** — all 4 orgs  
  Chaincode: `investment` (all 4 + Token + DataBridge)  
  Policy: `OR(...all 4 peers...)`

## Deploy

```bash
cd crowdfunding_dualchannel_v2
# Microfab + weft wallets/msp, then:
source ./deploy_chaincode.sh package
source ./deploy_chaincode.sh install-all
source ./deploy_chaincode.sh deploy all
```

## E2E

See [DUAL_CHANNEL_E2E_TESTING.md](./DUAL_CHANNEL_E2E_TESTING.md).

Example:

```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:7070 \
  --channelID governance-validation-channel -n governance \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:7070 \
  -c '{"function":"StartupContract:CreateCampaign","Args":[...]}'
```

## Layout

```
crowdfunding_dualchannel_v2/
  MICROFAB.txt
  collections_config_governance.json
  collections_config_investment.json
  deploy_chaincode.sh
  DUAL_CHANNEL_E2E_TESTING.md
  contracts/
    governance/   # combined package for governance-validation-channel
    investment/   # combined package for investment-execution-channel
```

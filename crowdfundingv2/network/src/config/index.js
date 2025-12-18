const path = require('path');

module.exports = {
    port: process.env.NETWORK_PORT || 4000,

    // Directories
    gatewaysDir: path.resolve(process.cwd(), '..', '_gateways'),
    walletsDir: path.resolve(process.cwd(), '..', '_wallets'),
    mspDir: path.resolve(process.cwd(), '..', '_msp'),

    // Fabric
    channelName: 'crowdfunding-channel',
    chaincodeName: 'crowdfunding',

    // Organizations
    orgs: {
        startup: {
            name: 'StartupOrg',
            mspId: 'StartupOrgMSP',
            caUrl: 'http://startuporgca-api.127-0-0-1.nip.io:9090',
            gatewayFile: 'startuporggateway.json',
            adminUser: 'startuporgadmin',
            adminPassword: 'startuporgadminpw',
        },
        investor: {
            name: 'InvestorOrg',
            mspId: 'InvestorOrgMSP',
            caUrl: 'http://investororgca-api.127-0-0-1.nip.io:9090',
            gatewayFile: 'investororggateway.json',
            adminUser: 'investororgadmin',
            adminPassword: 'investororgadminpw',
        },
        platform: {
            name: 'PlatformOrg',
            mspId: 'PlatformOrgMSP',
            caUrl: 'http://platformorgca-api.127-0-0-1.nip.io:9090',
            gatewayFile: 'platformorggateway.json',
            adminUser: 'platformorgadmin',
            adminPassword: 'platformorgadminpw',
        },
        validator: {
            name: 'ValidatorOrg',
            mspId: 'ValidatorOrgMSP',
            caUrl: 'http://validatororgca-api.127-0-0-1.nip.io:9090',
            gatewayFile: 'validatororggateway.json',
            adminUser: 'validatororgadmin',
            adminPassword: 'validatororgadminpw',
        },
    },
};

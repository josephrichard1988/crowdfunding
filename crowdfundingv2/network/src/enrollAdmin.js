const FabricCAServices = require('fabric-ca-client');
const { Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const config = require('./config');
const logger = require('./utils/logger');

/**
 * Enroll admin users for all organizations
 */
async function enrollAdmins() {
    for (const [orgKey, org] of Object.entries(config.orgs)) {
        try {
            logger.info(`Enrolling admin for ${org.name}...`);

            // Create wallet
            const walletPath = path.join(config.walletsDir, org.name);
            const wallet = await Wallets.newFileSystemWallet(walletPath);

            // Check if already enrolled
            const identity = await wallet.get(org.adminUser);
            if (identity) {
                logger.info(`Admin ${org.adminUser} already enrolled for ${org.name}`);
                continue;
            }

            // Connect to CA
            const ca = new FabricCAServices(org.caUrl);

            // Enroll admin
            const enrollment = await ca.enroll({
                enrollmentID: org.adminUser,
                enrollmentSecret: org.adminPassword,
            });

            // Create identity
            const x509Identity = {
                credentials: {
                    certificate: enrollment.certificate,
                    privateKey: enrollment.key.toBytes(),
                },
                mspId: org.mspId,
                type: 'X.509',
            };

            // Store in wallet
            await wallet.put(org.adminUser, x509Identity);
            logger.info(`✅ Successfully enrolled ${org.adminUser} for ${org.name}`);

        } catch (error) {
            logger.error(`❌ Failed to enroll admin for ${org.name}: ${error.message}`);
        }
    }
}

// Run if called directly
if (require.main === module) {
    enrollAdmins()
        .then(() => logger.info('Enrollment complete'))
        .catch(err => logger.error(`Enrollment failed: ${err.message}`));
}

module.exports = { enrollAdmins };

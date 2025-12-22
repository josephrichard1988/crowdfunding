import express from 'express';
import User from '../models/User.js';
import { generateToken, authMiddleware } from '../middleware/auth.js';

const router = express.Router();

// ============================================================================
// SIGNUP
// ============================================================================
router.post('/signup', async (req, res) => {
    try {
        const { email, password, name, role } = req.body;

        // Validate required fields
        if (!email || !password || !name || !role) {
            return res.status(400).json({
                error: 'Missing required fields: email, password, name, role'
            });
        }

        // Validate role
        const validRoles = ['STARTUP', 'INVESTOR', 'VALIDATOR', 'PLATFORM'];
        if (!validRoles.includes(role.toUpperCase())) {
            return res.status(400).json({
                error: `Invalid role. Must be one of: ${validRoles.join(', ')}`
            });
        }

        // Check if user exists
        const existingUser = await User.findOne({ email: email.toLowerCase() });
        if (existingUser) {
            return res.status(400).json({ error: 'Email already registered' });
        }

        // Create user
        const user = new User({
            email: email.toLowerCase(),
            password,
            name,
            role: role.toUpperCase()
        });

        await user.save();

        // Generate token
        const token = generateToken({
            userId: user._id,
            orgUserId: user.orgUserId,
            email: user.email,
            role: user.role,
            name: user.name
        });

        res.status(201).json({
            message: 'User created successfully',
            token,
            user: user.toJSON()
        });

    } catch (error) {
        console.error('Signup error:', error);
        res.status(500).json({ error: 'Failed to create user' });
    }
});

// ============================================================================
// LOGIN
// ============================================================================
router.post('/login', async (req, res) => {
    try {
        const { email, password } = req.body;

        if (!email || !password) {
            return res.status(400).json({ error: 'Email and password required' });
        }

        // Find user
        const user = await User.findOne({ email: email.toLowerCase() });
        if (!user) {
            return res.status(401).json({ error: 'Invalid credentials' });
        }

        // Check password
        if (!user.comparePassword(password)) {
            return res.status(401).json({ error: 'Invalid credentials' });
        }

        // Check if active
        if (!user.isActive) {
            return res.status(403).json({ error: 'Account is disabled' });
        }

        // Generate token
        const token = generateToken({
            userId: user._id,
            orgUserId: user.orgUserId,
            email: user.email,
            role: user.role,
            name: user.name
        });

        res.json({
            message: 'Login successful',
            token,
            user: user.toJSON()
        });

    } catch (error) {
        console.error('Login error:', error);
        res.status(500).json({ error: 'Login failed' });
    }
});

// ============================================================================
// GET CURRENT USER
// ============================================================================
router.get('/me', authMiddleware, async (req, res) => {
    try {
        const user = await User.findById(req.user.userId);

        if (!user) {
            return res.status(404).json({ error: 'User not found' });
        }

        res.json({ user: user.toJSON() });

    } catch (error) {
        console.error('Get user error:', error);
        res.status(500).json({ error: 'Failed to get user' });
    }
});

// ============================================================================
// UPDATE WALLET (sync from chaincode)
// ============================================================================
router.put('/wallet', authMiddleware, async (req, res) => {
    try {
        const { cftBalance, cfrtBalance, frozenCft } = req.body;

        const user = await User.findById(req.user.userId);
        if (!user) {
            return res.status(404).json({ error: 'User not found' });
        }

        user.wallet = {
            cftBalance: cftBalance ?? user.wallet.cftBalance,
            cfrtBalance: cfrtBalance ?? user.wallet.cfrtBalance,
            frozenCft: frozenCft ?? user.wallet.frozenCft,
            lastSynced: new Date()
        };

        await user.save();

        res.json({ wallet: user.wallet });

    } catch (error) {
        console.error('Update wallet error:', error);
        res.status(500).json({ error: 'Failed to update wallet' });
    }
});

// ============================================================================
// GET FEE SCHEDULE
// ============================================================================
router.get('/fees', (req, res) => {
    res.json({
        exchangeRate: {
            INR: parseFloat(process.env.INR_TO_CFT_RATE) || 2.5,
            USD: parseFloat(process.env.USD_TO_CFT_RATE) || 83.0
        },
        fees: {
            registrationFee: parseInt(process.env.REGISTRATION_FEE_CFT) || 250,
            campaignCreationFee: parseInt(process.env.CAMPAIGN_CREATION_FEE_CFT) || 1250,
            campaignPublishingFee: parseInt(process.env.CAMPAIGN_PUBLISHING_FEE_CFT) || 2500,
            validationFee: parseInt(process.env.VALIDATION_FEE_CFT) || 500,
            disputeFee: parseInt(process.env.DISPUTE_FEE_CFT) || 750,
            investmentFeePercent: parseFloat(process.env.INVESTMENT_FEE_PERCENT) || 5,
            withdrawalFeePercent: parseFloat(process.env.WITHDRAWAL_FEE_PERCENT) || 1
        }
    });
});

export default router;

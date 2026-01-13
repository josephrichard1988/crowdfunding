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

// ============================================================================
// STARTUP MANAGEMENT (for STARTUP role users)
// ============================================================================

// Create a new startup for the user
router.post('/startups', authMiddleware, async (req, res) => {
    try {
        const { name, description, startupId: providedStartupId, displayId: providedDisplayId } = req.body;

        if (!name) {
            return res.status(400).json({ error: 'Startup name is required' });
        }

        // Try to find user by _id first, then fallback to orgUserId (more resilient)
        let user = await User.findById(req.user.userId);
        if (!user && req.user.orgUserId) {
            console.log('[DEBUG /startups] Fallback: Looking up by orgUserId:', req.user.orgUserId);
            user = await User.findOne({ orgUserId: req.user.orgUserId });
        }
        console.log('[DEBUG /startups] User found:', user ? `${user.email} (role: ${user.role})` : 'NOT FOUND');

        if (!user) {
            return res.status(403).json({ error: 'Only STARTUP users can create startups', debug: 'User not found in DB' });
        }
        if (user.role !== 'STARTUP') {
            return res.status(403).json({ error: 'Only STARTUP users can create startups', debug: `User role is ${user.role}` });
        }

        let startupId = providedStartupId;
        let displayId = providedDisplayId;

        // If not provided, generate uniquely (fallback)
        if (!startupId) {
            const seq = (user.startups?.length || 0) + 1;
            displayId = `S-${String(seq).padStart(3, '0')}`;
            startupId = `STU_${user.orgUserId}_${String(seq).padStart(3, '0')}`;
        }

        // Add startup to user's startups array
        user.startups.push({
            startupId,
            displayId,
            name,
            description: description || '',
            campaigns: []
        });

        await user.save();

        res.status(201).json({
            success: true,
            data: {
                startupId,
                displayId,
                name,
                description
            }
        });
    } catch (error) {
        console.error('Create startup error:', error);
        res.status(500).json({ error: error.message });
    }
});

// Get user's startups
router.get('/startups', authMiddleware, async (req, res) => {
    try {
        const user = await User.findById(req.user.userId);
        if (!user) {
            return res.status(404).json({ error: 'User not found' });
        }

        res.json({
            success: true,
            data: user.startups || []
        });
    } catch (error) {
        console.error('Get startups error:', error);
        res.status(500).json({ error: error.message });
    }
});

// ============================================================================
// CAMPAIGN SYNC (sync Fabric campaign to MongoDB)
// ============================================================================

// Sync a campaign to MongoDB after creation in Fabric
router.post('/sync/campaign', authMiddleware, async (req, res) => {
    try {
        const { startupId, campaignId, displayId, projectName, status } = req.body;

        const user = await User.findById(req.user.userId);
        if (!user) {
            return res.status(404).json({ error: 'User not found' });
        }

        // Find the startup and add campaign
        const startup = user.startups.find(s => s.startupId === startupId);
        if (!startup) {
            return res.status(404).json({ error: 'Startup not found' });
        }

        // Check if campaign already exists
        const existingCampaign = startup.campaigns.find(c => c.campaignId === campaignId);
        if (existingCampaign) {
            // Update existing
            existingCampaign.status = status || existingCampaign.status;
            existingCampaign.projectName = projectName || existingCampaign.projectName;
        } else {
            // Add new
            startup.campaigns.push({
                campaignId,
                displayId,
                projectName,
                status: status || 'DRAFT',
                validationStatus: 'NOT_SUBMITTED'
            });
        }

        await user.save();

        res.json({ success: true, message: 'Campaign synced to MongoDB' });
    } catch (error) {
        console.error('Sync campaign error:', error);
        res.status(500).json({ error: error.message });
    }
});

// Update campaign status in MongoDB (for validation submission, approval, publishing)
router.post('/sync/campaign-status', authMiddleware, async (req, res) => {
    try {
        const { startupId, campaignId, status, validationStatus } = req.body;

        // Try to find user by _id first, then fallback to orgUserId
        let user = await User.findById(req.user.userId);
        if (!user && req.user.orgUserId) {
            user = await User.findOne({ orgUserId: req.user.orgUserId });
        }

        if (!user) {
            return res.status(404).json({ error: 'User not found' });
        }

        // Find the startup
        const startup = user.startups.find(s => s.startupId === startupId);
        if (!startup) {
            return res.status(404).json({ error: 'Startup not found' });
        }

        // Find and update the campaign
        const campaign = startup.campaigns.find(c => c.campaignId === campaignId);
        if (!campaign) {
            return res.status(404).json({ error: 'Campaign not found' });
        }

        // Update fields if provided
        if (status) campaign.status = status;
        if (validationStatus) campaign.validationStatus = validationStatus;

        await user.save();

        console.log(`[SYNC] Campaign ${campaignId} updated: status=${status}, validationStatus=${validationStatus}`);
        res.json({ success: true, message: 'Campaign status updated' });
    } catch (error) {
        console.error('Sync campaign status error:', error);
        res.status(500).json({ error: error.message });
    }
});

// Update startup status (e.g. for deletion)
router.post('/sync/startup-status', authMiddleware, async (req, res) => {
    try {
        const { startupId, status } = req.body;

        const user = await User.findById(req.user.userId);
        if (!user) {
            return res.status(404).json({ error: 'User not found' });
        }

        // Find the startup
        const startup = user.startups.find(s => s.startupId === startupId);
        if (!startup) {
            return res.status(404).json({ error: 'Startup not found' });
        }

        // Update status
        if (status) startup.status = status;

        await user.save();
        console.log(`[SYNC] Startup ${startupId} updated: status=${status}`);

        res.json({ success: true, message: 'Startup status updated' });
    } catch (error) {
        console.error('Sync startup status error:', error);
        res.status(500).json({ error: error.message });
    }
});

// ============================================================================
// ALLOCATION MANAGEMENT (for VALIDATOR/PLATFORM assignment)
// ============================================================================

// Get the validator/platform with least queue (for auto-allocation)
router.get('/allocation/next', authMiddleware, async (req, res) => {
    try {
        const { role } = req.query;  // 'VALIDATOR' or 'PLATFORM'

        if (!['VALIDATOR', 'PLATFORM'].includes(role)) {
            return res.status(400).json({ error: 'Role must be VALIDATOR or PLATFORM' });
        }

        // Find user with least assignedQueue
        const users = await User.find({ role, isActive: true })
            .select('orgUserId name assignedQueue')
            .lean();

        if (users.length === 0) {
            return res.status(404).json({ error: `No active ${role} users found` });
        }

        // Sort by queue length and get the one with least
        users.sort((a, b) => (a.assignedQueue?.length || 0) - (b.assignedQueue?.length || 0));
        const selected = users[0];

        res.json({
            success: true,
            data: {
                orgUserId: selected.orgUserId,
                name: selected.name,
                queueLength: selected.assignedQueue?.length || 0
            }
        });
    } catch (error) {
        console.error('Get next allocation error:', error);
        res.status(500).json({ error: error.message });
    }
});

// Assign campaign to validator/platform queue
router.post('/allocation/assign', authMiddleware, async (req, res) => {
    try {
        const { assigneeOrgUserId, campaignId, startupId, projectName, type } = req.body;

        if (!assigneeOrgUserId || !campaignId || !type) {
            return res.status(400).json({ error: 'assigneeOrgUserId, campaignId, and type are required' });
        }

        // Add to assignee's queue
        const assignee = await User.findOne({ orgUserId: assigneeOrgUserId });
        if (!assignee) {
            return res.status(404).json({ error: 'Assignee not found' });
        }

        assignee.assignedQueue.push({
            campaignId,
            startupId,
            projectName,
            type,  // 'VALIDATION' or 'PUBLISH'
            assignedAt: new Date()
        });

        await assignee.save();

        res.json({
            success: true,
            message: `Campaign assigned to ${assignee.name}`,
            data: { assigneeOrgUserId, campaignId, type }
        });
    } catch (error) {
        console.error('Assign campaign error:', error);
        res.status(500).json({ error: error.message });
    }
});

// Complete a task (remove from queue, add to completed)
router.post('/allocation/complete', authMiddleware, async (req, res) => {
    try {
        const { campaignId, type, result } = req.body;

        const user = await User.findById(req.user.userId);
        if (!user) {
            return res.status(404).json({ error: 'User not found' });
        }

        // Find and remove from queue
        const taskIndex = user.assignedQueue.findIndex(
            t => t.campaignId === campaignId && t.type === type
        );

        if (taskIndex === -1) {
            return res.status(404).json({ error: 'Task not found in queue' });
        }

        const task = user.assignedQueue[taskIndex];
        user.assignedQueue.splice(taskIndex, 1);

        // Add to completed
        user.completedTasks.push({
            campaignId,
            projectName: task.projectName,
            type,
            result,  // 'APPROVED', 'REJECTED', 'PUBLISHED'
            completedAt: new Date()
        });

        await user.save();

        res.json({ success: true, message: 'Task completed' });
    } catch (error) {
        console.error('Complete task error:', error);
        res.status(500).json({ error: error.message });
    }
});

// Get user's assigned queue
router.get('/queue', authMiddleware, async (req, res) => {
    try {
        const user = await User.findById(req.user.userId);
        if (!user) {
            return res.status(404).json({ error: 'User not found' });
        }

        res.json({
            success: true,
            data: {
                assignedQueue: user.assignedQueue || [],
                completedTasks: user.completedTasks || []
            }
        });
    } catch (error) {
        console.error('Get queue error:', error);
        res.status(500).json({ error: error.message });
    }
});

// ============================================================================
// CFT SUPPLY MANAGEMENT (Platform Admin Only)
// ============================================================================

// In-memory CFT supply (in production, use a proper settings collection)
let globalCftSupply = 0;

// Get available CFT supply
router.get('/cft-supply', authMiddleware, async (req, res) => {
    try {
        res.json({
            success: true,
            availableCft: globalCftSupply
        });
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

// Set CFT supply (Platform Admin only)
router.post('/cft-supply/set', authMiddleware, async (req, res) => {
    try {
        // Only Platform admins can set supply
        if (req.user.role !== 'PLATFORM') {
            return res.status(403).json({ error: 'Only Platform admins can set CFT supply' });
        }

        const { amount } = req.body;
        if (typeof amount !== 'number' || amount < 0) {
            return res.status(400).json({ error: 'Amount must be a positive number' });
        }

        globalCftSupply = amount;

        res.json({
            success: true,
            message: `CFT supply set to ${amount.toLocaleString()}`,
            availableCft: globalCftSupply
        });
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

// Deduct from CFT supply (after purchase)
router.post('/cft-supply/deduct', authMiddleware, async (req, res) => {
    try {
        const { amount } = req.body;
        if (typeof amount !== 'number' || amount <= 0) {
            return res.status(400).json({ error: 'Amount must be a positive number' });
        }

        if (amount > globalCftSupply) {
            return res.status(400).json({ error: 'Insufficient CFT supply' });
        }

        globalCftSupply -= amount;

        res.json({
            success: true,
            deducted: amount,
            remainingSupply: globalCftSupply
        });
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

export default router;


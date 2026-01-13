import mongoose from 'mongoose';
import crypto from 'crypto';

const userSchema = new mongoose.Schema({
    // Basic Info
    email: {
        type: String,
        required: true,
        unique: true,
        lowercase: true,
        trim: true
    },
    password: {
        type: String,
        required: true,
        minlength: 6
    },
    name: {
        type: String,
        required: true,
        trim: true
    },

    // Role - determines which dashboard they access
    role: {
        type: String,
        required: true,
        enum: ['STARTUP', 'INVESTOR', 'VALIDATOR', 'PLATFORM'],
        uppercase: true
    },

    // Organization/User ID (used in chaincode)
    orgUserId: {
        type: String,
        unique: true,
        sparse: true
    },

    // Startups owned by this user (STARTUP role)
    startups: [{
        startupId: { type: String, required: true },     // STU_{orgUserId}_{seq}
        displayId: { type: String },                      // S-001, S-002
        name: { type: String, required: true },
        description: { type: String },
        createdAt: { type: Date, default: Date.now },
        campaigns: [{
            campaignId: { type: String, required: true }, // CAMP_{startupId}_{seq}
            displayId: { type: String },                   // C-001
            projectName: { type: String },
            status: { type: String, default: 'DRAFT' },
            validationStatus: { type: String, default: 'NOT_SUBMITTED' },
            assignedValidatorId: { type: String },
            assignedPlatformId: { type: String },
            createdAt: { type: Date, default: Date.now }
        }]
    }],

    // Assigned work queue (VALIDATOR & PLATFORM roles)
    assignedQueue: [{
        campaignId: { type: String, required: true },
        startupId: { type: String },
        projectName: { type: String },
        type: { type: String, enum: ['VALIDATION', 'PUBLISH'], required: true },
        assignedAt: { type: Date, default: Date.now }
    }],

    // Completed work history
    completedTasks: [{
        campaignId: { type: String, required: true },
        projectName: { type: String },
        type: { type: String, enum: ['VALIDATION', 'PUBLISH'] },
        result: { type: String },  // APPROVED, REJECTED, PUBLISHED
        completedAt: { type: Date, default: Date.now }
    }],

    // Investments (INVESTOR role)
    investments: [{
        investmentId: { type: String },
        campaignId: { type: String },
        amount: { type: Number },
        currency: { type: String },
        status: { type: String },
        investedAt: { type: Date, default: Date.now }
    }],

    // Wallet Info (CFT/CFRT balances synced from chaincode)
    wallet: {
        cftBalance: { type: Number, default: 0 },
        cfrtBalance: { type: Number, default: 0 },
        frozenCft: { type: Number, default: 0 },
        lastSynced: { type: Date }
    },

    // ML Rating (synced from chaincode)
    mlRating: {
        overallScore: { type: Number, default: 70 },
        trustScore: { type: Number, default: 70 },
        feeTier: { type: String, default: 'STANDARD' },
        blacklisted: { type: Boolean, default: false }
    },

    // Status
    isActive: { type: Boolean, default: true },
    isVerified: { type: Boolean, default: false },

    createdAt: { type: Date, default: Date.now },
    updatedAt: { type: Date, default: Date.now }
});

// Generate unique org user ID before save
userSchema.pre('save', function (next) {
    if (!this.orgUserId) {
        const prefix = this.role.substring(0, 3).toUpperCase();
        const uniqueId = crypto.randomBytes(8).toString('hex').toUpperCase();
        this.orgUserId = `${prefix}_${uniqueId}`;
    }
    this.updatedAt = new Date();
    next();
});

// Hash password (simple hash - use bcrypt in production)
userSchema.pre('save', function (next) {
    if (this.isModified('password')) {
        // In production, use bcrypt
        this.password = crypto.createHash('sha256').update(this.password).digest('hex');
    }
    next();
});

// Compare password
userSchema.methods.comparePassword = function (candidatePassword) {
    const hashedCandidate = crypto.createHash('sha256').update(candidatePassword).digest('hex');
    return this.password === hashedCandidate;
};

// Remove password from JSON output
userSchema.methods.toJSON = function () {
    const obj = this.toObject();
    delete obj.password;
    return obj;
};

const User = mongoose.model('User', userSchema);

export default User;

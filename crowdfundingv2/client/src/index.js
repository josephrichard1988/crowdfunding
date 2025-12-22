import express from 'express';
import cors from 'cors';
import morgan from 'morgan';
import 'dotenv/config';
import config from './settings/index.js';
import connectDB from './database/db.js';

// Import routes
import startupRoutes from './routes/startup.routes.js';
import validatorRoutes from './routes/validator.routes.js';
import platformRoutes from './routes/platform.routes.js';
import investorRoutes from './routes/investor.routes.js';
import authRoutes from './routes/auth.routes.js';

const app = express();

// Connect to MongoDB
connectDB();

// Middleware
app.use(cors());
app.use(express.json());
app.use(morgan('dev'));

// Health check
app.get('/api/health', (req, res) => {
    res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

// API Routes
app.use('/api/auth', authRoutes);
app.use('/api/startup', startupRoutes);
app.use('/api/validator', validatorRoutes);
app.use('/api/platform', platformRoutes);
app.use('/api/investor', investorRoutes);

// Error handling middleware
app.use((err, req, res, next) => {
    console.error('Error:', err.message);
    res.status(err.status || 500).json({
        error: err.message || 'Internal Server Error',
        details: config.nodeEnv === 'development' ? err.stack : undefined,
    });
});

// Start server
const PORT = process.env.PORT || config.port;
app.listen(PORT, () => {
    console.log(`🚀 Crowdfunding API running on port ${PORT}`);
    console.log(`📡 Fabric Channel: ${config.fabric.channelName}`);
    console.log(`📦 Chaincode: ${config.fabric.chaincodeName}`);
});

export default app;


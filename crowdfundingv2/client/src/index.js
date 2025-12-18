import express from 'express';
import cors from 'cors';
import morgan from 'morgan';
import config from './settings/index.js';

// Import routes
import startupRoutes from './routes/startup.routes.js';
import validatorRoutes from './routes/validator.routes.js';
import platformRoutes from './routes/platform.routes.js';
import investorRoutes from './routes/investor.routes.js';

const app = express();

// Middleware
app.use(cors());
app.use(express.json());
app.use(morgan('dev'));

// Health check
app.get('/api/health', (req, res) => {
    res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

// API Routes
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
app.listen(config.port, () => {
    console.log(`🚀 Crowdfunding API running on port ${config.port}`);
    console.log(`📡 Fabric Channel: ${config.fabric.channelName}`);
    console.log(`📦 Chaincode: ${config.fabric.chaincodeName}`);
});

export default app;

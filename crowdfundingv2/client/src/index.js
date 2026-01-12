import express from 'express';
import cors from 'cors';
import morgan from 'morgan';
import 'dotenv/config';
import config from './settings/index.js';
import connectDB from './database/db.js';

// Import auth routes ONLY
import authRoutes from './routes/auth.routes.js';

const app = express();

// Connect to MongoDB
connectDB();

// Middleware
app.use(cors());
app.use(express.json());
app.use(morgan('dev'));

// Health check
app.get('/health', (req, res) => {
    res.json({
        status: 'ok',
        service: 'Authentication Server',
        timestamp: new Date().toISOString()
    });
});

// Auth Routes ONLY (no Fabric routes)
app.use('/api/auth', authRoutes);

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
    console.log(`🔐 Authentication Server running on port ${PORT}`);
    console.log(`📊 MongoDB: Connecting.....`);
    console.log(`⚠️  This server handles AUTH ONLY - Fabric API is on port 4000`);
});

export default app;

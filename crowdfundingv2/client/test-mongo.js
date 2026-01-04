// Test MongoDB connection
// Run this script to test MongoDB connection: node test-mongo.js

import mongoose from 'mongoose';
import dotenv from 'dotenv';

dotenv.config();

const testConnection = async () => {
    try {
        console.log('🔍 Testing MongoDB Atlas Connection...\n');

        const uri = process.env.MONGODB_URI;

        if (!uri) {
            console.error('❌ MONGODB_URI not found in .env file!');
            process.exit(1);
        }

        // Show partial URI for debugging (hide password)
        const uriPreview = uri.includes('@')
            ? uri.split('@')[0].substring(0, 30) + '***@' + uri.split('@')[1]
            : uri.substring(0, 50) + '...';

        console.log('📝 Connection URI preview:', uriPreview);
        console.log('📏 Connection string length:', uri.length, 'characters\n');

        // Check for common issues
        if (uri.includes('"') || uri.includes("'")) {
            console.warn('⚠️  WARNING: Connection string contains quotes - remove them from .env!');
        }

        if (!uri.includes('mongodb+srv://') && !uri.includes('mongodb://')) {
            console.error('❌ Invalid connection string format');
            process.exit(1);
        }

        // Attempt connection
        console.log('🔌 Attempting to connect...\n');
        await mongoose.connect(uri, {
            serverSelectionTimeoutMS: 10000,
            socketTimeoutMS: 45000,
        });

        console.log('✅ MongoDB Atlas connected successfully!');
        console.log('📊 Database:', mongoose.connection.name);
        console.log('🌐 Host:', mongoose.connection.host);

        await mongoose.disconnect();
        console.log('\n✅ Test complete - MongoDB connection is working!');
        process.exit(0);

    } catch (error) {
        console.error('\n❌ MongoDB Connection Failed!\n');
        console.error('Error type:', error.name);
        console.error('Error message:', error.message);

        if (error.message.includes('querySrv ETIMEOUT')) {
            console.error('\n🔧 Common fixes for DNS timeout:');
            console.error('   1. Check MongoDB Atlas IP whitelist (allow 0.0.0.0/0 for development)');
            console.error('   2. Verify connection string is complete (not truncated)');
            console.error('   3. Check if password has special characters (need URL encoding)');
            console.error('   4. Test connection using MongoDB Compass first');
        }

        console.error('\n📖 See mongodb_troubleshooting.md for detailed help');
        process.exit(1);
    }
};

testConnection();

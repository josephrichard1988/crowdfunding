import mongoose from 'mongoose';

const connectDB = async () => {
    try {
        const mongoURI = process.env.MONGODB_URI;

        await mongoose.connect(mongoURI);

        console.log('✅ MongoDB connected successfully');

        mongoose.connection.on('error', (err) => {
            console.error('MongoDB connection error:', err);
        });

        mongoose.connection.on('disconnected', () => {
            console.warn('MongoDB disconnected');
        });

    } catch (error) {
        console.error('❌ MongoDB connection failed:', error.message);
        console.warn('⚠️ Running without MongoDB - auth features disabled');
    }
};

export default connectDB;

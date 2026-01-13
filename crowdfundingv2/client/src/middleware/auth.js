import crypto from 'crypto';

// Simple JWT implementation (use jsonwebtoken in production)
const JWT_SECRET = process.env.JWT_SECRET || 'default-secret-change-me';
const JWT_EXPIRES_IN = process.env.JWT_EXPIRES_IN || '7d';

// Generate JWT token
export const generateToken = (payload) => {
    const header = { alg: 'HS256', typ: 'JWT' };
    const now = Math.floor(Date.now() / 1000);
    const exp = now + (7 * 24 * 60 * 60); // 7 days

    const tokenPayload = {
        ...payload,
        iat: now,
        exp: exp
    };

    const base64Header = Buffer.from(JSON.stringify(header)).toString('base64url');
    const base64Payload = Buffer.from(JSON.stringify(tokenPayload)).toString('base64url');

    const signature = crypto
        .createHmac('sha256', JWT_SECRET)
        .update(`${base64Header}.${base64Payload}`)
        .digest('base64url');

    return `${base64Header}.${base64Payload}.${signature}`;
};

// Verify JWT token
export const verifyToken = (token) => {
    try {
        const [base64Header, base64Payload, signature] = token.split('.');

        const expectedSignature = crypto
            .createHmac('sha256', JWT_SECRET)
            .update(`${base64Header}.${base64Payload}`)
            .digest('base64url');

        if (signature !== expectedSignature) {
            return null;
        }

        const payload = JSON.parse(Buffer.from(base64Payload, 'base64url').toString());

        // Check expiration
        if (payload.exp && payload.exp < Math.floor(Date.now() / 1000)) {
            return null;
        }

        return payload;
    } catch (error) {
        return null;
    }
};

// Auth middleware
export const authMiddleware = async (req, res, next) => {
    try {
        const authHeader = req.headers.authorization;

        if (!authHeader || !authHeader.startsWith('Bearer ')) {
            return res.status(401).json({ error: 'No token provided' });
        }

        const token = authHeader.substring(7);
        const decoded = verifyToken(token);

        if (!decoded) {
            return res.status(401).json({ error: 'Invalid or expired token' });
        }

        req.user = decoded;
        next();
    } catch (error) {
        res.status(401).json({ error: 'Authentication failed' });
    }
};

// Role guard middleware
export const roleGuard = (allowedRoles) => {
    return (req, res, next) => {
        if (!req.user) {
            return res.status(401).json({ error: 'Not authenticated' });
        }

        if (!allowedRoles.includes(req.user.role)) {
            return res.status(403).json({ error: 'Access denied' });
        }

        next();
    };
};

export default { generateToken, verifyToken, authMiddleware, roleGuard };

import * as jwt from 'jsonwebtoken';
import * as crypto from 'crypto';

const JWT_SECRET = 'your-echobetz-internal-secret';

// Generates the internal ECHOBETZ session token
export const generateSessionToken = (payload: any): string => {
    return jwt.sign(payload, JWT_SECRET, { expiresIn: '1h' });
};

export const verifySessionToken = (token: string): any | null => {
    try {
        return jwt.verify(token, JWT_SECRET);
    } catch (e) {
        return null;
    }
};

// Generates HMAC for internal S2S wallet calls (using the Partner Secret)
export const generateHmac = (payload: any, secret: string): string => {
    const data = JSON.stringify(payload);
    return crypto.createHmac('sha256', secret).update(data).digest('hex');
};

// Simplified partner database lookup
const PARTNERS = {
    'CASINO_ALPHA': {
        secret: 'your-secure-shared-secret', // Must match Casino's secret
        wallet_api_url: 'http://localhost:3000/api/wallet'
    }
};

export const getPartnerSecret = (id: string) => PARTNERS[id] ? PARTNERS[id].secret : null;
export const getPartnerWalletUrl = (id: string) => PARTNERS[id] ? PARTNERS[id].wallet_api_url : null;
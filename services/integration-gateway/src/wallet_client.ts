import axios from 'axios';
import { generateHmac, getPartnerSecret, getPartnerWalletUrl } from './auth'; // Imported from auth.ts
import * as crypto from 'crypto';

// Minimal Transaction Response structure
interface WalletResponse {
    status: string;
    balance: number;
    error?: string;
}

/**
 * Sends a signed debit request to the Casino's Wallet API.
 */
export async function debitExternalWallet(
    txId: string, 
    partnerId: string, 
    playerId: string, 
    amount: number
): Promise<WalletResponse> {
    const partnerSecret = getPartnerSecret(partnerId);
    const walletUrl = getPartnerWalletUrl(partnerId);
    
    if (!partnerSecret || !walletUrl) {
        throw new Error(`Partner configuration missing for ${partnerId}`);
    }

    const payload = {
        transaction_id: txId, 
        partner_id: partnerId,
        player_id: playerId,
        game_id: 'AURORA_STAR',
        amount: amount,
        currency: 'EUR',
        transaction_type: 'BET',
        // Nonce and Timestamp for Replay Protection (optional but recommended)
        nonce: crypto.randomBytes(16).toString('hex'),
        timestamp: Date.now() 
    };
    
    const signature = generateHmac(payload, partnerSecret);

    try {
        const response = await axios.post(`${walletUrl}/debit`, payload, {
            headers: {
                'X-Echobetz-Signature': signature,
                'Content-Type': 'application/json'
            }
        });
        
        return response.data as WalletResponse;
    } catch (error: any) {
        console.error('External Debit Failed:', error.response?.data || error.message);
        throw new Error(`Debit failed. Casino Error: ${error.response?.statusText || error.message}`);
    }
}

/**
 * Sends a signed credit request to the Casino's Wallet API.
 */
export async function creditExternalWallet(
    txId: string,
    relatedTxId: string,
    partnerId: string,
    playerId: string,
    amount: number
): Promise<WalletResponse> {
    const partnerSecret = getPartnerSecret(partnerId);
    const walletUrl = getPartnerWalletUrl(partnerId);

    // ... (Error checks for secret/URL omitted for brevity) ...

    const payload = {
        transaction_id: txId,
        related_transaction_id: relatedTxId,
        partner_id: partnerId,
        player_id: playerId,
        game_id: 'AURORA_STAR',
        amount: amount,
        currency: 'EUR',
        transaction_type: 'WIN',
        nonce: crypto.randomBytes(16).toString('hex'),
        timestamp: Date.now() 
    };

    const signature = generateHmac(payload, partnerSecret);

    try {
        const response = await axios.post(`${walletUrl}/credit`, payload, {
            headers: {
                'X-Echobetz-Signature': signature,
                'Content-Type': 'application/json'
            }
        });
        
        return response.data as WalletResponse;
    } catch (error: any) {
        console.error('External Credit Failed:', error.response?.data || error.message);
        throw new Error(`Credit failed. Casino Error: ${error.response?.statusText || error.message}`);
    }
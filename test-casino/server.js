const express = require('express');
const bodyParser = require('body-parser');
const crypto = require('crypto');
const app = express();
const port = 3000;

// Shared Secret (Must match the one used by ECHOBETZ Wallet Service)
const PARTNER_SECRET = 'your-secure-shared-secret'; 
const PLAYER_BALANCE = { 'PLAYER_UUID_HASHED_9876': 50000 };
const PROCESSED_TX = {}; // In-memory Idempotency Store

app.use(bodyParser.json());

// --- ECHOBETZ Wallet API Implementation (The Casino's API) ---

// HMAC Validation Middleware
const validateHmac = (req, res, next) => {
    const signature = req.headers['x-echobetz-signature'];
    const payload = JSON.stringify(req.body);
    const calculatedHmac = crypto.createHmac('sha256', PARTNER_SECRET)
        .update(payload)
        .digest('hex');

    if (calculatedHmac !== signature) {
        console.error("HMAC Mismatch!");
        return res.status(401).send({ error: "Invalid signature" });
    }
    next();
};

app.post('/api/wallet/debit', validateHmac, (req, res) => {
    const { transaction_id, player_id, amount } = req.body;

    // Idempotency Check
    if (PROCESSED_TX[transaction_id]) {
        return res.status(200).json({ status: "ok", balance: PLAYER_BALANCE[player_id], message: "Idempotent response" });
    }
    
    // Debit Logic
    if (PLAYER_BALANCE[player_id] < amount) {
        return res.status(400).json({ error: "Insufficient funds" });
    }

    PLAYER_BALANCE[player_id] -= amount;
    PROCESSED_TX[transaction_id] = true;
    console.log(`DEBIT: TxID=${transaction_id}, Player=${player_id}, Amount=${amount}, New Balance=${PLAYER_BALANCE[player_id]}`);
    
    res.status(200).json({ status: "ok", balance: PLAYER_BALANCE[player_id] });
});

app.post('/api/wallet/credit', validateHmac, (req, res) => {
    const { transaction_id, player_id, amount } = req.body;

    // Idempotency Check
    if (PROCESSED_TX[transaction_id]) {
        return res.status(200).json({ status: "ok", balance: PLAYER_BALANCE[player_id], message: "Idempotent response" });
    }

    // Credit Logic
    PLAYER_BALANCE[player_id] += amount;
    PROCESSED_TX[transaction_id] = true;
    console.log(`CREDIT: TxID=${transaction_id}, Player=${player_id}, Amount=${amount}, New Balance=${PLAYER_BALANCE[player_id]}`);
    
    res.status(200).json({ status: "ok", balance: PLAYER_BALANCE[player_id] });
});

// --- Test Casino Launch Page ---

// Function to generate the Launch URL signature
const generateLaunchToken = (params, secret) => {
    // In a real system, the token covers all query parameters.
    const data = params.partner_id + params.player_id + params.game_code;
    return crypto.createHmac('sha256', secret).update(data).digest('hex');
};

app.get('/launch', (req, res) => {
    const params = {
        partner_id: 'CASINO_ALPHA',
        game_code: 'AURORA_STAR',
        player_id: 'PLAYER_UUID_HASHED_9876', // Abstract ID
        currency: 'EUR',
        mode: 'REAL',
        ts: Date.now().toString()
    };

    params.token = generateLaunchToken(params, PARTNER_SECRET);
    
    const launchUrl = `http://localhost:8080/launch?${new URLSearchParams(params).toString()}`;
    
    res.send(`
        <h1>ECHOBETZ Test Casino</h1>
        <p>Current Balance: ${PLAYER_BALANCE['PLAYER_UUID_HASHED_9876']} EUR</p>
        <p>Player: ${params.player_id}</p>
        <a href="${launchUrl}" target="_blank">Launch AURORA STAR Slot</a>
    `);
});

app.listen(port, () => {
    console.log(`Test Casino listening on port ${port}`);
});
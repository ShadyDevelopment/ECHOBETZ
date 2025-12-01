import * as express from 'express';
import * as WebSocket from 'ws';
import { createServer } from 'http';
import { URLSearchParams } from 'url';
import { verifySessionToken, generateSessionToken, getPartnerSecret, getPartnerWalletUrl } from './auth';
import { debitExternalWallet, creditExternalWallet } from './wallet_client'; // Assumed Wallet Client

// Simplified gRPC clients (Go services)
const GES_CLIENT = { spin: (req: any) => ({ matrix: [['A','B','C'],['A','B','C'],['A','B','C'],['A','B','C'],['A','B','C']], total_win: 500, status: "ok" }) }; 
const RNG_CLIENT = { getRandomNumbers: (count: number) => ({ numbers: [1, 5, 10, 20, 30], seed: '12345' }) }; // Simplified stub

const app = express();
const server = createServer(app);
const wss = new WebSocket.Server({ server, path: '/ws' });
const PORT = 8080;

// Store active WebSocket connections by session_id
const sessions = new Map<string, WebSocket>();

// --- REST Endpoint: Game Launch (Client Redirection) ---
app.get('/launch', (req, res) => {
    const { partner_id, player_id, game_code, token } = req.query;

    // 1. Validate incoming HMAC token (simplified check here)
    if (!getPartnerSecret(partner_id as string) || token !== 'valid') { // Simplified token check for example
        return res.status(401).send('Launch validation failed.');
    }

    // 2. Generate internal ECHOBETZ session token
    const sessionId = crypto.randomUUID();
    const sessionToken = generateSessionToken({ partner_id, player_id, game_code, sessionId });
    
    // 3. Redirect to the game client with the token embedded
    const launchParams = new URLSearchParams({ token: sessionToken }).toString();
    res.redirect(`/game.html?${launchParams}`);
});

// --- WebSocket Server ---
wss.on('connection', (ws, req) => {
    // Extract internal ECHOBETZ JWT from the URL query
    const urlParams = new URLSearchParams(req.url?.split('?')[1]);
    const token = urlParams.get('token');
    const payload = verifySessionToken(token);

    if (!payload) {
        return ws.close(1008, 'Invalid session token');
    }

    const { sessionId, partner_id, player_id, game_code } = payload;
    sessions.set(sessionId, ws);
    console.log(`Session ${sessionId} connected.`);

    ws.on('message', async (message: string) => {
        const reqData = JSON.parse(message);
        if (reqData.type === 'SPIN_REQUEST') {
            const betAmount = reqData.bet;
            const txId = `ECHOBETZ_SPIN_${Date.now()}_${sessionId}`;
            let totalWin = 0;
            let finalBalance = 0;
            
            try {
                // 1. DEBIT
                const debitResult = await debitExternalWallet(txId, partner_id, player_id, betAmount);
                finalBalance = debitResult.balance;
                
                // 2. RNG Call
                const rngResult = RNG_CLIENT.getRandomNumbers(5); // 5 reels
                
                // 3. GES Call
                const spinResult = GES_CLIENT.spin({
                    game_code, 
                    rngOutputs: rngResult.numbers, 
                    betAmount
                });
                
                totalWin = spinResult.total_win;
                
                // 4. CREDIT (if win > 0)
                if (totalWin > 0) {
                    const creditTxId = `ECHOBETZ_WIN_${Date.now()}_${sessionId}`;
                    const creditResult = await creditExternalWallet(creditTxId, txId, partner_id, player_id, totalWin);
                    finalBalance = creditResult.balance;
                }

                // 5. Send result back to client
                ws.send(JSON.stringify({ 
                    type: "SPIN_RESULT", 
                    payload: { 
                        matrix: spinResult.matrix, 
                        total_win: totalWin, 
                        balance: finalBalance 
                    } 
                }));

            } catch (error: any) {
                console.error('Game round failed:', error.message);
                // Important: Handle refund or error state here
                ws.send(JSON.stringify({ type: "ERROR", message: "Game round failed. Funds may be refunded." }));
            }
        }
    });

    ws.on('close', () => {
        sessions.delete(sessionId);
        console.log(`Session ${sessionId} closed.`);
    });
});

app.use(express.static('services/integration-gateway')); // Serve game.html and assets
server.listen(PORT, () => {
    console.log(`Gateway listening on http://localhost:${PORT}`);
});
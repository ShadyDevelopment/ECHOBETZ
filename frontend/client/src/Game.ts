import * as PIXI from 'pixi.js';
import { NetworkClient } from './NetworkClient';
import { StateMachine, GameState } from './StateMachine';
import { ReelManager } from './ReelManager';
import { UIManager } from './UIManager';

export class Game {
    private app: PIXI.Application;
    private stateMachine: StateMachine;
    private networkClient: NetworkClient;
    private reelManager: ReelManager;
    private uiManager: UIManager;

    constructor(token: string) {
        this.app = new PIXI.Application({ width: 800, height: 600 });
        document.body.appendChild(this.app.view as any);

        this.stateMachine = new StateMachine();
        this.reelManager = new ReelManager(this.app.stage);
        this.uiManager = new UIManager(this.app.stage, this.handleSpinClick);
        
        // Connect to backend
        this.networkClient = new NetworkClient(token, this.handleBackendMessage);

        this.init();
    }

    private init(): void {
        // Load assets (symbols, background) here...
        this.stateMachine.transitionTo(GameState.IDLE);
    }

    private handleSpinClick = () => {
        if (!this.stateMachine.is(GameState.IDLE)) return;
        
        this.stateMachine.transitionTo(GameState.SPINNING);
        const betAmount = this.uiManager.getBetAmount(); // e.g., 100
        this.networkClient.sendSpin(betAmount);
        this.reelManager.startSpinAnimation();
    }

    private handleBackendMessage = (data: any) => {
        if (data.type === "SPIN_RESULT") {
            const { matrix, total_win, balance } = data.payload;
            
            // 1. Stop reels to show the determined outcome
            this.reelManager.stopReelsToMatrix(matrix, () => {
                // 2. Resolve Payouts once animation is done
                this.stateMachine.transitionTo(GameState.PAYOUT);
                this.uiManager.updateBalance(balance);
                
                if (total_win > 0) {
                    this.uiManager.showWin(total_win);
                    // Start win presentation animations...
                }
                
                // 3. Return to Idle
                setTimeout(() => {
                    this.stateMachine.transitionTo(GameState.IDLE);
                }, total_win > 0 ? 3000 : 500); // Wait for win animation
            });
        }
    }
}

// Global initialization from the launch page (index.html)
window.onload = () => {
    // Get the JWT token embedded by the Gateway
    const token = new URLSearchParams(window.location.search).get('token');
    if (token) {
        new Game(token);
    } else {
        document.body.innerHTML = 'Error: Missing session token.';
    }
};
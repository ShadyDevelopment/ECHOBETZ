import * as PIXI from 'pixi.js';

const REEL_WIDTH = 150;
const SYMBOL_SIZE = 120;
const SYMBOL_COLORS: { [key: string]: number } = {
    "S_HIGH_A": 0xFF0000, // Red
    "S_WILD": 0x00FF00, // Green
    "S_SCATTER": 0x0000FF, // Blue
    "S_MID_C": 0xFFFF00, // Yellow
    "S_LOW_E": 0xAAAAAA  // Gray
    // ... add all symbols
};

export class ReelManager extends PIXI.Container {
    private reelContainers: PIXI.Container[] = [];

    constructor(parentStage: PIXI.Container) {
        super();
        parentStage.addChild(this);
        this.createReels();
    }

    private createReels(): void {
        for (let i = 0; i < 5; i++) {
            const reel = new PIXI.Container();
            reel.x = i * REEL_WIDTH;
            this.reelContainers.push(reel);
            this.addChild(reel);
            
            // Initial dummy symbols
            for (let j = 0; j < 3; j++) {
                const symbol = this.createSymbol('S_LOW_E');
                symbol.y = j * SYMBOL_SIZE;
                reel.addChild(symbol);
            }
        }
        this.position.set(50, 50);
    }

    private createSymbol(symbolId: string): PIXI.Graphics {
        const color = SYMBOL_COLORS[symbolId] || 0xCCCCCC;
        const graphics = new PIXI.Graphics();
        graphics.beginFill(color);
        graphics.drawRect(0, 0, SYMBOL_SIZE, SYMBOL_SIZE);
        graphics.endFill();
        return graphics;
    }

    // Placeholder: Starts a visual spin effect
    public startSpinAnimation(): void {
        console.log("Reels spinning...");
        // In a real game: apply high-speed movement or motion blur to symbols
    }

    // Crucial: Stops the reels based on the final matrix from the backend
    public stopReelsToMatrix(matrix: string[][], callback: () => void): void {
        console.log("Stopping reels to matrix:", matrix);
        // matrix is Rows x Reels (e.g., [[R1S1, R2S1, ...], [R1S2, ...], ...])

        // For simplicity, we instantly update the visuals without animation
        for (let r = 0; r < 5; r++) { // Iterate through reels
            const reel = this.reelContainers[r];
            reel.removeChildren(); // Clear old symbols
            
            for (let i = 0; i < 3; i++) { // Iterate through 3 visible rows
                // Note: matrix[i] is the symbol at row i, reel r
                const symbolId = matrix[i][r]; 
                const newSymbol = this.createSymbol(symbolId);
                newSymbol.y = i * SYMBOL_SIZE;
                reel.addChild(newSymbol);
            }
        }
        
        // In a real game: this would involve a timed animation loop for each reel
        callback(); 
    }
}
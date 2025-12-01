import * as PIXI from 'pixi.js';

const REEL_WIDTH = 150;
const SYMBOL_SIZE = 120;
const SYMBOL_COLORS: { [key: string]: number } = {
    "S_HIGH_A": 0xFF0000, // Red
    "S_WILD": 0x00FF00, // Green
    "S_SCATTER": 0x0000FF, // Blue
    "S_MID_C": 0xFFFF00, // Yellow
    "S_LOW_D": 0xFF8800, // Orange
    "S_LOW_E": 0xAAAAAA  // Gray
};

/**
 * Manages the visual representation and animation of the 5x3 reel grid.
 */
export class ReelManager extends PIXI.Container {
    private reelContainers: PIXI.Container[] = [];

    constructor(parentStage: PIXI.Container) {
        super();
        parentStage.addChild(this);
        // Position the entire grid area on the screen
        this.position.set(50, 50); 
        this.createReels();
    }

    private createReels(): void {
        for (let i = 0; i < 5; i++) { // 5 Reels (columns)
            const reel = new PIXI.Container();
            reel.x = i * REEL_WIDTH;
            this.reelContainers.push(reel);
            this.addChild(reel);
            
            // Populate with 3 initial symbols
            for (let j = 0; j < 3; j++) { 
                const symbol = this.createSymbol('S_LOW_E'); // Default symbol
                symbol.y = j * SYMBOL_SIZE;
                reel.addChild(symbol);
            }
        }
    }

    private createSymbol(symbolId: string): PIXI.Graphics {
        const color = SYMBOL_COLORS[symbolId] || 0xCCCCCC;
        const graphics = new PIXI.Graphics();
        graphics.beginFill(color);
        graphics.drawRect(0, 0, SYMBOL_SIZE, SYMBOL_SIZE);
        graphics.endFill();
        return graphics;
    }

    /**
     * Starts a visual spin animation (placeholder).
     */
    public startSpinAnimation(): void {
        console.log("Reels spinning...");
        // Implement complex visual animation logic here (e.g., using PIXI.ticker or GSAP)
    }

    /**
     * Stops the reels and updates the display to match the final symbol matrix from the backend.
     * @param matrix The 3x5 symbol grid ([Row][Reel]).
     * @param callback Function to execute after the visual stop sequence is complete.
     */
    public stopReelsToMatrix(matrix: string[][], callback: () => void): void {
        console.log("Stopping reels to matrix:", matrix);
        
        for (let r = 0; r < 5; r++) { // Iterate through Reels
            const reel = this.reelContainers[r];
            reel.removeChildren(); // Clear old symbols
            
            for (let i = 0; i < 3; i++) { // Iterate through 3 visible Rows
                // matrix[i][r] gets the symbol at Row i, Reel r
                const symbolId = matrix[i][r]; 
                const newSymbol = this.createSymbol(symbolId);
                newSymbol.y = i * SYMBOL_SIZE;
                reel.addChild(newSymbol);
            }
        }
        
        callback(); // Signal that the visual update is complete
    }
}
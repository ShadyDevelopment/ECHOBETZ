import * as PIXI from 'pixi.js';

/**
 * Manages the user interface elements: balance, bet, and spin button.
 */
export class UIManager extends PIXI.Container {
    private balanceText: PIXI.Text;
    private spinButton: PIXI.Graphics;
    private betAmount: number = 100;
    
    constructor(parentStage: PIXI.Container, spinCallback: () => void) {
        super();
        parentStage.addChild(this);
        
        // Balance Display
        this.balanceText = new PIXI.Text('Balance: ---', { fill: 0xFFFFFF, fontSize: 24 });
        this.balanceText.position.set(50, 550);
        this.addChild(this.balanceText);

        // Spin Button
        this.spinButton = this.createButton(spinCallback);
        this.spinButton.position.set(700, 540);
        this.addChild(this.spinButton);
    }

    private createButton(callback: () => void): PIXI.Graphics {
        const btn = new PIXI.Graphics();
        btn.beginFill(0xCC33CC); 
        btn.drawRect(0, 0, 80, 40);
        btn.endFill();
        
        const text = new PIXI.Text('SPIN', { fill: 0xFFFFFF, fontSize: 18 });
        text.anchor.set(0.5);
        text.position.set(40, 20); 
        btn.addChild(text);
        
        btn.interactive = true;
        btn.buttonMode = true;
        btn.on('pointerdown', callback);
        return btn;
    }

    /**
     * Updates the displayed player balance.
     * @param newBalance The updated balance value.
     */
    public updateBalance(newBalance: number): void {
        this.balanceText.text = `Balance: ${newBalance.toFixed(2)} EUR`;
    }
    
    /**
     * Displays a win notification (placeholder).
     * @param amount The win amount.
     */
    public showWin(amount: number): void {
        console.log(`WIN: ${amount} EUR!`);
        // Add win animation/text display here
    }
    
    public getBetAmount(): number {
        return this.betAmount;
    }
}
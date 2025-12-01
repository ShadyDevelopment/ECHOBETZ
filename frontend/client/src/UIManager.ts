import * as PIXI from 'pixi.js';

export class UIManager extends PIXI.Container {
    private balanceText: PIXI.Text;
    private spinButton: PIXI.Graphics;
    private betAmount: number = 100;
    
    constructor(parentStage: PIXI.Container, spinCallback: () => void) {
        super();
        parentStage.addChild(this);
        
        this.balanceText = new PIXI.Text('Balance: ---', { fill: 0xFFFFFF });
        this.balanceText.position.set(50, 500);
        this.addChild(this.balanceText);

        this.spinButton = this.createButton(spinCallback);
        this.spinButton.position.set(700, 500);
        this.addChild(this.spinButton);
    }

    private createButton(callback: () => void): PIXI.Graphics {
        const btn = new PIXI.Graphics();
        btn.beginFill(0xCC33CC);
        btn.drawRect(0, 0, 80, 40);
        btn.endFill();
        
        const text = new PIXI.Text('SPIN', { fill: 0xFFFFFF, fontSize: 18 });
        text.position.set(10, 10);
        btn.addChild(text);
        
        btn.interactive = true;
        btn.buttonMode = true;
        btn.on('pointerdown', callback);
        return btn;
    }

    public updateBalance(newBalance: number): void {
        this.balanceText.text = `Balance: ${newBalance.toFixed(2)} EUR`;
    }
    
    public showWin(amount: number): void {
        console.log(`WIN: ${amount} EUR!`);
        // In a real game: play win sound, show big text animation
    }
    
    public getBetAmount(): number {
        return this.betAmount;
    }
}
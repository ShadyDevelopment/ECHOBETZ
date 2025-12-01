// Manages the game's state transitions
export enum GameState {
    INIT = 'INIT',
    IDLE = 'IDLE',
    SPINNING = 'SPINNING',
    RESOLVING = 'RESOLVING',
    PAYOUT = 'PAYOUT',
    BONUS = 'BONUS'
}

export class StateMachine {
    private currentState: GameState = GameState.INIT;
    
    public transitionTo(newState: GameState): void {
        console.log(`State Transition: ${this.currentState} -> ${newState}`);
        this.currentState = newState;
    }
    
    public is(state: GameState): boolean {
        return this.currentState === state;
    }
}
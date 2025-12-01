// Uses WebSocket for persistent, low-latency communication
export class NetworkClient {
    private ws: WebSocket;
    private gameCallback: (data: any) => void;

    constructor(token: string, onMessage: (data: any) => void) {
        // Assume Gateway is running on port 8080
        this.ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
        this.gameCallback = onMessage;

        this.ws.onopen = () => console.log('WebSocket Connected');
        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.gameCallback(data);
        };
        this.ws.onclose = () => console.log('WebSocket Disconnected');
        this.ws.onerror = (error) => console.error('WebSocket Error:', error);
    }

    public sendSpin(betAmount: number): void {
        const request = {
            type: "SPIN_REQUEST",
            bet: betAmount,
            // session_id is derived from the JWT on the server side
        };
        if (this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(request));
        }
    }
}
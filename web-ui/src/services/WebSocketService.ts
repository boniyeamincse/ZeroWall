/**
 * WebSocketService manages real-time telemetry streaming from the ZeroWall backend.
 */
class WebSocketService {
    private socket: WebSocket | null = null;
    private listeners: ((data: any) => void)[] = [];

    connect(url: string) {
        console.log(`Connecting to telemetry hub at ${url}...`);
        this.socket = new WebSocket(url);

        this.socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.listeners.forEach(l => l(data));
        };

        this.socket.onclose = () => {
            console.log("WebSocket disconnected. Retrying in 5s...");
            setTimeout(() => this.connect(url), 5000);
        };
    }

    addListener(callback: (data: any) => void) {
        this.listeners.push(callback);
    }

    send(message: any) {
        if (this.socket?.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(message));
        }
    }
}

export const telemetry = new WebSocketService();

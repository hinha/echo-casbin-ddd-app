// WebSocketService.ts
// This service handles WebSocket connections to the backend

class WebSocketService {
  private socket: WebSocket | null = null;
  private url: string;

  constructor(url: string = 'ws://localhost:8081/ws/users') {
    this.url = url;
  }

  // Connect to the WebSocket server
  connect(): Promise<WebSocket> {
    return new Promise((resolve, reject) => {
      this.socket = new WebSocket(this.url);

      this.socket.onopen = () => {
        console.log('WebSocket connection established');
        resolve(this.socket as WebSocket);
      };

      this.socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        reject(error);
      };

      this.socket.onclose = (event) => {
        console.log('WebSocket connection closed:', event.code, event.reason);
      };
    });
  }

  // Send a message to the WebSocket server
  sendMessage(type: string, payload: any): void {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not connected');
      return;
    }

    const message = JSON.stringify({
      type,
      payload
    });

    this.socket.send(message);
  }

  // Add a message listener
  addMessageListener(callback: (data: any) => void): void {
    if (!this.socket) {
      console.error('WebSocket is not initialized');
      return;
    }

    this.socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        callback(data);
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };
  }

  // Disconnect from the WebSocket server
  disconnect(): void {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }
}

export default WebSocketService;
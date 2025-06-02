// WebSocketService.ts
// This service handles WebSocket connections to the backend

class WebSocketService {
  private socket: WebSocket | null = null;
  private url: string;
  private reconnectAttempts: number = 0;
  private maxReconnectAttempts: number = 5;

  constructor(url: string = "ws://localhost:8081/ws/users") {
    this.url = url;
  }

  // Connect to the WebSocket server
  connect(): Promise<WebSocket> {
    return new Promise((resolve, reject) => {
      try {
        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
          console.log("WebSocket connection established to:", this.url);
          this.reconnectAttempts = 0;
          resolve(this.socket as WebSocket);
        };

        this.socket.onerror = (event: Event) => {
          const errorMessage = this.getErrorMessage(event);
          console.error("WebSocket error:", errorMessage);
          reject(new Error(errorMessage));
        };

        this.socket.onclose = (event: CloseEvent) => {
          const reason = event.reason || "Unknown reason";
          console.log(
            `WebSocket connection closed: Code ${event.code}, Reason: ${reason}`,
          );

          if (
            !event.wasClean &&
            this.reconnectAttempts < this.maxReconnectAttempts
          ) {
            this.handleReconnect();
          }
        };
      } catch (error) {
        console.error("Failed to create WebSocket connection:", error);
        reject(new Error(`Failed to create WebSocket connection: ${error}`));
      }
    });
  }

  // Get meaningful error message from WebSocket event
  private getErrorMessage(event: Event): string {
    if (event.type === "error") {
      // Check if the server is not available
      if (
        this.url.includes("localhost:8081") ||
        this.url.includes("localhost:8080")
      ) {
        return "Cannot connect to WebSocket server. Please ensure the backend server is running on the correct port.";
      }
      return "WebSocket connection failed. Please check your network connection and server availability.";
    }
    return "Unknown WebSocket error occurred";
  }

  // Handle reconnection logic
  private handleReconnect(): void {
    this.reconnectAttempts++;
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 10000); // Exponential backoff

    console.log(
      `Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`,
    );

    setTimeout(() => {
      this.connect().catch((error) => {
        console.error("Reconnection failed:", error);
      });
    }, delay);
  }

  // Send a message to the WebSocket server
  sendMessage(type: string, payload: any): void {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      console.error("WebSocket is not connected");
      return;
    }

    const message = JSON.stringify({
      type,
      payload,
    });

    this.socket.send(message);
  }

  // Add a message listener
  addMessageListener(callback: (data: any) => void): void {
    if (!this.socket) {
      console.error("WebSocket is not initialized");
      return;
    }

    this.socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        callback(data);
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
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

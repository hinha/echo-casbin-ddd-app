// MockWebSocketService.ts
// This service provides a mock WebSocket implementation for development
// when the backend is not available

class MockWebSocketService {
  private url: string;
  private connected: boolean = false;
  private messageListeners: Array<(data: any) => void> = [];

  constructor(url: string = "ws://localhost:8081/ws/users") {
    this.url = url;
  }

  // Simulate WebSocket connection
  connect(): Promise<MockWebSocketService> {
    return new Promise((resolve) => {
      // Simulate connection delay
      setTimeout(() => {
        console.log("Mock WebSocket connection established to:", this.url);
        this.connected = true;
        resolve(this);
      }, 500);
    });
  }

  // Simulate sending a message
  sendMessage(type: string, payload: any): void {
    if (!this.connected) {
      console.error("Mock WebSocket is not connected");
      return;
    }

    console.log("Mock WebSocket sending message:", { type, payload });

    // Simulate server response based on message type
    if (type === "auth_request") {
      setTimeout(() => {
        this.simulateAuthResponse(payload);
      }, 1000);
    }
  }

  // Simulate authentication response
  private simulateAuthResponse(payload: any): void {
    const { username, password } = payload;

    // Mock authentication logic
    const isValidUser = username === "admin" && password === "password";

    const response = {
      type: "auth_response",
      payload: {
        success: isValidUser,
        message: isValidUser ? "Login successful" : "Invalid credentials",
        token: isValidUser ? "mock-jwt-token-" + Date.now() : null,
        user: isValidUser ? { id: 1, username, role: "admin" } : null,
      },
    };

    // Notify all listeners
    this.messageListeners.forEach((listener) => {
      try {
        listener(response);
      } catch (error) {
        console.error("Error in mock message listener:", error);
      }
    });
  }

  // Add a message listener
  addMessageListener(callback: (data: any) => void): void {
    this.messageListeners.push(callback);
  }

  // Remove a message listener
  removeMessageListener(callback: (data: any) => void): void {
    const index = this.messageListeners.indexOf(callback);
    if (index > -1) {
      this.messageListeners.splice(index, 1);
    }
  }

  // Disconnect from mock WebSocket
  disconnect(): void {
    console.log("Mock WebSocket disconnected");
    this.connected = false;
    this.messageListeners = [];
  }

  // Check if connected
  isConnected(): boolean {
    return this.connected;
  }
}

export default MockWebSocketService;

// WebSocketServiceFactory.ts
// Factory to create WebSocket service with automatic fallback to mock service

import WebSocketService from "./WebSocketService";
import MockWebSocketService from "./MockWebSocketService";

export interface IWebSocketService {
  connect(): Promise<any>;
  sendMessage(type: string, payload: any): void;
  addMessageListener(callback: (data: any) => void): void;
  disconnect(): void;
}

class WebSocketServiceFactory {
  private static instance: IWebSocketService | null = null;
  private static useMock: boolean = false;

  static async createService(
    url: string = "ws://localhost:8081/ws/users",
  ): Promise<IWebSocketService> {
    // Return existing instance if available
    if (this.instance) {
      return this.instance;
    }

    // Try real WebSocket service first
    const realService = new WebSocketService(url);

    try {
      await this.testConnection(url);
      console.log("Using real WebSocket service");
      this.instance = realService;
      this.useMock = false;
      return realService;
    } catch (error) {
      console.warn(
        "Real WebSocket service unavailable, falling back to mock service",
      );
      console.log(
        "To use real WebSocket service, ensure backend server is running on:",
        url,
      );
      const mockService = new MockWebSocketService(url);
      this.instance = mockService;
      this.useMock = true;
      return mockService;
    }
  }

  // Test if WebSocket connection is possible
  private static testConnection(url: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const testSocket = new WebSocket(url);
      const timeout = setTimeout(() => {
        testSocket.close();
        reject(new Error("Connection timeout"));
      }, 2000);

      testSocket.onopen = () => {
        clearTimeout(timeout);
        testSocket.close();
        resolve();
      };

      testSocket.onerror = () => {
        clearTimeout(timeout);
        reject(new Error("Connection failed"));
      };
    });
  }

  static isMockService(): boolean {
    return this.useMock;
  }

  static reset(): void {
    if (this.instance) {
      this.instance.disconnect();
      this.instance = null;
    }
    this.useMock = false;
  }
}

export default WebSocketServiceFactory;

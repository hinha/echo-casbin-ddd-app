import React, { useState, useEffect } from "react";
import WebSocketServiceFactory, {
  IWebSocketService,
} from "../services/WebSocketServiceFactory";

interface LoginProps {
  onLoginSuccess: (token: string) => void;
}

const Login: React.FC<LoginProps> = ({ onLoginSuccess }) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [wsService, setWsService] = useState<IWebSocketService | null>(null);
  const [serviceInfo, setServiceInfo] = useState<string>("");

  useEffect(() => {
    // Initialize WebSocket service
    const initializeService = async () => {
      try {
        const service = await WebSocketServiceFactory.createService();
        setWsService(service);

        // Set service info for user
        if (WebSocketServiceFactory.isMockService()) {
          setServiceInfo(
            "Demo mode: Using mock authentication (username: admin, password: password)",
          );
        } else {
          setServiceInfo("Connected to backend server");
        }
      } catch (error) {
        console.error("Failed to initialize WebSocket service:", error);
        setError("Failed to initialize connection service");
      }
    };

    initializeService();

    // Clean up WebSocket connection when component unmounts
    return () => {
      if (wsService) {
        wsService.disconnect();
      }
    };
  }, []);

  const handleLogin = async (e?: React.FormEvent | React.MouseEvent) => {
    if (e) e.preventDefault();
    setError(null);
    setLoading(true);

    if (!wsService) {
      setError("WebSocket service not initialized. Please refresh the page.");
      setLoading(false);
      return;
    }

    try {
      // Connect to WebSocket
      await wsService.connect();

      // Add message listener for authentication response
      wsService.addMessageListener((data) => {
        if (data.type === "auth_response") {
          setLoading(false);

          if (data.payload.success) {
            // Login successful
            onLoginSuccess(data.payload.token);
          } else {
            // Login failed
            setError(data.payload.message || "Authentication failed");
          }
        }
      });

      // Send login request
      wsService.sendMessage("auth_request", {
        username,
        password,
      });
    } catch (err) {
      setLoading(false);

      // Provide more specific error messages
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError(
          "Failed to connect to server. Please ensure the backend server is running and try again.",
        );
      }

      console.error("WebSocket connection error:", err);
    }
  };

  return (
    <div className="login-container">
      <div className="login-wrapper">
        <div className="login-inner">
          <div className="login-form-container">
            <h2>Login</h2>
            {error && <div className="error-message">{error}</div>}
            <form onSubmit={handleLogin}>
              <div className="form-group">
                <label htmlFor="username" className="form-label">
                  Username
                </label>
                <input
                  type="text"
                  id="username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                  className="form-input"
                />
              </div>
              <div className="form-group">
                <label htmlFor="password" className="form-label">
                  Password
                </label>
                <input
                  type="password"
                  id="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  className="form-input"
                />
              </div>
            </form>
            <button
              onClick={handleLogin}
              disabled={loading}
              className="login-button"
            >
              {loading ? "Logging in..." : "Login"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;

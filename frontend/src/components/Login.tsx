import React, { useState, useEffect } from 'react';
import WebSocketService from '../services/WebSocketService';

interface LoginProps {
  onLoginSuccess: (token: string) => void;
}

const Login: React.FC<LoginProps> = ({ onLoginSuccess }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [wsService] = useState(() => new WebSocketService());

  useEffect(() => {
    // Clean up WebSocket connection when component unmounts
    return () => {
      wsService.disconnect();
    };
  }, [wsService]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      // Connect to WebSocket
      await wsService.connect();

      // Add message listener for authentication response
      wsService.addMessageListener((data) => {
        if (data.type === 'auth_response') {
          setLoading(false);
          
          if (data.payload.success) {
            // Login successful
            onLoginSuccess(data.payload.token);
          } else {
            // Login failed
            setError(data.payload.message || 'Authentication failed');
          }
        }
      });

      // Send login request
      wsService.sendMessage('auth_request', {
        username,
        password
      });
    } catch (err) {
      setLoading(false);
      setError('Failed to connect to server. Please try again.');
      console.error('WebSocket connection error:', err);
    }
  };

  return (
    <div className="login-container">
      <h2>Login</h2>
      {error && <div className="error-message">{error}</div>}
      <form onSubmit={handleLogin}>
        <div className="form-group">
          <label htmlFor="username">Username</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit" disabled={loading}>
          {loading ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </div>
  );
};

export default Login;
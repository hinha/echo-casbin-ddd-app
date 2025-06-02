import React, { useState } from 'react';
import Login from './components/Login';

const App: React.FC = () => {
  const [token, setToken] = useState<string | null>(null);

  const handleLoginSuccess = (newToken: string) => {
    setToken(newToken);
    // Store token in localStorage for persistence
    localStorage.setItem('authToken', newToken);
  };

  // Check if user is already logged in
  React.useEffect(() => {
    const storedToken = localStorage.getItem('authToken');
    if (storedToken) {
      setToken(storedToken);
    }
  }, []);

  const handleLogout = () => {
    setToken(null);
    localStorage.removeItem('authToken');
  };

  return (
    <div className="app">
      {token ? (
        <div className="authenticated-content">
          <header>
            <h1>Welcome to the Dashboard</h1>
            <button onClick={handleLogout}>Logout</button>
          </header>
          <main>
            <p>You are now logged in. Your token is: {token.substring(0, 10)}...</p>
          </main>
        </div>
      ) : (
        <Login onLoginSuccess={handleLoginSuccess} />
      )}
    </div>
  );
};

export default App;

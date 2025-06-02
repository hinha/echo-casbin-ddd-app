# Frontend for Echo Casbin DDD App

This is the frontend application for the Echo Casbin DDD App. It provides a login interface that connects to the backend WebSocket API.

## Prerequisites

- Node.js (v14 or later)
- Yarn or npm

## Getting Started

1. Install dependencies:

```bash
cd frontend
yarn install
# or
npm install
```

2. Start the development server:

```bash
yarn dev
# or
npm run dev
```

This will start the development server on port 5173 by default. You can access the application at http://localhost:5173.

## Features

- WebSocket connection to backend server
- Login authentication via WebSocket
- Token-based authentication
- Persistent login state

## Configuration

The WebSocket server URL is configured in `src/services/WebSocketService.ts`. By default, it connects to `ws://localhost:8081/ws/users`. If your backend is running on a different port or host, you'll need to update this URL.

## Project Structure

- `src/components/Login.tsx`: Login component with WebSocket authentication
- `src/services/WebSocketService.ts`: Service for WebSocket communication
- `src/App.tsx`: Main application component
- `src/global.css`: Global styles

## Building for Production

To build the application for production:

```bash
yarn build
# or
npm run build
```

This will generate a production-ready build in the `dist` directory.
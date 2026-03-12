# ZeroWall Web GUI Development Plan

This document outlines the design and implementation strategy for the ZeroWall administrative interface.

## 1. Core Architecture
- **Framework**: React 18+ with TypeScript.
- **Build Tool**: Vite (for fast HMR and optimized production builds).
- **Styling**: Tailwind CSS for layout + Vanilla CSS for high-performance custom animations.
- **State Management**: 
    - **React Query**: For server state and caching (API interactions).
    - **Context API**: For global UI state (theming, sidebar toggle).
- **Routing**: React Router 6.

## 2. Design Philosophy
- **Premium Aesthetics**: Dark mode by default (Slate 900/950).
- **Dynamic Elements**: Glassmorphism effects, smooth transitions, and real-time micro-animations.
- **Data Visualization**: D3.js or Chart.js for real-time traffic and system health metrics.

## 3. Key Components
- **Dashboard**: Real-time widgets for CPU, RAM, Interface throughput, and IDS alerts.
- **Firewall Rule Editor**: Advanced CRUD interface with drag-and-drop prioritization.
- **Log Streamer**: WebSocket-powered live view of system/firewall logs with multi-field filtering.
- **Setup Wizard**: Step-by-step walkthrough for first-time installation (WAN, LAN, Admin password).
- **Certification Manager**: GUI for managing OpenSSL-based CAs, Certificates, and CRLs.

## 4. Implementation Phases
1. **Foundation**: Build the layout shell, sidebar, and authentication context.
2. **Dashboard Logic**: Integrate RRDtool graphs and WebSocket telemetry.
3. **Module UIs**: Individually build Interfaces, Firewall, VPN, and Services pages.
4. **Interactive Tools**: Implement the TUI-like console features and packet capture UI.
5. **Polish**: Responsive design optimization and final UX refinement.

## 5. Integration with Backend
- All communication via RESTful API at `/api/v1/`.
- Real-time updates via WebSocket at `/ws/stats`.
- Authenticated via JWT (stored in HttpOnly cookies or Secure localStorage).

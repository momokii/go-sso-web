# Single Sign-On (SSO) Dashboard Web App

This repository contains the source code for a **Single Sign-On (SSO) Web Application**. The primary goal of this application is to centralize authentication and session management across multiple web applications through a unified dashboard. Users can seamlessly access various apps without needing to log in multiple times, ensuring a consistent and secure experience.  

## Features  
- **Centralized Authentication**: Login once and access multiple applications without re-authentication.  
- **Unified Dashboard**: A simple, user-friendly interface to manage and access all integrated applications.  
- **Central Session Management**:  
  - A **centralized session store** is implemented using PostgreSQL.  
  - The central session acts as a single source of truth for all application sessions.  
  - Each integrated app maintains a local session, which is validated against the central session to ensure synchronization.  
- **Minimalist Frontend**: Built using **HTML**, **Bootstrap**, and **jQuery** for a clean and responsive design.  
- **Backend Framework**: Powered by **Golang Fiber**, ensuring fast and efficient performance.  

## Purpose  
This project was developed as part of my portfolio to showcase a practical implementation of a centralized authentication system with robust session management. The app is designed to simplify user management while ensuring secure and scalable authentication across multiple applications.  

## Tech Stack  
- **Backend**: Golang Fiber framework  
- **Frontend**: HTML, Bootstrap, and jQuery  
- **Database**: PostgreSQL for session storage and centralized session management  
- **Authentication**: Centralized Single Sign-On (SSO) logic  

## How It Works  
1. **Login Process**:  
   - Users authenticate via the SSO application.  
   - A central session is created and stored in the PostgreSQL database.  
   - The session ID is shared securely with the local session for the accessed application.  
2. **Session Management**:  
   - The central session acts as the **master session**, ensuring all local sessions remain synchronized.  
   - If a user logs out or the central session expires, all local sessions for the user are invalidated automatically.  
3. **Dashboard Navigation**: Users can navigate through the dashboard to access available applications seamlessly.  

## How to Run

1. **Clone the repository:**
   ```bash
   git clone github.com/momokii/go-sso-web
   ```

2. **Install Go:**
   - Make sure Go is installed. You can check by running:
     ```bash
     go version
     ```

3. **Run the server:**
   - Start the server using the following command:
     ```bash
     go run main.go
     ```

4. **Optional: Use Air for Hot Reloading:**
   - If you want hot reloading during development, you can use [Air](https://github.com/cosmtrek/air).
   - Start the server with Air by running:
     ```bash
     air
     ```

5. **Access the website:**
   - Open your browser and go to `http://localhost:3001` (or the specified port).

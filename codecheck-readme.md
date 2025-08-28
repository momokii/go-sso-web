# Go CodeCheck Desktop

<p align="center">
  <img src="./frontend/src/assets/images/logo-universal.png" alt="Go CodeCheck Desktop Logo" width="175" style="border-radius: 15px;">
</p>

A comprehensive desktop application for automated security vulnerability scanning and code analysis, built with the Wails framework. This tool provides an intuitive interface for managing code repositories, performing security scans using Semgrep, and analyzing vulnerability trends over time.

## 🚀 Key Features

### 🔐 Authentication & User Management
- **Secure Session-Based Authentication**: Login with username/password credentials
- **First-Time Setup Flow**: Mandatory password change for new users ensures security
- **User Profile Management**: Update username and password through intuitive modals
- **Professional Error Handling**: Enhanced UX with informative modal feedback and loading states

### 📁 Repository Management
- **CRUD Operations**: Complete repository lifecycle management (Create, Read, Update, Delete)
- **Path Validation**: Real-time verification of repository paths before saving
- **Search & Filter**: Advanced filtering and pagination for large repository collections
- **Repository Details**: Store and manage repository metadata including name, description, and path

### 🔍 Security Scanning
- **Semgrep Integration**: Powered by industry-standard Semgrep security analysis engine
- **Docker-Based Execution**: Containerized scanning ensures consistency and isolation
- **Automated Dependency Management**: Automatic Docker image pulling and environment setup
- **Real-Time Scan Progress**: Live updates with scanning stages (preparing, scanning, processing)
- **Comprehensive Vulnerability Detection**: Identifies security issues, code quality problems, and potential vulnerabilities

### 📊 Scan History & Analysis
- **Complete Scan History**: Persistent storage of all scan results with detailed metadata
- **Advanced Filtering**: Filter scans by repository, date range, status, and search terms
- **Detailed Vulnerability Reports**: In-depth analysis of each vulnerability with:
  - Severity levels and impact assessment
  - File locations and line numbers
  - Code snippets and context
  - Semgrep rule information and recommendations
- **Scan Comparison Engine**: Side-by-side comparison of two scans from the same repository showing:
  - **New Vulnerabilities**: Issues introduced since the baseline scan
  - **Fixed Issues**: Vulnerabilities resolved between scans
  - **Persistent Problems**: Unresolved issues requiring attention
  - **Statistical Analysis**: Trend analysis and security posture improvements

### 💾 Data Management
- **SQLite Database**: Embedded database for local data persistence
- **Automated Migrations**: Database schema versioning and automatic updates
- **Relational Data Model**: Properly normalized schema with foreign key relationships
- **Data Integrity**: Transactional operations ensure data consistency

### 🎨 Modern User Interface
- **Svelte Frontend**: Reactive, component-based UI architecture
- **DaisyUI Design System**: Professional, accessible design components
- **TailwindCSS Styling**: Utility-first CSS framework for consistent design
- **Responsive Layout**: Optimized for various screen sizes and resolutions
- **Loading States**: Professional loading overlays and progress indicators
- **Modal System**: Contextual dialogs for user interactions and confirmations

## 🛠 Technical Stack

### Backend
- **Go**: High-performance backend with robust error handling
- **Wails Framework**: Go-based desktop application framework
- **SQLite**: Embedded database for local data storage
- **Docker Integration**: Containerized security scanning with Semgrep
- **Repository Pattern**: Clean architecture with separated concerns

### Frontend
- **Svelte**: Reactive JavaScript framework
- **Vite**: Fast build tool and development server
- **DaisyUI**: Semantic component library

### Infrastructure
- **Docker**: Container runtime for Semgrep execution
- **Semgrep**: Open-source static analysis security scanner
- **Cross-Platform**: Windows, macOS, and Linux support

## 📋 Prerequisites

Before installing Go CodeCheck Desktop, ensure you have:

- **Go 1.19+**: Backend runtime and build tools
- **Node.js 16+**: Frontend development and build process
- **Docker Desktop**: Required for security scanning functionality
- **Wails CLI**: Desktop application framework
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

## 🚀 Installation & Setup

### From Releases

1. Download the latest release for your platform from the [Releases](https://github.com/your-username/go-codecheck-desktop/releases) page
2. Extract the archive (if applicable)
3. Run the application executable
4. Complete the initial setup and authentication flow

### Building from Source

#### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/go-codecheck-desktop.git
cd go-codecheck-desktop
```

#### 2. Install Dependencies
```bash
# Install Go dependencies
go mod tidy

# Install Node.js dependencies
cd frontend
npm install
cd ..
```

#### 3. Database Setup
The application automatically initializes the SQLite database on first run with the required schema including:
- Users table for authentication
- Repositories table for project management
- Scans table for vulnerability tracking

#### 4. Docker Configuration
Ensure Docker Desktop is running. The application will automatically:
- Check Docker availability
- Pull the Semgrep Docker image when needed
- Manage containerized scanning processes

#### 5. Build the Application
```bash
# Development build
wails dev

# Production build
wails build
```

The built application will be in the `build/bin/` directory.

## 🎯 Usage Guide

### Initial Setup
1. **Launch the Application**: Run the built executable or development server
2. **First-Time Login**: Use the default credentials `admin/admin` (username: admin, password: admin) for initial access
3. **Password Setup**: Complete the mandatory password change for security - you will be automatically prompted to change the default password
4. **Environment Check**: Verify Docker and Semgrep availability from the Scan page

### Repository Management
1. **Navigate to Repositories**: Access the repository management interface
2. **Add Repository**: Click "Add Repository" and provide:
   - Repository name and description
   - Full path to the code directory
   - Path validation confirms directory existence
3. **Manage Repositories**: Edit, update, or remove repositories as needed

### Security Scanning
1. **Access Scan Interface**: Navigate to the scanning page
2. **Select Repository**: Choose from your registered repositories
3. **Initiate Scan**: Click "Scan Now" to begin the security analysis
4. **Monitor Progress**: Track real-time scanning progress through multiple stages
5. **Review Results**: Analyze detailed vulnerability reports upon completion

### History & Analysis
1. **View Scan History**: Access comprehensive scan records with filtering options
2. **Detailed Reports**: Click on any scan to view in-depth vulnerability analysis
3. **Compare Scans**: Select two scans from the same repository to:
   - Identify new security issues
   - Track resolved vulnerabilities
   - Monitor security posture trends
   - Generate comparison reports

### User Management
Access user settings to:
- Update login credentials
- Modify user profile information
- Change passwords with validation

## 📁 Project Structure

```
go-codecheck-desktop/
├── backend/
│   ├── core/
│   │   ├── databases/          # SQLite setup and migrations
│   │   ├── docker/            # Docker integration
│   │   ├── handlers/          # API request handlers
│   │   ├── models/            # Data models and structures
│   │   ├── repository/        # Data access layer
│   │   └── semgrep/           # Semgrep integration
│   └── pkg/
│       ├── parser/            # Report parsing utilities
│       └── utils/             # Helper functions
├── frontend/
│   ├── src/
│   │   ├── components/        # Reusable UI components
│   │   ├── stores/            # State management
│   │   └── views/             # Main application pages
│   └── wailsjs/               # Generated Wails bindings
├── build/                     # Build artifacts and resources
└── design/                    # Application design assets
```

## 🔧 Development

### Live Development Mode
```bash
wails dev
```
This starts the application with hot-reload capabilities for both frontend and backend changes. The application is also accessible in your browser at http://localhost:34115 for debugging with devtools.

### Building for Production
```bash
wails build
```
Creates optimized binaries for your platform in the `build/bin/` directory.

### Database Migrations
The application automatically handles database migrations. To add new migrations:
1. Create SQL files in `backend/core/databases/migrations/`
2. Update the migration logic in the database initialization code

### How It Works
1. **Repository Registration**: Users register code repositories with metadata
2. **Docker Integration**: Application manages Semgrep container execution
3. **Temporary Workspace**: Target directories are copied to temporary locations for scanning
4. **Security Analysis**: Semgrep performs comprehensive static analysis
5. **Report Processing**: JSON reports are parsed and stored in the database
6. **Comparative Analysis**: Historical data enables trend analysis and comparison

### Adding New Features
1. **Backend Changes**: Modify Go code in the `backend/` directory
2. **Frontend Changes**: Update Svelte components in `frontend/src/`
3. **API Integration**: Update Wails bindings and handlers as needed
4. **Database Schema**: Add migrations for data model changes

## 🤝 Contributing

We welcome contributions to Go CodeCheck Desktop! Please follow these guidelines:

### Getting Started
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes with appropriate tests
4. Commit your changes: `git commit -m 'Add amazing feature'`
5. Push to the branch: `git push origin feature/amazing-feature`
6. Open a Pull Request

### Code Standards
- Follow Go best practices and formatting (`gofmt`)
- Use meaningful commit messages
- Include documentation for new features
- Ensure all tests pass before submitting
- Maintain consistent code style with existing codebase

### Bug Reports
When reporting bugs, please include:
- Operating system and version
- Go and Node.js versions
- Steps to reproduce the issue
- Expected vs. actual behavior
- Relevant logs or error messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **[Semgrep](https://semgrep.dev/)**: For providing the excellent static analysis security scanner
- **[Wails Team](https://wails.io/)**: For creating an outstanding Go-based desktop application framework
- **[Svelte](https://svelte.dev/)**: For the reactive frontend framework
- **[DaisyUI](https://daisyui.com/)**: For the beautiful and accessible UI components
- **[TailwindCSS](https://tailwindcss.com/)**: For the utility-first CSS framework

---

**Go CodeCheck Desktop** - Making code security analysis accessible, comprehensive, and actionable. 🛡️
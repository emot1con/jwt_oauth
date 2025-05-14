@echo off
echo Starting Auth System...

REM Check command line args for docker option
if "%1"=="docker" (
    echo Running in Docker mode...
    
    REM Build and start Docker containers
    powershell -Command "cd c:\Users\arqan\Documents\github_project\auth && make up_build"
    
    echo Docker containers are starting...
) else (
    echo Running in local mode...
    
    REM Start backend server
    start powershell -Command "cd c:\Users\arqan\Documents\github_project\auth && go run cmd/main.go"
    
    REM Wait a moment for backend to initialize
    timeout /t 5
    
    REM Start frontend server
    start powershell -Command "cd c:\Users\arqan\Documents\github_project\auth\frontend && go run server.go"
    
    echo Both servers are starting...
)

echo Backend: http://localhost:8080
echo Frontend: http://localhost:3000

REM Open browser to the frontend
timeout /t 2
start http://localhost:3000

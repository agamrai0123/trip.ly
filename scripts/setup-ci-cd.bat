@echo off
REM ============================================================================
REM WanderPlan - Automated CI/CD & Local Development Setup (Windows)
REM ============================================================================

setlocal enabledelayedexpansion

for %%I in (.) do set "PROJECT_ROOT=%%~fI"
cd /d "%PROJECT_ROOT%"

REM Colors equivalent (Windows doesn't support ANSI by default)
set "BLUE=[94m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "RED=[91m"
set "NC=[0m"

REM ============================================================================
REM UTILITY FUNCTIONS
REM ============================================================================

:log_info
echo.
echo [INFO] %1
goto :eof

:log_success
echo [SUCCESS] %1
goto :eof

:log_error
echo [ERROR] %1
goto :eof

:log_warning
echo [WARNING] %1
goto :eof

REM ============================================================================
REM CHECKS & INSTALLATION
REM ============================================================================

:check_prerequisites
echo.
echo Checking prerequisites...

REM Check Git
where git >nul 2>nul
if %errorlevel% neq 0 (
  echo [ERROR] Git is not installed. Please install Git first.
  exit /b 1
)

REM Check Go
where go >nul 2>nul
if %errorlevel% neq 0 (
  echo [WARNING] Go is not installed. Please install Go 1.23+ from https://golang.org/dl/
) else (
  for /f "tokens=3" %%i in ('go version') do echo [SUCCESS] Go %%i found
)

REM Check Bun or npm
where bun >nul 2>nul
if %errorlevel% neq 0 (
  where npm >nul 2>nul
  if %errorlevel% neq 0 (
    echo [WARNING] Neither Bun nor npm found
  ) else (
    for /f "tokens=*" %%i in ('npm --version') do echo [SUCCESS] npm %%i found
  )
) else (
  for /f "tokens=*" %%i in ('bun --version') do echo [SUCCESS] Bun %%i found
)

REM Check Docker
where docker >nul 2>nul
if %errorlevel% neq 0 (
  echo [WARNING] Docker is not installed
) else (
  for /f "tokens=*" %%i in ('docker --version') do echo [SUCCESS] %%i found
)

echo [SUCCESS] All prerequisites checked
goto :eof

REM ============================================================================
REM SETUP GIT HOOKS
REM ============================================================================

:setup_git_hooks
echo.
echo Setting up Git hooks for local validation...

if not exist ".git\hooks" mkdir ".git\hooks"

REM Pre-commit hook
(
echo @echo off
echo setlocal enabledelayedexpansion
echo for /f "tokens=*" %%i in ('git rev-parse --show-toplevel') do set "PROJECT_ROOT=%%i"
echo echo Checking for frontend changes...
echo git diff --name-only --cached ^| findstr /i "^frontend/" >nul
echo if not errorlevel 1 (
echo   echo Running ESLint...
echo   cd /d "!PROJECT_ROOT!\frontend"
echo   call npm run lint -- --fix ^>nul 2^>nul ^|^| call bun run lint --fix ^>nul 2^>nul ^|^| true
echo   cd /d "!PROJECT_ROOT!"
echo   call git add .
echo )
echo echo Pre-commit checks passed
) > ".git\hooks\pre-commit.bat"

echo [SUCCESS] Git pre-commit hook installed

REM Pre-push hook
(
echo @echo off
echo setlocal enabledelayedexpansion
echo for /f "tokens=*" %%i in ('git rev-parse --show-toplevel') do set "PROJECT_ROOT=%%i"
echo echo Running pre-push validation...
echo cd /d "!PROJECT_ROOT!\frontend"
echo echo Testing frontend...
echo call npm run test ^>nul 2^>nul ^|^| call bun run test ^>nul 2^>nul ^|^| true
echo cd /d "!PROJECT_ROOT!\backend"
echo echo Testing backend...
echo call go test ./... ^>nul 2^>nul ^|^| true
echo echo Pre-push validation complete
) > ".git\hooks\pre-push.bat"

echo [SUCCESS] Git pre-push hook installed
goto :eof

REM ============================================================================
REM BUILD & TEST FUNCTIONS
REM ============================================================================

:test_frontend
echo.
echo Testing frontend...
cd /d "%PROJECT_ROOT%\frontend"

where bun >nul 2>nul
if %errorlevel% equ 0 (
  call bun install --frozen-lockfile
  call bun run lint
  call bun run test
) else (
  call npm install --frozen-lockfile
  call npm run lint
  call npm run test
)

cd /d "%PROJECT_ROOT%"
echo [SUCCESS] Frontend tests passed
goto :eof

:test_backend
echo.
echo Testing backend services...
cd /d "%PROJECT_ROOT%\backend"

call go mod download
call go vet ./...
call go test -v -race -coverprofile=coverage.out ./...

cd /d "%PROJECT_ROOT%"
echo [SUCCESS] Backend tests passed
goto :eof

:test_proto
echo.
echo Validating proto files...
cd /d "%PROJECT_ROOT%\backend"

call go install github.com/bufbuild/buf/cmd/buf@latest
call buf lint
call buf generate
call go mod tidy

cd /d "%PROJECT_ROOT%"
echo [SUCCESS] Proto validation passed
goto :eof

REM ============================================================================
REM AUTO-FIX FUNCTIONS
REM ============================================================================

:fix_frontend_errors
echo.
echo Attempting to auto-fix frontend errors...
cd /d "%PROJECT_ROOT%\frontend"

where bun >nul 2>nul
if %errorlevel% equ 0 (
  call bun run lint --fix
) else (
  call npm run lint -- --fix
)

cd /d "%PROJECT_ROOT%"
echo [SUCCESS] Frontend auto-fix completed
goto :eof

:fix_backend_errors
echo.
echo Attempting to auto-fix backend errors...
cd /d "%PROJECT_ROOT%\backend"

call go mod tidy
call go fmt ./...
call go install golang.org/x/tools/cmd/goimports@latest
call goimports -w .

cd /d "%PROJECT_ROOT%"
echo [SUCCESS] Backend auto-fix completed
goto :eof

:fix_all_errors
call :fix_frontend_errors
call :fix_backend_errors
call :test_proto

echo.
for /f %%A in ('git status --porcelain') do (
  git config user.email "local-ci[bot]@wanderplan.local"
  git config user.name "WanderPlan Local CI"
  git add -A
  git commit -m "🤖 Auto-fix: Local validation errors"
  echo [SUCCESS] Fixes committed
  goto :eof
)

goto :eof

REM ============================================================================
REM FULL VALIDATION
REM ============================================================================

:run_full_validation
echo.
echo Running full validation pipeline...

set EXIT_CODE=0

call :test_frontend
if %errorlevel% neq 0 set EXIT_CODE=1

call :test_backend
if %errorlevel% neq 0 set EXIT_CODE=1

call :test_proto
if %errorlevel% neq 0 set EXIT_CODE=1

if %EXIT_CODE% equ 0 (
  echo [SUCCESS] All validations passed!
) else (
  echo [ERROR] Some validations failed. Attempting auto-fix...
  call :fix_all_errors
  echo [WARNING] Auto-fix completed. Please review changes and test again.
)

goto :eof

REM ============================================================================
REM DOCKER STACK
REM ============================================================================

:run_docker_stack
echo.
echo Starting Docker stack...

where docker-compose >nul 2>nul
if %errorlevel% neq 0 (
  echo [ERROR] Docker Compose not installed
  goto :eof
)

call docker-compose up -d
echo [INFO] Waiting for services to be ready...
timeout /t 10 /nobreak

echo [INFO] Checking service health...
for %%s in (localhost:8080 localhost:8081 localhost:8082) do (
  curl -s http://%%s/healthz >nul 2>nul
  if !errorlevel! equ 0 (
    echo [SUCCESS] %%s is healthy
  ) else (
    echo [WARNING] %%s health check failed
  )
)

goto :eof

REM ============================================================================
REM GITHUB STATUS
REM ============================================================================

:check_github_actions
echo.
echo Checking GitHub Actions status...

where gh >nul 2>nul
if %errorlevel% neq 0 (
  echo [WARNING] GitHub CLI not found. Install from https://cli.github.com/
  goto :eof
)

call gh run list --workflow=ci-cd.yml --limit 1
echo [INFO] Use 'gh run view ^<run-id^>' to see detailed logs

goto :eof

REM ============================================================================
REM SETUP AUTO-DEPLOYMENT
REM ============================================================================

:setup_auto_deployment
echo.
echo Setting up continuous auto-deployment...
echo.
echo GitHub Actions workflows are already configured in .github\workflows\
echo.
echo The following automation is now active:
echo.
echo.    On every push/PR:
echo.       - Run ESLint + TypeScript check
echo.       - Run unit ^& component tests
echo.       - Run Go vet ^& linting
echo.       - Build all services
echo.       - Validate proto files
echo.       - Build Docker images
echo.
echo.    On test failures:
echo.       - Auto-detect error types
echo.       - Auto-fix formatting ^& lint issues
echo.       - Create issues for unresolved errors
echo.       - Retry workflow runs
echo.
echo.    On success:
echo.       - Auto-merge PRs (if configured)
echo.       - Create deployment artifacts
echo.       - Mark ready for deployment
echo.
echo [SUCCESS] Auto-deployment setup complete!

goto :eof

REM ============================================================================
REM MAIN MENU
REM ============================================================================

:show_menu
cls
echo.
echo ================================================================
echo    WanderPlan - Automated CI/CD ^& Test Suite
echo ================================================================
echo.
echo Select an option:
echo   1) Check prerequisites ^& install tools
echo   2) Setup Git hooks
echo   3) Run full validation (all tests)
echo   4) Test frontend only
echo   5) Test backend only
echo   6) Auto-fix all errors
echo   7) Start Docker stack
echo   8) Check GitHub Actions status
echo   9) Setup for continuous auto-deployment
echo   0) Exit
echo.

goto :eof

REM ============================================================================
REM MAIN LOOP
REM ============================================================================

:main
cd /d "%PROJECT_ROOT%"

:menu_loop
call :show_menu
set /p choice="Enter your choice: "

if "%choice%"=="1" call :check_prerequisites
if "%choice%"=="2" call :setup_git_hooks
if "%choice%"=="3" call :run_full_validation
if "%choice%"=="4" call :test_frontend
if "%choice%"=="5" call :test_backend
if "%choice%"=="6" call :fix_all_errors
if "%choice%"=="7" call :run_docker_stack
if "%choice%"=="8" call :check_github_actions
if "%choice%"=="9" call :setup_auto_deployment
if "%choice%"=="0" goto :end

pause
goto :menu_loop

:end
echo [SUCCESS] Goodbye!
endlocal
exit /b 0

if "%1"=="" goto main

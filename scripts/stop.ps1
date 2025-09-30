# PowerShell script to stop the application
Write-Host "Stopping RiskNexus application..." -ForegroundColor Yellow

docker-compose down

Write-Host "Application stopped." -ForegroundColor Green

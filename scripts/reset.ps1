# PowerShell script to reset the application (remove all data)
Write-Host "Resetting RiskNexus application..." -ForegroundColor Red
Write-Host "This will remove all data including the database!" -ForegroundColor Red

$confirmation = Read-Host "Are you sure? (y/N)"
if ($confirmation -eq 'y' -or $confirmation -eq 'Y') {
    Write-Host "Stopping and removing all containers and volumes..." -ForegroundColor Yellow
    docker-compose down -v
    docker system prune -f
    
    Write-Host "Application reset complete." -ForegroundColor Green
    Write-Host "Run .\scripts\start.ps1 to start fresh." -ForegroundColor Cyan
} else {
    Write-Host "Reset cancelled." -ForegroundColor Yellow
}

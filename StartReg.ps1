if ($args.Count -eq 0) {
    Write-Host "Please provide the path to bytebrother.exe as an argument."
    Exit
}

$registryKeyPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run"

$programName = "ByteBrother"

# args[0] absolute path to bytebrother.exe
$programPath = Resolve-Path $args[0]

New-ItemProperty -Path $registryKeyPath -Name $programName -Value $programPath -PropertyType String -Force

Write-Host "Added $programPath to startup registry."

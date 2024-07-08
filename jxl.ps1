# Define the file extension and the associated application
$extension = ".jxl"
$appId = "AppX43hnxtbyyps62jhe9sqpdzxn1790zetc"

# Define registry paths
$extensionKey = "HKCU:\Software\Classes\$extension"
$userChoiceKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\FileExts\$extension\UserChoice"

# Create or set the default value for the file extension
if (-not (Test-Path $extensionKey)) {
    New-Item -Path $extensionKey
}
Set-ItemProperty -Path $extensionKey -Name "(default)" -Value "AppX43hnxtbyyps62jhe9sqpdzxn1790zetc"

# Set the user choice for the file extension
if (-not (Test-Path $userChoiceKey)) {
    New-Item -Path $userChoiceKey
}
Set-ItemProperty -Path $userChoiceKey -Name "ProgId" -Value "AppX43hnxtbyyps62jhe9sqpdzxn1790zetc"
Set-ItemProperty -Path $userChoiceKey -Name "Hash" -Value ""

Write-Output "File association for $extension set to Windows Photos"

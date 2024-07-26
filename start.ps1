$host.UI.RawUI.WindowTitle = "Aussteller Sync"

while ($true) {
    Clear-Host
    .\ausstellersync.exe

    $timeout = 4 # Timeout in Stunden
    $timeout = $timeout * 60 * 60 # Timeout in Sekunden
    Timeout /T $timeout
}
Dieses Programm synchronisiert alle Artikel auf dem SAGE für die Seite: https://aussteller.computer-extra.de

Zuerst Go runterladen und installieren
https://go.dev/dl/

Danach das Repo runterladen
https://github.com/computerextra/AusstellerSync

In den Ordner gehen, wo das Zeug runtergeladen wurde:
Terminal in dem Ordner öffnen:
```pwsh
go get .
```
Warten bis fertig.

Danach
```
go build
```

In dem Ordner ist nun eine ```sync.exe``` Datei.

Wenn diese Ausgeführt wird, passiert alles von alleine.

Die oberen Schritte müssen nur einmalig durchgeführt werden, danach kann direkt die ```sync.exe``` geöffnet werden.

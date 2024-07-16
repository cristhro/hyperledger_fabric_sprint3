### Cambios importantes
Ficheros que se han modificado: 
https://github.com/cristhro/hyperledger_fabric_sprint3/compare/52021838230d04edf024c0f67c1bc81bdccf2eac...7faf5e2b866966c2519232e37c4cc9f46095cec8
### 
# Prerequisitos
### Instalar fabric\-samples
```warp-runnable-command
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.3.2 1.5.2

```
### Configuración de Módulo en Go
```warp-runnable-command
cd curso/chaincodes/test/go
go mod init test.go 
go mod tidy
go mod vendor
```
# Ejecución del script \`\`
`test_universidades_chaincode.sh`
Ir al directorio `/curso`
```warp-runnable-command
cd curso
```
Dar permiso de ejecución al script
```warp-runnable-command
chmod +x test_universidades_chaincode.sh

```
Ejecutamos el script
```warp-runnable-command
./test_universidades_chaincode.sh
```

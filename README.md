# The Blues
Trabajo práctico de la materia Teoría del Lenguaje.

## Instrucciones de uso
Es necesario tener un archivo `.env` con el siguiente formato:

```go
USER_DB="..."
PASSWORD_DB="..."
URI_DB="mongodb+srv://%s:%s@..."
```
> `URI_DB` debe tener `%s` (el primero es `USER_DB` y el otro `PASSWORD_DB`) porque se usa con `fmt.Sprintf`

Para iniciar el servidor, correr:
```
go run main.go
```

## Tipos de subasta
- Subasta dinamica ascendente(normal)
- Subasta dinamica descendente(se pone un precio muy elevado y se empiezan a hacer ofertas menores)
- Subasta por tiempo
- Subasta donde el primero paga lo de el segundo
- Subasta de servicio o producto de proveedores, elige la mas baja.

## Aspectos a implementar
- Controlador
- Usuarios
- Subasta
- Conversion de monedas
- UI

## Como implementarlo
- Con las librerias nativas
- Con librerias como gorilla/gin.

## Endpoints
- authUser
- updateOffer
- closeAuction
- createAuction
- getAuctions
- getAuction
- getUsers
- getUser
- createUser

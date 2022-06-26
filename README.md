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
```sh
go run main.go
```

o bien:

```sh
go build # Compila y genera el ejecutable
./tp-tdl
```

## Tipos de subasta
- Subasta dinamica ascendente (normal)
- Subasta por tiempo
- ~~Subasta dinamica descendente(se pone un precio muy elevado y se empiezan a hacer ofertas menores)~~
- ~~Subasta donde el primero paga lo de el segundo~~
- ~~Subasta de servicio o producto de proveedores, elige la mas baja.~~

## Endpoints

URL | Método | Descripción
--|:--:|--
`/` | GET | Pantalla de inicio con login y registro
`/users` | POST | Crea un nuevo usuario
`/login` | POST | Se conecta al sistema
`/profile` | GET | Información del usuario
`/logout` | POST | Se desconecta del sistema
`/create-auction` | GET | Devuelve el formulario para crear una subasta nueva
`/create-auction` | POST | Crea una nueva subasta
`/auctions` | GET | Obtiene todas las subastas disponibles
`/auctions/{auction-id}` | GET | Obtiene la información de una subasta
`/auctions/{auction-id}` | DELETE | Elimina una subasta
`/auctions/{auction-id}` | POST | Actualiza la oferta de la subasta
`/auctions/seller/{auction-id}` | POST | Termina una subasta
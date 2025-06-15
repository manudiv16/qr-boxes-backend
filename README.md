# API REST en Go con integración Clerk

Esta es una API REST simple implementada en Go que se integra con Clerk para autenticación y está diseñada para trabajar con un frontend Astro.

## Estructura del proyecto

```
backend/
├── auth/               # Paquete de autenticación con Clerk
├── utils/              # Utilidades compartidas
├── .env                # Variables de entorno
├── go.mod              # Definición del módulo Go y dependencias
├── go.sum              # Checksums de dependencias
├── main.go             # Archivo principal
└── README.md           # Este archivo
```

## Requisitos

- Go 1.18 o superior
- Cuenta en Clerk y claves de API

## Configuración

1. Asegúrate de tener las variables de entorno configuradas en el archivo `.env`:

```
CLERK_SECRET_KEY=sk_test_yourClerkSecretKey
PORT=8080
FRONTEND_URL=http://localhost:4321
```

## Dependencias

Este proyecto utiliza:

- `github.com/clerkinc/clerk-sdk-go` - SDK oficial de Clerk para Go
- `github.com/rs/cors` - Middleware CORS para Go
- `github.com/joho/godotenv` - Carga de variables de entorno desde archivo `.env`

## Ejecución

Para ejecutar la aplicación:

```bash
cd backend
go mod tidy  # Instala las dependencias
go run main.go
```

El servidor se iniciará por defecto en http://localhost:8080

## Endpoints disponibles

### Endpoints públicos

#### GET /api/health
Comprueba el estado del API.

**Respuesta**:
```json
{
  "success": true,
  "data": {
    "status": "ok",
    "timestamp": "2023-04-21T12:34:56Z",
    "version": "1.0.0"
  }
}
```

### Endpoints protegidos (requieren autenticación)

Todos los endpoints protegidos requieren un encabezado de autorización con un token JWT válido de Clerk:

```
Authorization: Bearer <clerk_session_token>
```

#### GET /api/user/profile
Devuelve los datos del perfil del usuario autenticado.

**Respuesta**:
```json
{
  "success": true,
  "data": {
    "id": "user_123456789",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "timestamp": 1697839645
  }
}
```

## Integración con Frontend

Para integrar esta API con tu frontend Astro + Clerk:

1. Asegúrate de que el valor de `FRONTEND_URL` coincida con la URL donde se ejecuta tu frontend Astro.

2. En tu aplicación frontend, usa el token de sesión de Clerk para autenticación:

```javascript
// Ejemplo de cómo hacer una solicitud autenticada desde el frontend
import { useAuth } from '@clerk/clerk-react';

const { getToken } = useAuth();

async function fetchUserProfile() {
  const token = await getToken();
  const response = await fetch('http://localhost:8080/api/user/profile', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return await response.json();
}
```

## Desarrollo

Esta API puede expandirse fácilmente:

- Agrega nuevas rutas en `main.go`
- Implementa controladores adicionales para manejar la lógica de negocio
- Integra una base de datos para almacenamiento persistente
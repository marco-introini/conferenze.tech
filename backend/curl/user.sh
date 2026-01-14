# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"marco@marcointroini.it","password":"password"}'

# Il response contiene un token. Usarlo per chiamate protette:
curl http://localhost:8080/api/me \
  -H "Authorization: Bearer <token>"

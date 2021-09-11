
### Add User
```
vault write auth/userpass/users/nipe \
    password=harefa \
    policies=service
vault login -method=userpass \
    username=nipe \
    password=harefa
```

### Enable Plugin
```sh
vault secrets enable -path=postgres database

vault write postgres/config/products \
    plugin_name=postgresql-database-plugin \
    allowed_roles="*" \
    connection_url="postgresql://{{username}}:{{password}}@hostdb:5432?sslmode=disable" \
    username="postgres" \
    password="password"

vault write postgres/roles/nginx \
  db_name=products \
  creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
  GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"{{name}}\"" \
  default_ttl="1m" \
  max_ttl="24h"

```



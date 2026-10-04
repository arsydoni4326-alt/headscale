# Troubleshooting

This guide covers common issues you may encounter with Headscale and Headplane, along with their solutions.

## Authentication Issues

### Password Authentication (Headplane Web UI)

#### Invalid password error

**Symptoms:**

- Login form shows "Invalid password" error
- HTTP 401 Unauthorized response

**Possible causes and solutions:**

1. **Password mismatch:**
   - Verify the password matches your `config.yaml` or environment variable exactly
   - Check for leading/trailing whitespace in the password
   - Passwords are case-sensitive

2. **Configuration not loaded:**
   - Restart Headscale after changing the password
   - Verify Headscale is reading the correct configuration file
   - Check logs for configuration parsing errors: `journalctl -u headscale -n 50`

3. **Environment variable precedence:**
   - If `HEADSCALE_HEADPLANE_PASSWORD` is set, it overrides `config.yaml`
   - Check environment variables: `env | grep HEADSCALE_HEADPLANE`
   - Unset the environment variable or update it to match your intended password

4. **YAML syntax error:**
   - Verify the `headplane.password` field is correctly indented in `config.yaml`
   - Run `headscale configtest` to validate your configuration
   - Example correct syntax:
     ```yaml
     headplane:
       password: "your-password-here"
     ```

**Diagnostic commands:**

```bash
# Verify configuration is valid
headscale configtest

# Check which config file is being used
headscale version

# Check environment variables
env | grep HEADSCALE_HEADPLANE

# View recent logs
journalctl -u headscale -n 50
```

#### Rate limiting (429 Too Many Requests)

**Symptoms:**

- Login attempts return HTTP 429
- Error message: "Rate limit exceeded. Please try again later."

**Causes:**

- More than 5 failed login attempts within 1 minute from the same IP address

**Solutions:**

1. **Wait and retry:**
   - Wait 60 seconds for the rate limit window to reset
   - Try logging in again with the correct password

2. **Check for automated scripts:**
   - Verify no scripts are making excessive login attempts
   - Review cron jobs or monitoring tools that might be hitting the login endpoint

3. **Shared IP address:**
   - If behind a proxy or NAT, multiple users may share the same IP
   - Consider increasing the rate limit in a future Headscale update (not currently configurable)

4. **Verify password before retrying:**
   - Double-check the password before attempting again
   - Each failed attempt counts toward the rate limit

#### Session expired

**Symptoms:**

- Redirected to login page after being logged in
- HTTP 401 Unauthorized on API requests
- "Session expired" message in UI

**Causes:**

- Sessions expire after 24 hours
- System clock skew between client and server

**Solutions:**

1. **Log in again:**
   - Sessions are designed to expire after 24 hours for security
   - Simply log in again to get a new session token

2. **Check system clock:**
   - Verify both client and server clocks are synchronized
   - Use NTP to keep clocks in sync:
     ```bash
     timedatectl status
     systemctl status systemd-timesyncd
     ```

3. **Browser issues:**
   - Clear browser cache and cookies for the Headscale domain
   - Try logging in with an incognito/private window

#### Password change not taking effect

**Symptoms:**

- Changed password in config but old password still works
- Changed password but can't log in with new one

**Solutions:**

1. **Restart Headscale:**
   - Configuration changes require a restart
   ```bash
   systemctl restart headscale
   ```

2. **Verify configuration precedence:**
   - Environment variable takes precedence over config file
   - If `HEADSCALE_HEADPLANE_PASSWORD` is set, update it or unset it
   ```bash
   # Check current env var
   env | grep HEADSCALE_HEADPLANE_PASSWORD
   
   # Unset if needed
   unset HEADSCALE_HEADPLANE_PASSWORD
   ```

3. **Validate configuration syntax:**
   ```bash
   headscale configtest
   ```

4. **Check which config file is loaded:**
   - Headscale searches multiple paths in order:
     - `/etc/headscale/config.yaml`
     - `~/.headscale/config.yaml`
     - `./config.yaml` (current directory)
   - Verify you're editing the correct file
   - Use `headscale --config /path/to/config.yaml` to specify explicitly

5. **Existing sessions remain valid:**
   - Changing the password does not invalidate existing sessions
   - Users with active sessions can continue using Headplane for up to 24 hours
   - To force logout, restart Headscale or wait for session expiry

### API Key Authentication

#### API key not working

**Symptoms:**

- HTTP 401 Unauthorized when using API key
- "Invalid API key" error

**Solutions:**

1. **Verify API key format:**
   - API keys should be used in the `Authorization: Bearer <KEY>` header
   - Example:
     ```bash
     curl -H "Authorization: Bearer YOUR_API_KEY" \
         https://headscale.example.com/api/v1/user
     ```

2. **Check API key expiry:**
   ```bash
   headscale apikeys list
   ```
   - Look for your API key and check its expiration date
   - Create a new key if expired

3. **API key vs. session token:**
   - Session tokens (from password login) start with `hp_` and expire after 24 hours
   - API keys are longer-lived (default 90 days)
   - Don't confuse the two token types

4. **Verify API key wasn't expired:**
   ```bash
   headscale apikeys list
   ```
   - If the key is no longer in the list, it was expired
   - Create a new API key

## Configuration Issues

### Configuration file not found

**Symptoms:**

- "Configuration file not found" error on startup
- Headscale fails to start

**Solutions:**

1. **Create configuration file:**
   - Download example config:
     ```bash
     wget -O /etc/headscale/config.yaml \
         https://raw.githubusercontent.com/juanfont/headscale/main/config-example.yaml
     ```

2. **Specify config path explicitly:**
   ```bash
   headscale --config /path/to/config.yaml serve
   ```

3. **Use environment variable:**
   ```bash
   export HEADSCALE_CONFIG=/path/to/config.yaml
   headscale serve
   ```

### YAML syntax errors

**Symptoms:**

- "Error parsing configuration file"
- Headscale fails to start with parsing error

**Solutions:**

1. **Validate configuration:**
   ```bash
   headscale configtest
   ```

2. **Common YAML errors:**
   - Incorrect indentation (use 2 spaces, not tabs)
   - Missing quotes around special characters
   - Incorrect nesting of fields

3. **Check headplane section syntax:**
   ```yaml
   # Correct:
   headplane:
     password: "your-password"
   
   # Wrong (missing colon):
   headplane
     password: "your-password"
   
   # Wrong (incorrect indentation):
   headplane:
   password: "your-password"
   ```

### Environment variable not taking effect

**Symptoms:**

- Set `HEADSCALE_HEADPLANE_PASSWORD` but password doesn't work
- Configuration seems to be ignored

**Solutions:**

1. **Verify environment variable is set:**
   ```bash
   env | grep HEADSCALE_HEADPLANE_PASSWORD
   ```

2. **Restart Headscale after setting:**
   ```bash
   export HEADSCALE_HEADPLANE_PASSWORD="your-password"
   systemctl restart headscale
   ```

3. **Check systemd service file:**
   - If running as a systemd service, env vars must be set in the service file
   - Edit `/etc/systemd/system/headscale.service`:
     ```ini
     [Service]
     Environment="HEADSCALE_HEADPLANE_PASSWORD=your-password"
     ```
   - Reload and restart:
     ```bash
     systemctl daemon-reload
     systemctl restart headscale
     ```

4. **Use systemd environment file:**
   ```bash
   # Create /etc/headscale/environment
   echo 'HEADSCALE_HEADPLANE_PASSWORD=your-password' > /etc/headscale/environment
   chmod 600 /etc/headscale/environment
   ```
   
   Update service file:
   ```ini
   [Service]
   EnvironmentFile=/etc/headscale/environment
   ```

## Database Migration Issues

### Missing headplane_users table

**Symptoms:**

- Error in logs: `SQL logic error: no such table: headplane_users (1)`
- HTTP 500 errors when accessing `/api/v1/headplane/users`
- Frontend error when visiting `/admin/admin/users`: `Cannot read properties of undefined (reading 'length')`
- Headscale starts but user management features don't work

**Cause:**

The `headplane_users` table was not created during database initialization or migration. This table is required for the multi-user Headplane authentication system introduced in recent versions.

**How migrations work:**

Headscale uses an automatic migration system that runs on server startup. Migrations are defined in code (`hscontrol/db/db.go`) and tracked in the `migrations` table. The `headplane_users` table should be created automatically by migration ID `202610031721-create-headplane-users`.

**Solution 1: Restart Headscale (let auto-migration retry)**

The simplest fix is to restart Headscale and let the migration system retry:

```bash
# For systemd
systemctl restart headscale

# For Docker
docker restart headscale

# Or kill and restart the process
pkill headscale
headscale serve
```

Check the logs during startup for migration-related messages:

```bash
# For systemd
journalctl -u headscale -f

# For Docker
docker logs -f headscale
```

Look for messages like:
- `Created default admin user from config password`
- `running migration 202610031721-create-headplane-users`

**Solution 2: Manually create the table (SQLite)**

If auto-migration fails, you can manually create the table using SQL. This is a safe operation that will not affect existing data.

1. **Stop Headscale:**

   ```bash
   systemctl stop headscale
   # or
   docker stop headscale
   ```

2. **Locate your database file:**

   Check your `config.yaml` for the database path:
   
   ```bash
   grep -A 5 "^database:" /etc/headscale/config.yaml
   ```
   
   Common locations:
   - `/var/lib/headscale/db.sqlite`
   - `/etc/headscale/db.sqlite`
   - `./db.sqlite` (if running from source)

3. **Backup your database:**

   ```bash
   cp /var/lib/headscale/db.sqlite /var/lib/headscale/db.sqlite.backup-$(date +%Y%m%d-%H%M%S)
   ```

4. **Open the database and check if migration is recorded:**

   ```bash
   sqlite3 /var/lib/headscale/db.sqlite
   ```
   
   Inside sqlite3:
   
   ```sql
   -- Check if migration was attempted
   SELECT * FROM migrations WHERE id = '202610031721-create-headplane-users';
   ```
   
   - If it returns a row, the migration was recorded but the table creation failed
   - If it returns nothing, the migration was never run

5. **Check if the table already exists:**

   ```sql
   SELECT name FROM sqlite_master WHERE type='table' AND name='headplane_users';
   ```
   
   - If it returns `headplane_users`, the table exists (migration issue is elsewhere)
   - If it returns nothing, proceed to create the table

6. **Manually create the headplane_users table:**

   ```sql
   -- Create the headplane_users table
   CREATE TABLE IF NOT EXISTS headplane_users (
       id INTEGER PRIMARY KEY AUTOINCREMENT,
       username TEXT NOT NULL UNIQUE,
       password_hash TEXT NOT NULL,
       role TEXT NOT NULL DEFAULT 'user',
       created_at DATETIME,
       updated_at DATETIME
   );
   
   -- Create index for username lookups
   CREATE UNIQUE INDEX IF NOT EXISTS idx_headplane_users_username 
       ON headplane_users(username);
   ```

7. **Create a default admin user (if needed):**

   If you have a password configured in `config.yaml` under `headplane.password`, you need to create an admin user manually. The password must be bcrypt-hashed.
   
   Generate bcrypt hash for your password (choose one method):
   
   ```bash
   # Exit sqlite3 first (type .quit)
   
   # Method 1: Using htpasswd (if available)
   htpasswd -nBC 12 "" | tr -d ':\n'
   
   # Method 2: Using Python
   python3 -c "import bcrypt; print(bcrypt.hashpw(b'your-password', bcrypt.gensalt(12)).decode())"
   
   # Method 3: Using Go (if available)
   go run -<<'EOF'
   package main
   import (
       "fmt"
       "golang.org/x/crypto/bcrypt"
   )
   func main() {
       hash, _ := bcrypt.GenerateFromPassword([]byte("your-password"), 12)
       fmt.Println(string(hash))
   }
   EOF
   ```
   
   Then insert the admin user:
   
   ```bash
   sqlite3 /var/lib/headscale/db.sqlite
   ```
   
   ```sql
   -- Replace $2a$12$... with your bcrypt hash from above
   INSERT INTO headplane_users (username, password_hash, role, created_at, updated_at)
   VALUES ('admin', '$2a$12$YOUR_BCRYPT_HASH_HERE', 'admin', datetime('now'), datetime('now'));
   ```

8. **Record the migration (if it wasn't recorded):**

   ```sql
   -- Only run this if the migration wasn't in the migrations table
   INSERT OR IGNORE INTO migrations (id) VALUES ('202610031721-create-headplane-users');
   ```

9. **Verify the table and data:**

   ```sql
   -- Check table structure
   .schema headplane_users
   
   -- Check if admin user exists
   SELECT id, username, role FROM headplane_users;
   
   -- Exit sqlite3
   .quit
   ```

10. **Start Headscale:**

    ```bash
    systemctl start headscale
    # or
    docker start headscale
    ```

11. **Verify the fix:**

    ```bash
    # Check logs for errors
    journalctl -u headscale -n 50
    
    # Test the API endpoint (replace YOUR_API_KEY)
    curl -H "Authorization: YOUR_API_KEY" http://localhost:8080/api/v1/headplane/users
    
    # Or visit the Headplane UI and try to log in
    # Navigate to http://your-headscale:8080/admin/admin/users
    ```

**Solution 3: Manually create the table (PostgreSQL)**

If you're using PostgreSQL instead of SQLite:

1. **Connect to your database:**

   ```bash
   psql -h localhost -U headscale -d headscale
   ```

2. **Check if table exists:**

   ```sql
   \dt headplane_users
   ```

3. **Create the table:**

   ```sql
   CREATE TABLE IF NOT EXISTS headplane_users (
       id SERIAL PRIMARY KEY,
       username TEXT NOT NULL UNIQUE,
       password_hash TEXT NOT NULL,
       role TEXT NOT NULL DEFAULT 'user',
       created_at TIMESTAMP,
       updated_at TIMESTAMP
   );
   
   CREATE UNIQUE INDEX IF NOT EXISTS idx_headplane_users_username 
       ON headplane_users(username);
   ```

4. **Create admin user:**

   Generate bcrypt hash (same as SQLite instructions above), then:
   
   ```sql
   INSERT INTO headplane_users (username, password_hash, role, created_at, updated_at)
   VALUES ('admin', '$2a$12$YOUR_BCRYPT_HASH_HERE', 'admin', NOW(), NOW());
   ```

5. **Record migration:**

   ```sql
   INSERT INTO migrations (id) VALUES ('202610031721-create-headplane-users')
   ON CONFLICT DO NOTHING;
   ```

**Verification checklist:**

After applying the fix:

- [ ] `headplane_users` table exists in database
- [ ] At least one admin user exists in the table
- [ ] Migration `202610031721-create-headplane-users` is recorded in `migrations` table
- [ ] Headscale starts without errors
- [ ] No `no such table: headplane_users` errors in logs
- [ ] `/api/v1/headplane/users` endpoint returns HTTP 200
- [ ] Headplane UI `/admin/admin/users` page loads without errors

**Prevention:**

To avoid this issue in the future:

1. **Always backup your database before upgrading:**
   ```bash
   cp /var/lib/headscale/db.sqlite /var/lib/headscale/db.sqlite.backup-$(date +%Y%m%d)
   ```

2. **Check logs during startup:**
   ```bash
   journalctl -u headscale -f
   ```

3. **Verify migrations after upgrade:**
   ```bash
   sqlite3 /var/lib/headscale/db.sqlite "SELECT * FROM migrations ORDER BY id;"
   ```

4. **Keep your database file writable:**
   ```bash
   ls -la /var/lib/headscale/db.sqlite
   # Should be owned by headscale user with read/write permissions
   ```

### Migration failed during upgrade

**Symptoms:**

- Headscale won't start after upgrade
- Error: `migration failed: ...`
- Database is in an inconsistent state

**Solution:**

1. **Restore from backup:**
   
   ```bash
   systemctl stop headscale
   cp /var/lib/headscale/db.sqlite.backup /var/lib/headscale/db.sqlite
   systemctl start headscale
   ```

2. **If no backup exists, try to repair:**
   
   Check the specific migration error in logs and consult the [Missing headplane_users table](#missing-headplane_users-table) section above.

3. **Report the issue:**
   
   If manual table creation doesn't work, file a bug report with:
   - The exact error message from logs
   - Output of `.schema` from sqlite3
   - Output of `SELECT * FROM migrations;`

## Network and Connectivity Issues

### Cannot connect to Headplane web interface

**Symptoms:**

- Browser shows "Connection refused" or "Unable to connect"
- Headplane UI doesn't load

**Solutions:**

1. **Verify Headscale is running:**
   ```bash
   systemctl status headscale
   ```

2. **Check listen address:**
   - Verify `listen_addr` in config.yaml
   - For remote access, use `0.0.0.0:8080` not `127.0.0.1:8080`

3. **Check firewall rules:**
   ```bash
   # Allow port 8080
   ufw allow 8080/tcp
   
   # Or for firewalld
   firewall-cmd --add-port=8080/tcp --permanent
   firewall-cmd --reload
   ```

4. **Verify port is listening:**
   ```bash
   ss -tlnp | grep 8080
   ```

### HTTPS/TLS issues

**Symptoms:**

- "Your connection is not private" browser warning
- TLS certificate errors

**Solutions:**

See the [TLS Configuration Guide](ref/tls.md) for detailed HTTPS setup instructions.

Quick checks:

1. **Verify TLS is configured:**
   - Check `tls_cert_path` and `tls_key_path` in config.yaml
   - Or verify Let's Encrypt config

2. **Test certificate:**
   ```bash
   openssl s_client -connect headscale.example.com:443
   ```

3. **Check certificate expiry:**
   ```bash
   echo | openssl s_client -connect headscale.example.com:443 2>/dev/null | \
       openssl x509 -noout -dates
   ```

## General Troubleshooting Steps

### Enable debug logging

Add to your `config.yaml`:

```yaml
log:
  level: debug
  format: text
```

Restart Headscale and check logs:

```bash
systemctl restart headscale
journalctl -u headscale -f
```

### Check Headscale version

```bash
headscale version
```

Ensure you're running a version that supports password authentication (v0.30.0+).

### Review logs

```bash
# View recent logs
journalctl -u headscale -n 100

# Follow logs in real-time
journalctl -u headscale -f

# View logs for specific time range
journalctl -u headscale --since "2026-10-02 10:00:00"

# Search logs for errors
journalctl -u headscale | grep -i error
```

## Getting Help

If you've tried the solutions above and still have issues:

1. **Check the documentation:**
   - [Authentication Guide](usage/authentication.md)
   - [Configuration Reference](ref/configuration.md)
   - [API Reference](ref/api.md)

2. **Search existing issues:**
   - [Headscale GitHub Issues](https://github.com/juanfont/headscale/issues)

3. **Join the community:**
   - [Discord server](https://discord.gg/c84AZQhmpx) for real-time help

4. **File a bug report:**
   - Include Headscale version (`headscale version`)
   - Include relevant logs (with sensitive info redacted)
   - Include configuration (with passwords/secrets redacted)
   - Describe steps to reproduce the issue



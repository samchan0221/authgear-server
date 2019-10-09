package template

/* #nosec */
const templateResetPasswordHTML = `<!DOCTYPE html>

<head>
  <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
{{ if .error }}
<p>{{ .error }}</p>
{{ end }}
<form method="POST" action="{{ .action_url }}">
  <label for="password">New Password</label>
  <input type="password" name="password"><br>
  <label for="confirm">Confirm Password</label>
  <input type="password" name="confirm"><br>
  <input type="hidden" name="code" value="{{ .code }}">
  <input type="hidden" name="user_id" value="{{ .user_id }}">
  <input type="hidden" name="expire_at" value="{{ .expire_at }}">
  <input type="submit" value="Submit">
</form>`

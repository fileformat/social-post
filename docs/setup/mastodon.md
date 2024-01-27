---
title: Mastodon setup for social-post CLI
---

You need an account at a Mastodon server.  The hostname of the server will go it the `MASTODON_SERVER` environment variable.  Note that it does **not** include any `https://` prefix or `/@...` suffix.

Once you have your account:
1. Preferences
2. Development
3. New Application
4. Fill in:
   * Application name: `social-post CLI`
   * Application website: `https://social-post.marcuse.info/`
   * Redirect URI: don't change (should default to `urn:ietf:wg:oauth:2.0:oob`)
   * Permissions: don't change (should already include `read` and `write`)
5. Save
6. Note the "Your access token" value.  This goes in the `MASTODON_USER_TOKEN` environment variable.
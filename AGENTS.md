## Flutter migration rules

- `mobile-app/` is the Flutter migration target.
- `desktop-app/` is the primary reference for features and API usage.
- The server provides both legacy APIs and API2.
- **All Flutter server communication MUST use API2, consistent with `desktop-app/`.**
- Never use legacy API endpoints in Flutter.
- `mobile-app/legacy-android/` is read-only and is only a reference for legacy mobile behavior.
- Prefer a single Codex agent; do not spawn subagents for Flutter migration.
- Keep migration progress in `mobile-app/MIGRATION.md`.
- Do not modify server, desktop-app, or legacy-android code without explicit approval.

## Git repository isolation

The current working directory is the server Git repository.

`mobile-app/` is a separate, independent Git repository.

For Flutter migration tasks:
- Keep Codex working at the current repository root to access server API2 and `desktop-app/`.
- Only modify and commit files within the `mobile-app/` repository, excluding `mobile-app/legacy-android/`.
- Always use `git -C mobile-app` for Git operations.
- Never stage, commit, reset, clean, rebase, or otherwise modify the server Git repository.
- Never stage or commit unrelated changes.
- Do not push automatically.


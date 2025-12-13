# GitVault

An easy way to align data between GitHub and your computer/server/NAS/*

1. Fetches every repo name (public & private) from GitHub
2. Cross checks if repo already exists, is new or it is soft deleted
    a. If soft deleted, does nothing
    b. If already exists, fetches the changes
    c. If it is new, it clones the repo

Distributed as Docker CLI

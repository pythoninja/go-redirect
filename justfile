#!/usr/bin/env just --justfile

switch branch:
    @echo 'Fetching and switching to {{branch}}'
    git fetch origin && git switch --create '{{branch}}'

master-pull:
    git switch master && git pull

branch-push:
    git push --set-upstream origin HEAD
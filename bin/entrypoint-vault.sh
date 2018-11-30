#!/bin/bash

set -e

if [ ! -z "${VAULT_URL}" ] && [ ! -z "${VAULT_TOKEN}" ]
then
    set -a && source <(curl -k -s -H "X-Vault-Token:${VAULT_TOKEN}" "${VAULT_URL}" | jq -r '.data|to_entries|map("\(.key)=\(.value)")[]') && set +a
fi

exec "$@"
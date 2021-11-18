#!/bin/sh

docker compose exec temporal-admin-tools tctl admin cluster add-search-attributes --name CandidateEmail --type Text
docker compose exec temporal-admin-tools tctl admin cluster add-search-attributes --name BackgroundCheckStatus --type keyword

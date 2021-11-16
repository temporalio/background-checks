#!/bin/sh

tctl admin cluster add-search-attributes --name CandidateEmail --type String
tctl admin cluster add-search-attributes --name BackgroundCheckStatus --type keyword

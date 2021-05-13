#!/bin/bash
export TRANSACTION_POOL_BASE_URL=http://localhost:8010

go install
~/go/bin/transaction-spawner
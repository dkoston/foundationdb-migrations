'use strict'

const fs = require('fs')

fs.writeFileSync('/go/src/github.com/cryptowalkio/goose/node_migration_complete', '1')

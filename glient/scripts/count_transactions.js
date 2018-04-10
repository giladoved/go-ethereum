
function count_transactions() {
  var blockNum = 0

  var trans = 0
  var blocks = 0

  console.log(new Date(), 'Starting to count transactions..')

  var print_every = 1000
  while (true) {
    blck = blockNum++
    // Hack
    if (blck > 200000)
        break
    
    block = web3.eth.getBlock(blck)
    if (!block)
      break

    trans = trans + block.transactions.length
    blocks++
    if (blocks % print_every == 0)
      console.log(new Date(), blocks, 'blocks processed. Number of transactions till block number', blck, ':', trans, 'transactions.')
  }

  console.log(new Date(), 'Exiting process.')
  console.log('Total blocks processed:', blocks, 'Total transactions:', trans)
}


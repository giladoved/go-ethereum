
function count_accounts() {
  var blockNum = 0

  var trans = 0
  var blocks = 0
  var addresses = {}

  console.log(new Date(), 'Starting to count accounts..')

  var print_every = 1000
  while (true) {
    blck = blockNum++
    // Hack
//    if (blck > 200000)
//        break
    
    try {
      block = web3.eth.getBlock(blck)
      if (!block)
        break
    } catch (err) {
      console.log(err)
      break
    } finally {
      // do nothing
    }

    trans = trans + block.transactions.length
    for(var i = 0; i < block.transactions.length; i++) {
      var tx = web3.eth.getTransaction(block.transactions[i])
      if (parseInt(tx.value) > 0) {
        addresses[tx.to] = true
      }
    }
    blocks++
    if (blocks % print_every == 0)
      console.log(new Date(), blocks, 'blocks processed. Number of transactions till block number', blck, ':', trans, 'transactions. Number of accounts:', Object.keys(addresses).length)
  }

  console.log('Total blocks processed:', blocks, 'Total transactions:', trans)  

  //let positiveAddresses = []
  var numAddresses = Object.keys(addresses).length
  console.log(new Date(), 'Total number of accounts:', numAddresses)
  var positiveAddresses = 0
  
  for (address in addresses) {
    try {
      var balance = web3.eth.getBalance(address)
      if (balance > 0) {
        //positiveAddresses.push(address)
        positiveAddresses++
      }
      console.log(address)
    } catch (err) {
      console.log(err)
    } finally {
      // do nothing
    }
  }
  console.log(new Date(), 'Number of accounts with positive balance:', positiveAddresses)
  console.log(new Date(), 'Exiting process.')
}


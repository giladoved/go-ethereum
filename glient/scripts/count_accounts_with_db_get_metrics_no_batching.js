
function count_accounts_with_db_get_metrics() {
  var blockNum = 0

  var trans = 0
  var blocks = 0
  var addresses = {}
  
  var gets = debug.metrics(true).eth.db.chaindata.user.gets.Overall
  var gets_now = 0

  var getsTimerMeter = debug.metrics(true).eth.db.chaindata.user.getsTimerMeter.Overall
  var getsTimerMeterNow = 0
  
  var getBlockGets = 0
  var getTransactionGets = 0
  var getBalanceGets = 0

  var getBlockGetTime = 0
  var getTransactionGetTime = 0
  var getBalanceGetTime = 0

  var zeroTransactionBlocks = 0

  console.log(new Date(), 'Starting to count accounts..')

  var print_every = 1000
  while (true) {
    blck = blockNum++
    
    try {
      block = web3.eth.getBlock(blck)

      gets_now = debug.metrics(true).eth.db.chaindata.user.gets.Overall
      getBlockGets += gets_now - gets
      gets = gets_now
 
      getsTimerMeterNow = debug.metrics(true).eth.db.chaindata.user.getsTimerMeter.Overall
      getBlockGetTime += getsTimerMeterNow - getsTimerMeter
      getsTimerMeter = getsTimerMeterNow

      if (!block)
        break
    } catch (err) {
      console.log(err)
      break
    } finally {
      // do nothing
    }

    trans = trans + block.transactions.length
    if (block.transactions.length == 0) {
    	zeroTransactionBlocks++
    }
    for(var i = 0; i < block.transactions.length; i++) {
      var tx = web3.eth.getTransaction(block.transactions[i])

      gets_now = debug.metrics(true).eth.db.chaindata.user.gets.Overall
      getTransactionGets += gets_now - gets
      gets = gets_now

      getsTimerMeterNow = debug.metrics(true).eth.db.chaindata.user.getsTimerMeter.Overall
      getTransactionGetTime += getsTimerMeterNow - getsTimerMeter
      getsTimerMeter = getsTimerMeterNow

      if (parseInt(tx.value) > 0) {
        addresses[tx.to] = true
      }
    }

    blocks++
    if (blocks % print_every == 0) {
      console.log(new Date(), blocks, 'blocks processed. Number of transactions till block number', blck, ':', trans, 'transactions. Number of accounts:', Object.keys(addresses).length)
      console.log(new Date(), 'Get metrics - getBlockGets:', getBlockGets, 'getTransactionGets:', getTransactionGets, 'getBalanceGets:', getBalanceGets)
      console.log(new Date(), 'Get Timer Meter metrics - getBlockGetTimer:', getBlockGetTime, 'getTransactionGetTimer:', getTransactionGetTime, 'getBalanceGetTimer:', getBalanceGetTime)
    }
  }

  console.log(new Date(), 'Total blocks processed:', blocks, 'Total transactions:', trans, 'Number of blocks with zero transactions:', zeroTransactionBlocks)  

  //let positiveAddresses = []
  var numAddresses = Object.keys(addresses).length
  console.log(new Date(), 'Total number of accounts:', numAddresses)
  var positiveAddresses = 0
  
  for (address in addresses) {
    try {
      var balance = web3.eth.getBalance(address)

      gets_now = debug.metrics(true).eth.db.chaindata.user.gets.Overall
      getBalanceGets += gets_now - gets
      gets = gets_now

      getsTimerMeterNow = debug.metrics(true).eth.db.chaindata.user.getsTimerMeter.Overall
      getBalanceGetTime += getsTimerMeterNow - getsTimerMeter
      getsTimerMeter = getsTimerMeterNow

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
  console.log(new Date(), 'Get metrics - getBlockGets:', getBlockGets, 'getTransactionGets:', getTransactionGets, 'getBalanceGets:', getBalanceGets)
  console.log(new Date(), 'Get Timer Meter metrics - getBlockGetTimer:', getBlockGetTime, 'getTransactionGetTimer:', getTransactionGetTime, 'getBalanceGetTimer:', getBalanceGetTime)
  console.log(new Date(), 'Number of accounts with positive balance:', positiveAddresses)
  console.log(new Date(), 'Exiting process.')
}


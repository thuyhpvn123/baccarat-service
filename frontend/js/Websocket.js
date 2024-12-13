var output = document.getElementById("log-content");
var socket = new WebSocket("ws://127.0.0.1:8999/ws");
var socketActive = false;
var $createResult = document.getElementById('create-result');
var deck ;
var playerArr=[];
// var keyArr=[] //mảng gồm 52 key của tất cả lá bài sau khi generate key
var keyArrAPlayer=[] //mảng gồm 13 key của 13 lá bài của mỗi người chơi sau khi gọi get-key-for-player
var encryptedCards=[]
var encryptedCardTest1,encryptedCardTest2,encryptedCardTest3
var addressPlayer
console.log("Imported");
// * Websocket
// Connect to server successfully
var walletAddress = "aa39344b158f4004cac70bb4ace871a9b54baa1e";

var messageForm = {
  command:"get-wallet-pagination",
  limit:2,
  page:1,
};

socket.onopen = (msg) => {
  socketActive = true;
  socket.send(JSON.stringify({ command: "test", value: "ping" }));

};

// WS connection's closed
socket.onclose = (event) => {
  console.log("WS Connection is closed: ", event);
};

// WS connection having errors
socket.onerror = (error) => {
  console.log("Socket Error: ", error);
};
socket.onmessage = (msg) => {
  var data12 = JSON.parse(msg.data);
  output.innerHTML += "Server: " + msg.data + "\n";
   switch (data12.command) {
    case "deal-cards":
      break;
    // case "get-cards":
    //   document.getElementById("cards").innerHTML = data12.data;
    case "shuffle-card":
      console.log("shuffle-card successfully")
      DealCards()
      break;
    case "get-cards":
      encryptedCards=data12.data

      encryptedCardTest1=encryptedCards[0]
      encryptedCardTest2=encryptedCards[1]
      encryptedCardTest3=encryptedCards[2]
      break;
    case "verified-sign":
      addressPlayer=data12.data
      console.log("addressPlayer",addressPlayer)
      getKeyCards(addressPlayer)

    case "get-key-for-player":
      keyArrAPlayer=data12.data
        console.log("keyArrAPlayer:",keyArrAPlayer)
      break;
  }

}
var sendMessage = (msg) => {
  console.log(msg);
  socket.send(JSON.stringify(msg));

};
// //shuffle deck- mảng key gồm 52 key
var deal = () => {
    var ms={
        "roomNumber":"42"
      }

    var Msg = {
        command: "deal-cards",
        value: ms,  
    }
    sendMessage(Msg);
    console.log("deal cards")
};
var $dealCards = document.getElementById('dealCards');

$dealCards.addEventListener('click', async(e) => {

    e.preventDefault()
      try{
        await deal();
        console.log("deal cards successfully ")
        

      }catch{
        console.log(e)
        $createResult.innerHTML = `Ooops... there was an error while trying to encrypt cards`;
      }
    // }
  })
//   
var getCards =(value)=>{
  console.log("address:",value)
  var inputs=JSON.stringify(
    {
      "internalType": "address",
      "name": "player",
      "type": "address",
      "value":value
    });
    var functionInputs=[inputs] ;

    var dataCall = 
    {
      'from':   "45c75cfb8e20a8631c134555fa5d61fcf3e602f2",
      'priKey': "36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb",
      'to':    contract,
      amount:            "",
      "function-name":"getPlayerCards",
      inputArray:functionInputs,
      gas:1000000,
      gasPrice:10,
      timeUse:1000,
      relatedAddresses:[],
    }
  ;
  GetCardsMessage(dataCall);
}
var GetCardsMessage = (ms) => {

  var setMsg = {
      command: "get-cards",
      value: ms,  
  }
  sendMessage(setMsg);
};

  //call function getKeys
var $getKey = document.getElementById('getKey');

$getKey.addEventListener('submit', async(e) => {
  
      e.preventDefault()
      var flag =1,publicKey,hash,sign
      pubKey = $('#publicKey').val()
      hash = $('#hash').val()
      sign = $('#sign').val()
      // addressPlayer = $('#addressPlayer').val()
      if( publicKey ==''|| hash ==''|| sign ==''){
        flag=0
        $('.error_getKey').html("Please type player address")
      }else{
        $('.error_getKey').html("")
      }
      if(flag==1  ){
        try{
          getKeyCards(pubKey,hash,sign)
        }catch{
          console.log(e)
          $createResult.innerHTML = `Ooops... there was an error while trying to set players`;
        }
      }
  })
  var getKeyCards =(pubKey,hash,sign)=>{
    var ms ={
      "hash":hash,
      "sign":sign,
      "pubKey":pubKey,
      "roomNumber":"42"
    }
    var setMsg = {
      command: "get-key-for-player",
      value: ms,  
    }
    sendMessage(setMsg);
  }
//  var getKeyCards =(addressPlayer)=>{
//   console.log("get Keys")
//   var inputs=JSON.stringify(
//     {
//       "internalType": "address",
//       "name": "player",
//       "type": "address",
//       "value":addressPlayer
//     });

//   var functionInputs=[inputs] ;

//   var callData ={
      
//       'from':   "e730d4572f20a4d701ebb80b8b5afa99b36d5e49",
//       'priKey': "224759c1c17f62960245721b2d8b6f2623295c6bfb397ca8fe35f95e55756aa3",
//       'to':    contract,
//       amount:            "",
//       "function-name":"getPlayerCards",
//       inputArray:functionInputs,
//       gas:1000000,
//       gasPrice:10,
//       timeUse:1000,
//       relatedAddresses:[],
//     };
//   GetKeyPlayerMessage(callData);
//   console.log("get key message sent ")

//  }
//  var GetKeyPlayerMessage = (ms) => {
    
//   var setMsg = {
//       command: "get-key-for-player",
//       value: ms,  
//   }
//   sendMessage(setMsg);
// };




//test decrypt
var $test = document.getElementById('test');

$test.addEventListener('click', async(e) => {
  document.getElementById('encryptedCard1').value=encryptedCardTest1
  document.getElementById('encryptedCard2').value=encryptedCardTest2
  document.getElementById('encryptedCard3').value=encryptedCardTest3

  document.getElementById('key1').value=keyArrAPlayer[0]
  document.getElementById('key2').value=keyArrAPlayer[1]
  document.getElementById('key3').value=keyArrAPlayer[2]

})
//decrypt cards
var decrypt = (ms) => {

  var setDecryptMsg = {
      command: "decrypt-cards",
      value: ms,  
  }
  sendMessage(setDecryptMsg);
};
var $decryptCards = document.getElementById('decrypt');

$decryptCards.addEventListener('submit', async(e) => {

    e.preventDefault()
    console.log("decrypt cards")
    var card1,card2,card3,key1,key2,key3
    card1 = $('#encryptedCard1').val()
    card2 = $('#encryptedCard2').val()
    card3 = $('#encryptedCard3').val()

    key1 = $('#key1').val()
    key2 = $('#key2').val()
    key3 = $('#key3').val()
    var keyArrDecode =[]
    var cardArr=[]
    console.log("encrypted cards1:",card1)
    console.log("key1:",key1)
    keyArrDecode.push(key1,key2,key3)
    cardArr.push(card1,card2,card3)
    var ms = {
      "encrypted-cards": cardArr,
      "encrypted-keys": keyArrDecode,
    };
    console.log("playerCards:",cardArr)
    console.log("keyArr:",keyArrDecode)
    // for (i=0;i<4;i++){
    //   try{
        decrypt(ms);
        console.log("decrypt cards sent ")
    //   }catch{
    //     console.log(e)
    //     $createResult.innerHTML = `Ooops... there was an error while trying to encrypt cards`;
    //   }
    // }
})

//get sign
var $sign = document.getElementById('getSign');
$sign.addEventListener('submit', async(e) => {
    
  e.preventDefault()
  var 
  privateKey = $('#privateKey').val()
  address = $('#addressSign').val()

  var ms={
    "privateKey":privateKey,
    "address":address
  }
  var setDecryptMsg = {
    command: "get-sign",
    value: ms,  
}
sendMessage(setDecryptMsg);
})
//admin comfirm result
var comfirmMsg = (roomNumber) => {
  var ms={
      "roomid":roomNumber
    }

  var Msg = {
      command: "comfirm",
      value: ms,  
  }
  sendMessage(Msg);
  console.log("comfirm")
};
var $comfirm = document.getElementById('comfirm');

$comfirm.addEventListener('submit', async(e) => {

  e.preventDefault()
  var roomNumber
  roomNumber = $('#roomidComfirm').val()
    try{
      await comfirmMsg(roomNumber);
      console.log("comfirm successfully ")
      

    }catch{
      console.log(e)
      $createResult.innerHTML = `Ooops... there was an error while trying to comfirm`;
    }
  // }
})


// var callVerify={
//   "hash":"335d3df52df4852c962277ca563e4a7cb67f450a31f9b3fe80543f262bf3f71c",
//   "sign":"a72c0b66b936b36525fdc1f07c94266a8280abd6acb04d1d87477d8ea790c98593084649f5be7fb7653c8ffcc34dd3e5067b2ffac7feaa85e5e574c972b7363dc0ac7df937c82cdce2ccae7cf6c2df265a3e8c9188cda23d872da08a0e2fa5a4",
//   "pubKey":"927c8f340255d5a2d7b839ce4859f7dc33c64e5967fe40dde023c0751a53dff3925eec1eaaeb60f83c203c85e9c5b223",
// }
// var messageForm11 = {
//   command:"verify-json-app",
//   value: callVerify,
// };
// addressS := common.FromHex(addressString)
// hashPub := crypto.Keccak256([]byte(pubKey))
// fmt.Println("hash:", hex.EncodeToString(hashPub))
// hash = fmt.Sprintf("%x", crypto.Keccak256([]byte(pubKey)))


// connect wallet
// var $connectNode = document.getElementById('register');

// $connectNode.addEventListener('submit', async(e) => {

//     e.preventDefault()
//     console.log("connect node")
//     var flag =1,
//     fromAddress = $('#address').val()
//     prikey = $('#prikey').val()

//     if( fromAddress ==''){
//       flag=0
//       $('.error_address').html("Please type address of the account")
//     }else{
//       $('.error_address').html("")
//     }
//     if( prikey ==''){
//       flag=0
//       $('.error_prikey').html("Please type private key of the account")
//     }else{
//       $('.error_prikey').html("")
//     }

//     if(flag==1  ){
//       try{
//         await connectWallet(fromAddress,prikey);
//         console.log("connected wallet ")
//         $createResult.innerHTML =` connect wallet ${fromAddress} successfully` ;
//         $('#wallet-id-user').html(`${fromAddress}`) ;
//       }catch{
//         console.log(e)
//         $createResult.innerHTML = `Ooops... there was an error while trying to register wallet`;
//       }
//     }
//   })
 

// var connectWallet = (address,priKey) => {
//     var wallet={
//         "address":address,
//         "priKey":priKey
//       }
//     var connectWallet = {
//         command: "connect-wallet",
//         value: wallet,  
//     }
//     sendMessage(connectWallet);
//     console.log("connectWallet")
// };
// shuffleCards in contract
// var shuffleCards =()=>{

//     var functionInputs=[] ;

//     var dataCall = 
//     {
//       // 'from':   "45c75cfb8e20a8631c134555fa5d61fcf3e602f2",
//       // 'priKey': "36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb",
//       // 'to':    contract,
//       amount:            "",
//       "function-name":"shuffleCards",
//       inputArray:functionInputs,
//       // gas:1000000,
//       // gasPrice:10,
//       // timeUse:1000,
//       relatedAddresses:[],
//     }
//   ;
//   ShuffleCardsMessage(dataCall);
// }
// var ShuffleCardsMessage = (ms) => {

//   var setMsg = {
//       command: "shuffle-cards",
//       value: ms,  
//   }
//   sendMessage(setMsg);
// };
// getPlayerCards

// 
//generate keys
// var $generateKeys = document.getElementById('generateKey');

// $generateKeys.addEventListener('click', async(e) => {

//     e.preventDefault()
//     console.log("generate keys")
//       try{
//           generate();
//           console.log("generated keys for player ",i+1)
        
//       }catch{
//         console.log(e)
//         $createResult.innerHTML = `Ooops... there was an error while trying to generate keys`;
//       }
// })

//deal cards
  // //call function deal cards
  // var $dealCards = document.getElementById('dealCards');

  // $dealCards.addEventListener('click', async(e) => {
  
  //     e.preventDefault()
     
  //       try{
  //       await shuffleCards()
        

  //       }catch{
  //         console.log(e)
  //         $createResult.innerHTML = `Ooops... there was an error while trying to deal Cards`;
  //       }
      
  // })
//   var DealCards =()=>{
//     console.log("deal cards")
//     var functionInputs=[] ;

//     var dataCall = 
//       {
//         'from':   "45c75cfb8e20a8631c134555fa5d61fcf3e602f2",
//         'priKey': "36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb",
//         'to':    contract,
//         amount:            "",
//         "function-name":"dealCards",
//         inputArray:functionInputs,
//         gas:1000000,
//         gasPrice:10,
//         timeUse:1000,
//         relatedAddresses:[],
//       }
//     ;
//     setDealMessage(dataCall);
//     console.log("deal Cards successfully ")
//   }
//   var setDealMessage = (ms) => {

//     var setMsg = {
//         command: "deal-cards",
//         value: ms,  
//     }
//     sendMessage(setMsg);
// };

// var setDeck =(deck)=>{
  //   var inputs=JSON.stringify(
  //     {
  //       "internalType": "string[]",
  //       "name": "cardsArr",
  //       "type": "string[]",
  //       "value":deck
  //     });
  //     var functionInputs=[inputs] ;
  
  //     var dataCall = 
  //     {
  //       'from':   "45c75cfb8e20a8631c134555fa5d61fcf3e602f2",
  //       'priKey': "36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb",
  //       'to':    contract,
  //       amount:            "",
  //       "function-name":"setDeck",
  //       inputArray:functionInputs,
  //       gas:1000000,
  //       gasPrice:10,
  //       timeUse:1000,
  //       relatedAddresses:[],
  //     }
  //   ;
  //   setDeckMessage(dataCall);
  // }
  // var setDeckMessage = (ms) => {
  
  //     var setMsg = {
  //         command: "set-Deck",
  //         value: ms,  
  //     }
  //     sendMessage(setMsg);
  // };
//   //call function set players
//   var $setPlayers = document.getElementById('setPlayers');

//   $setPlayers.addEventListener('submit', async(e) => {
  
//       e.preventDefault()
//       var flag =1,player1,player2,player3,player4
//       player1 = $('#player1').val()
//       player2 = $('#player2').val()
//       player3 = $('#player3').val()
//       player4 = $('#player4').val()
  
//       if( player1 ==''|| player2 ==''|| player3 ==''|| player4 ==''){
//         flag=0
//         $('.error_player').html("Please type player address")
//       }else{
//         $('.error_player').html("")
//       }
//       if(flag==1  ){
//         try{
//           console.log("set Players")
          
//           var inputs=JSON.stringify(
//           {
//             "internalType": "address[]",
//             "name": "addrs",
//             "type": "address[]",
//             "value":[player1,player2,player3,player4]
//           });
//           var functionInputs=[inputs] ;

//           var dataCall = 
//           {
//             'from':   "45c75cfb8e20a8631c134555fa5d61fcf3e602f2",
//             'priKey': "36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb",
//             'to':    contract,
//             amount:            "",
//             "function-name":"setPlayers",
//             inputArray:functionInputs,
//             gas:1000000,
//             gasPrice:10,
//             timeUse:1000,
//             relatedAddresses:[],
//           }
//         ;
//         await setPlayerMessage(dataCall);
//         console.log("set Player successfully ")

//         }catch{
//           console.log(e)
//           $createResult.innerHTML = `Ooops... there was an error while trying to set players`;
//         }
//       }
//   })
//   var setPlayerMessage = (ms) => {

//     var setMsg = {
//         command: "set-players",
//         value: ms,  
//     }
//     sendMessage(setMsg);
// };
// var generate = () => {
  //     // var ms={
  //     //     "numPlayers":number
  //     //   }
  //     var generateMsg = {
  //         command: "generate-keys",
  //         value: "",  
  //     }
  //     sendMessage(generateMsg);
  //     console.log("generate-keys")
  // };
  
  
//SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.4;

contract ThreeCardsMain {
    uint8 constant NUM_CARDS = 52;
    uint constant cardsPerPlayer = 3;
    // uint8 private limit = 10;
    address  public admin;
    uint public totalRoom; //total room created
    uint public roomNumber=1;
    uint timeLimit = 1 minutes;
    enum State{CREATED,STARTED,STOPPED,COMFIRMED,CLOSED}

    struct Room
    {
        uint roomID; 
        uint8  NUM_PLAYERS ;
        string[] deck;
        uint  bid ;
        State state;
        uint8 numberComfirm;
        address  bidder;
        uint registerEndTime;
        address [] attendants; //all players join
        mapping(address => Player) players;
        Player [] playersArr;
        address [] chosenPlayers;
        uint playersAmount; //max 17 người
        uint limitTime;
    }
    mapping(uint=>mapping(address => uint)) public registeredFunds;
    mapping(address => uint)withdrawFunds;
    struct RoomCopy
    {
        uint roomID; 
        uint registerEndTime;
        uint8  NUM_PLAYERS ;
        bytes[] deck; 
        uint  bid ;
        State state;
        uint8 numberComfirm;
        address  bidder;
        address [] attendants; //all players join 
        address [] chosenPlayers; 
    }

    mapping(uint => Room) private roomStructs;
    uint[]  public totalActiveRoom; //

    mapping(address => uint[]) private userLive; //  mapping địa chỉ player ra những room đang tham gia
    mapping(address => uint) private mTotalUserLive; //mapping địa chỉ player ra tổng số player đang trong các room đang live
    mapping(address => uint[]) private userHistory; //mapping đia chỉ player ra mảng gồm những room id của user đó
    mapping(address => uint8) private mTotalUserHistory;// mapping địa chỉ player ra tổng số lần add history của user

    struct Player
    {
        uint playerId;
        address addr;
        string[] cards;
        string[]keys;
        string[]decodedCards;
        uint []ranks;
        uint totalRank;
        uint lastDigit;
        bool win;
    }
    struct PlayerCopy
    {
        uint playerId;
        address addr;
        bytes[] cards;
        bytes[]keys;
        bytes[]decodedCards;
        uint []ranks;
        uint totalRank;
        uint lastDigit;
        bool win;
    }

    // State public state=State.CREATED;
 
    //The time for joining has already ended.
    error RegisterTimeEnded();
    // The time for joining has not ended yet.
    error RegisterTimeNotYetEnded();
    modifier onlyAdmin(){
        require(msg.sender == admin,"only admin can call this");
        _;
    }
    constructor() payable { 
        admin = msg.sender;
    }
    function setNewAdmin(address _newAdmin)public onlyAdmin
    {
        admin = _newAdmin;
    }
    function setTimeLimit(uint _newTimeLimit) public onlyAdmin{
        timeLimit = _newTimeLimit;
    }
    function getChosenPlayerAddr(uint _roomNumber)public view returns(address[] memory){
        return roomStructs[_roomNumber].chosenPlayers;
    }
    event CreatedRoom(uint roomNumber,uint bidAmount);
    function CreateRoom(uint registerTime,uint bidAmount, uint _playersAmount)public payable returns (uint)
        {
        uint registerEndTime = block.timestamp + registerTime;
        require(msg.value == bidAmount* (_playersAmount-1),"money value is not true");
        require(msg.sender != admin,"admin are not allowed to join");
        require(_playersAmount <=17,"maximum 17 players in the game");
        Room storage room = roomStructs[roomNumber];
        {
            room.roomID = roomNumber;
            room.registerEndTime = registerEndTime;
            room.bidder = msg.sender;
            room.bid= bidAmount;
            room.state= State.CREATED;
            room.NUM_PLAYERS =0 ;
            room.numberComfirm=0 ;
            room.playersAmount = _playersAmount;
        }
        register(roomNumber);
        addRoom(roomNumber);
        emit CreatedRoom(roomNumber,bidAmount);
        roomNumber++;
        return (roomNumber-1);

    }
    // function SetPage(uint8 _number) external onlyAdmin {
    //     limit = _number;
    // }
       //Create New Game Room
    function addRoom(uint  _roomNumber) private returns(bool success)
    {   
        totalRoom++;
        //add new room in array
        totalActiveRoom.push(_roomNumber);
        //erase closed room from array
        sortForActiveRoom();
        return true;
    }
    function increaseRegisterTime(uint  _roomNumber,uint extraTime)public  onlyAdmin
    {
        Room storage room = roomStructs[_roomNumber];
        require(block.timestamp < room.registerEndTime, "only can increase time when joining time has not over yet");
        room.registerEndTime += extraTime;
    }
    //To Add player in Room
    event Registerd(uint roomNumber,address player);
    function register(uint  _roomNumber)public payable returns(bool message)
    {
        Room storage room = roomStructs[_roomNumber];
        require(room.state == State.CREATED,"invalid state");
        address sender = msg.sender;
        if (block.timestamp > room.registerEndTime){
        revert RegisterTimeEnded();
        }
        if (sender != room.bidder ){
        require(msg.value == room.bid,"neccessary to send bid amount");
        }
        require(sender != admin,"admin are not allowed to join");
        require(room.playersArr.length <room.playersAmount,"over maximum players in the game");
        require(room.players[sender].playerId == 0,"already set this address");

                //To check Room is created or not!
        uint roomNum = room.roomID;
        if(roomNum != 0){
            uint playerCounter = room.playersArr.length;
            if(playerCounter <= 17){
                room.players[sender].addr = sender;
                room.players[sender].playerId = room.NUM_PLAYERS+1;
                room.players[sender].cards = new string[](3);
                //update room info
                room.playersArr.push(room.players[sender]);
                room.attendants.push(room.players[sender].addr);
                registeredFunds[_roomNumber][msg.sender] +=  msg.value;
                room.NUM_PLAYERS++;
                //add player to mapping userLive and erase other players from closed room
                AddLive(_roomNumber,sender);
                emit Registerd(_roomNumber,sender);
                return true;
            }
            else{
                revert("Enter in new Room, its full!!");
            }
        }
        else{
            revert("Room is not created!!");
        }

    }
    //if after register time no player register except for bidder
    function closeRoom(uint  _roomNumber) public 
    {
        Room storage room = roomStructs[_roomNumber];
        require(block.timestamp > room.registerEndTime,"close room only after register time");
        require(room.state == State.CREATED,"invalid state");
        require(room.playersArr.length == 1,"only room with no player except for bidder can be closed ");
        require(room.bidder == msg.sender,"only the bidder can close his room");
        room.state = State.CLOSED;
        // sortForLive(room.bidder);
        // sortForActiveRoom();
        transferFund(_roomNumber);
    }
    // when bidder dont choose players to compare in limited time afer setDeck
    function playerCloseRoom(uint  _roomNumber) public
    {
        Room storage room = roomStructs[_roomNumber];
        require(room.state == State.STARTED,"invalid state");
        require(block.timestamp>room.limitTime && room.chosenPlayers.length == 0,"only call when bidder didnt choose players in time limit");
        require(room.players[msg.sender].addr == msg.sender && msg.sender != room.bidder,"only players of room can call this");
        room.state = State.CLOSED;
        // sortForLive(room.bidder);
        // sortForActiveRoom();      
        transferFund(_roomNumber);
    }

    event SetDeck(uint roomNumber,string result);
    function setDeck(uint  _roomNumber,string[] memory cardsArr)public onlyAdmin returns(bool)
    {         
        Room storage room = roomStructs[_roomNumber];
        require(room.roomID >=1,"room number does not exist");
        require(block.timestamp > room.registerEndTime,"set deck only after register time over");
        require(room.state == State.CREATED,"invalid state");
        require(cardsArr.length == NUM_CARDS, "Insufficient cards in the deck.");
        room.deck = cardsArr;
        //shuffle
        uint256 seed = uint(keccak256(abi.encodePacked(block.timestamp, msg.sender)));      
        for (uint8 i = 51; i > 0; i--) {
            uint8 j = uint8(uint256(keccak256(abi.encodePacked(block.timestamp, msg.sender, seed, i)))) % (i + 1);
            string memory temp = room.deck[i];
            room.deck[i] = room.deck[j];
            room.deck[j] = temp;
        }
        // Deal cards starting from the dealer's right
        uint dealer = 0; // Dealer's index
        uint cardIndex = 0;
        for (uint card = 0; card < cardsPerPlayer; card++) {
            for (uint player = dealer; player < room.NUM_PLAYERS; player++) {
                // playerCards[player][card] = deck[cardIndex];
                room.playersArr[player].cards[card] = room.deck[cardIndex];
                cardIndex++;
            }
        }
        room.state = State.STARTED;
        room.limitTime = block.timestamp + timeLimit;
        emit SetDeck(_roomNumber,"success");
       return true;
    }
    function chosePlayerCompare(uint  _roomNumber,address [] memory chosenArr)public  
    {
        Room storage room = roomStructs[_roomNumber];
        require(room.state == State.STARTED,"invalid state");
        require(msg.sender == room.bidder, "only bidder can chose players to compare");
        require(chosenArr.length <=16," choose maximum 16 people to compare");
        for(uint8 i=0; i<chosenArr.length;i++){
            if(chosenArr[i]==room.players[chosenArr[i]].addr){
                room.chosenPlayers.push(chosenArr[i]);
            }
        }        
        room.state = State.STOPPED;
    }

    function comfirm(uint  _roomNumber,address addr,string []memory decodedKey,string []memory decodedCardsArr )public onlyAdmin 
    {
        Room storage room = roomStructs[_roomNumber];
        
        require(room.state == State.STOPPED,"invalid state");
        require(room.players[addr].playerId != 0,"there is not this player exist in room ");
        // require(room.chosenPlayers != 0,"there is not this player exist in room ");
        require(room.players[addr].decodedCards.length ==0,"alreadly comfirmed this player address");
        require(cardsPerPlayer == decodedKey.length && decodedKey.length == decodedCardsArr.length ,"number of info comfirm array is wrong");
        room.players[addr].decodedCards = decodedCardsArr;
        room.players[addr].keys = decodedKey;
        uint lastdigit = checkLastDigit(room,decodedCardsArr,addr);
        uint idx = room.players[addr].playerId -1;
        room.playersArr[idx].decodedCards =decodedCardsArr;
        room.playersArr[idx].lastDigit =lastdigit;
        room.players[addr].lastDigit =lastdigit;
        //check address comfirm is in chosenPlayers array
        for(uint8 i=0; i<room.chosenPlayers.length;i++){
            if(room.chosenPlayers[i] == addr){
                room.numberComfirm++;
            }
        }
        if(addr == room.bidder){
            room.numberComfirm++;
        }
        if(room.numberComfirm == (room.chosenPlayers.length+1)){
            room.state = State.COMFIRMED ;
            compare(_roomNumber);
        }
    }
    function checkLastDigit(Room storage room,string [] memory decodedCards, address addr)private returns (uint) 
    {
        for(uint8 i=0;i <decodedCards.length;i++){  
            string memory decodedCard = decodedCards[i];  
            if(keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("1S"))||keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("1C"))||keccak256(abi.encodePacked(decodedCard)) ==keccak256(abi.encodePacked("1D"))||keccak256(abi.encodePacked(decodedCard))==keccak256(abi.encodePacked("1H"))){
                room.players[addr].ranks.push(1);
            }
            else if(keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("2S"))||keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("2C"))||keccak256(abi.encodePacked(decodedCard)) ==keccak256(abi.encodePacked("2D"))||keccak256(abi.encodePacked(decodedCard))==keccak256(abi.encodePacked("2H"))){
                room.players[addr].ranks.push(2);
            }
            else if(keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("3S"))||keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("3C"))||keccak256(abi.encodePacked(decodedCard)) ==keccak256(abi.encodePacked("3D"))||keccak256(abi.encodePacked(decodedCard))==keccak256(abi.encodePacked("3H"))){

                room.players[addr].ranks.push(3);
            }
            else if(keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("4S"))||keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("4C"))||keccak256(abi.encodePacked(decodedCard)) ==keccak256(abi.encodePacked("4D"))||keccak256(abi.encodePacked(decodedCard))==keccak256(abi.encodePacked("4H"))){

                room.players[addr].ranks.push(4);
            }
            else if(keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("5S"))||keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("5C"))||keccak256(abi.encodePacked(decodedCard)) ==keccak256(abi.encodePacked("5D"))||keccak256(abi.encodePacked(decodedCard))==keccak256(abi.encodePacked("5H"))){

                room.players[addr].ranks.push(5);
            }
            else if(keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("6S"))||keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("6C"))||keccak256(abi.encodePacked(decodedCard)) ==keccak256(abi.encodePacked("6D"))||keccak256(abi.encodePacked(decodedCard))==keccak256(abi.encodePacked("6H"))){

                room.players[addr].ranks.push(6);
            }
            else if(keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("7S"))||keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("7C"))||keccak256(abi.encodePacked(decodedCard)) ==keccak256(abi.encodePacked("7D"))||keccak256(abi.encodePacked(decodedCard))==keccak256(abi.encodePacked("7H"))){

                room.players[addr].ranks.push(7);
            }
            else if(keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("8S"))||keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("8C"))||keccak256(abi.encodePacked(decodedCard)) ==keccak256(abi.encodePacked("8D"))||keccak256(abi.encodePacked(decodedCard))==keccak256(abi.encodePacked("8H"))){

                room.players[addr].ranks.push(8);
            }
            else if(keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("9S"))||keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("9C"))||keccak256(abi.encodePacked(decodedCard)) ==keccak256(abi.encodePacked("9D"))||keccak256(abi.encodePacked(decodedCard))==keccak256(abi.encodePacked("9H"))){

                room.players[addr].ranks.push(9);
            }
            else{
                room.players[addr].ranks.push(10);
            }
        }
        uint totalRank = calculateTotalRank(room,room.players[addr].ranks,addr);
        uint lastDigit ;
        
        if (totalRank == 30 ) {          
            bool specialcase = checkSpecialCase(decodedCards);
            if(specialcase == true){
                lastDigit = 10 ;  
            }else{
                lastDigit = 0;
            }
            
        }else{
            lastDigit = getLastDigitRank(totalRank);
        }
        return lastDigit;
    }

    function checkSpecialCase(string [] memory decodedCards)private pure  returns (bool)
    {
        for(uint8 i=0;i <decodedCards.length;i++){  
            string memory decodedCard = decodedCards[i];  
            if(keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("10S"))||keccak256(abi.encodePacked(decodedCard)) == keccak256(abi.encodePacked("10C"))||keccak256(abi.encodePacked(decodedCard)) ==keccak256(abi.encodePacked("10D"))||keccak256(abi.encodePacked(decodedCard))==keccak256(abi.encodePacked("10H"))){
                return false;
            }else{
                return true;
            }
        }
    }
    function calculateTotalRank(Room storage room,uint[] memory hand, address addr) private returns (uint)
    {
        require(hand.length==cardsPerPlayer,"each player has 3 cards");
        uint totalRank = 0;
        for (uint i = 0; i < hand.length; i++) {
            totalRank += hand[i];
        }
        room.players[addr].totalRank = totalRank;
        uint idx = room.players[addr].playerId -1;
        room.playersArr[idx].totalRank =totalRank;
        return totalRank;
    }
    function getLastDigitRank(uint _totalRank)private pure returns (uint)
    {
        uint lastdigit;
        if (_totalRank<10){
            lastdigit = _totalRank;
        }else if(_totalRank == 10|| _totalRank == 20 ){
            lastdigit = 0;
        }else if(10<_totalRank&& _totalRank<20){
            lastdigit = _totalRank -10;
        }else{
            lastdigit = _totalRank -20;
        }
        return lastdigit;
    }
    event CompareToBidder (address player, string win);
    function compare(uint  _roomNumber) private 
    {
        Room storage room = roomStructs[_roomNumber];
        require(room.state == State.COMFIRMED,"invalid state");
        uint base = room.players[room.bidder].lastDigit ;
        // uint[] memory temp =new uint[](3);
        // temp[0]=base;
        for (uint i = 0; i < room.chosenPlayers.length; i++) {
            
            address chosenAddr =room.chosenPlayers[i];
            Player storage player = room.players[chosenAddr];
            
            // temp[i+1]=player.lastDigit;
            if(player.addr != room.bidder){
                if(player.lastDigit > base){
                    player.win = true;
                    // room.playersArr[player.playerId].win =true;
                    room.players[player.addr].win = true;
                    registeredFunds[_roomNumber][player.addr] += room.bid;
                    registeredFunds[_roomNumber][room.bidder] -= room.bid;
                    emit CompareToBidder(player.addr,"win");

                }else if(player.lastDigit < base){
                    player.win = false;
                    room.players[player.addr].win = false;
                    registeredFunds[_roomNumber][player.addr] -= room.bid;
                    registeredFunds[_roomNumber][room.bidder] += room.bid;
                    emit CompareToBidder(player.addr,"lose");
                }else{
                    player.win = false;
                    room.players[player.addr].win = false;
                    emit CompareToBidder(player.addr,"draw");
                }             
            }
            AddHistory(_roomNumber,chosenAddr);
        }  
        AddHistory(_roomNumber,room.bidder);

        room.state = State.CLOSED;  
        transferFund(_roomNumber);
        // for (uint i = 0; i < room.attendants.length; i++) {
        //     sortForLive(room.attendants[i]);
        //     uint amount = registeredFunds[room.attendants[i]];
        //     withdrawFund[room.attendants[i]] = amount;

        // }
        // sortForActiveRoom();

        // return temp;   
    }
    function transferFund(uint  _roomNumber)private {
        Room storage room = roomStructs[_roomNumber];
        require(room.state == State.CLOSED,"invalid state");
         for (uint i = 0; i < room.attendants.length; i++) {
            sortForLive(room.attendants[i]);
            uint amount = registeredFunds[_roomNumber][room.attendants[i]];
            registeredFunds[_roomNumber][room.attendants[i]]=0;
            withdrawFunds[room.attendants[i]] += amount;
        }
        sortForActiveRoom();
    }
    function balanceOfRegisteredFunds(address player) public view returns(uint totalRegisterFund)
    {
        for(uint i; i<mTotalUserLive[player]; i++){
            totalRegisterFund += registeredFunds[userLive[player][i]][player];
        }
        return  totalRegisterFund;
    }

    function balanceOfwithdrawFunds(address player) public view returns(uint)
    {
        return  withdrawFunds[player];
    }

    function withdrawFund() external 
    {
        require(withdrawFunds[msg.sender] > 0);

        uint256 funds = withdrawFunds[msg.sender];
        withdrawFunds[msg.sender] = 0;

        (bool success, ) = msg.sender.call{value: funds}("");
        require(success, "MTD transfer failed");
    }
    function getAllPlayersAdd(uint  _roomNumber) public view returns (address[] memory )
    {
        return roomStructs[_roomNumber].attendants;
    }
    function getPlayer(uint  _roomNumber,address _addr)public view returns ( Player memory)
    {
        return roomStructs[_roomNumber].players[_addr];
    }
    function getPlayerBytes(uint  _roomNumber,address _addr)public view returns ( PlayerCopy memory y)
    {
        Player memory player = roomStructs[_roomNumber].players[_addr];
        bytes[] memory cardsBytes = new bytes[](3);
        bytes[] memory keysBytes = new bytes[](3);
        bytes[] memory decodedCardsBytes = new bytes[](3);
        for(uint i;i<3;i++){
            cardsBytes[i] = bytes(player.cards[i]);
            keysBytes[i] = bytes(player.keys[i]);
            decodedCardsBytes[i] = bytes(player.decodedCards[i]);
        }        
        y.playerId = player.playerId;
        y.addr = player.addr;
        y.cards = cardsBytes;
        y.keys = keysBytes;
        y.decodedCards = decodedCardsBytes;
        y.ranks = player.ranks;
        y.totalRank = player.totalRank;
        y.lastDigit = player.lastDigit;
        y.win = player.win;

    }
    function getPlayerCards(uint  _roomNumber,address player) public view returns (string[] memory) 
    {
        uint id = roomStructs[_roomNumber].players[player].playerId -1;
        return roomStructs[_roomNumber].playersArr[id].cards;
    }
    function getACard(uint  _roomNumber,address player,uint card) public view returns (string memory) 
    {
        uint id = roomStructs[_roomNumber].players[player].playerId -1;
        return roomStructs[_roomNumber].playersArr[id].cards[card];
    }
        //Get Room Details
    function getRoom(uint  _roomNumber) public view returns
    (
        address bidder,
        uint roomID,
        uint8 NUM_PLAYERS,
        uint bid ,
        State state,
        address [] memory attendants, //all players join
        uint  registerEndTime,
        address [] memory chosenPlayers,
        uint playersAmount
    )
    {
        Room storage room = roomStructs[_roomNumber];
        //To check Room is created or not!
        uint roomNum = room.roomID;
        if(roomNum != 0){
            return
            (
                room.bidder,
                room.roomID,
                room.NUM_PLAYERS,
                room.bid,
                room.state,
                room.attendants,
                room.registerEndTime,
                room.chosenPlayers,
                room.playersAmount
            );
        }
        else{
            revert("Room is not created!!");
        }
    }
    //quản lý room live
    //nếu status là closed thì gọi hàm xóa room
    function sortForActiveRoom() private {
        for (uint i = 0; i < totalActiveRoom.length; i++) {
            Room storage room = roomStructs[totalActiveRoom[i]];
            if (
                room.state == State.CLOSED                
            ) {
                adjustActive(totalActiveRoom[i]);
            }
        }
    }
    //xóa room trong mảng totalActiveRoom
    function adjustActive(uint _idRoom) private {
        for (uint256 i = 0; i < totalActiveRoom.length; i++) {
            if (totalActiveRoom[i] == _idRoom) {
                totalActiveRoom[i] = totalActiveRoom[
                    totalActiveRoom.length - 1
                ];
                totalActiveRoom.pop();
                // totalRoom--;
            }
        }
    }
    //quản lý user live
    //add vào map những room người chơi tạo ra, tăng  xóa người chơi 
    function AddLive(uint _roomNumber, address _user) private {
        Room storage room = roomStructs[_roomNumber];
        require(
            room.bidder == _user || room.players[_user].addr == _user,
            "Cannot Add Live"
        );
        mTotalUserLive[_user]++;
        userLive[_user].push(_roomNumber);
        sortForLive(_user);
    }
    //loop qua các room có user tham gia nếu có closed state hoặc không phải người chơi thì xóa người chơi khỏi mảng live 
    function sortForLive(address _user) private {
        for (uint256 i = 0; i < userLive[_user].length; i++) {
            Room storage room = roomStructs[userLive[_user][i]];
            if (
                room.state == State.CLOSED
            ) {
                userLive[_user][i] = userLive[_user][
                    userLive[_user].length - 1
                ];
                userLive[_user].pop();
                mTotalUserLive[_user]--;
            }
        }
    }
    //quản lý history (danh sách các room đã hết live được đẩy qua history)
    function AddHistory(uint _roomNumber, address _user) private {
        // Room storage room = roomStructs[_roomNumber];
        // require(
        //     room.bidder == _user || room.players[_user].addr == _user,
        //     "Cannot Add History"
        // );
        mTotalUserHistory[_user]++;
        userHistory[_user].push(_roomNumber);
        sortForHistory(_user);
    }
    //nếu history của 1 user lớn hơn 10 room thì gọi hàm xóa id room cũ nhất trong mảng history
    function sortForHistory(address _user) private {
        if (userHistory[_user].length > 10) {
            adjustHash(0, _user);
            mTotalUserHistory[_user]--;
        }
    }
    //xóa id room đầu tiên trong mảng history
    function adjustHash(uint256 index, address _user) private {
        require(index < userHistory[_user].length, "Invalid index");
        for (uint256 i = index; i < userHistory[_user].length - 1; i++) {
            userHistory[_user][i] = userHistory[_user][i + 1];
        }
        userHistory[_user].pop();
    }

    function GetHistory(address _user) external view returns (RoomCopy[] memory arrayRoom) 
    {
        arrayRoom = new RoomCopy[](userHistory[_user].length);
        for (uint i = 0; i < arrayRoom.length; i++) {
            arrayRoom[i] = convert(userHistory[_user][i]);
        }
    }
    function TotalActiveRoomLength()public view returns(uint){
        return totalActiveRoom.length;
    }
    function getTotalActiveRoom(uint index)public view returns(uint){
        return totalActiveRoom[index];
    }
    function getUserLiveLength(address addr)public view returns(uint){
        return userLive[addr].length;
    }
    function getUserLiveItem(address addr,uint index)public view returns(uint){
        return userLive[addr][index];
    }
    function convert(uint x) public view returns (RoomCopy memory y) {
        Room storage room = roomStructs[x];
        bytes[] memory deck = new bytes[](room.deck.length);
        for (uint i;i<room.deck.length;i++){
            deck[i] = bytes(room.deck[i]);
        }       
        y.roomID=room.roomID; 
        y.NUM_PLAYERS=room.NUM_PLAYERS ;
        y.deck=deck;
        y.bid=room.bid ;
        y.state=room.state;
        y.numberComfirm=room.numberComfirm;
        y.bidder=room.bidder;
        y.attendants=room.attendants; //all players join
        y.playersArr=room.playersArr;
        y.registerEndTime=room.registerEndTime;
        y.chosenPlayers=room.chosenPlayers;
        return y;
    }
}
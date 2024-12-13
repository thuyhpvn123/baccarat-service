Luong bài cào (Baccarat):
1.CreateRoom: Người chơi tạo room sẽ trở thành nhà cái(bidder) gọi hàm CreateRoom với registerTime là thời gian tối đa của room để người chơi khác được tham gia. 
-bidAmount là tiền cược của mỗi người chơi thì join room. 
-playerAmount: số người chơi được tham gia phòng(trừ người tạo phòng)
-Người tạo room phải chuyển số tiền = bidAmount* (số người được tham gia-1). Số phòng được tạo bắt đầu từ 1, tự động tăng. Một người có thể tạo nhiều phòng. Tối đa 17 người tham gia.

Sau khi tạo phòng thì có thể lấy thông tin phòng bằng hàm getRoom(số phòng)
2.register: Người chơi tham gia room nào thì gọi hàm register và chọn số phòng

3.closeRoom: Sau khi hết thời gian đang kí mà không ai tham gia thì bidder được gọi hàm closeRoom để đóng phòng .

4.deal-cards: FE gửi lệnh xuống server để tạo 52 key cho 52 lá bài, tạo bộ bài, xào bài, mã hóa các lá bài theo 52 key đã tạo rồi rồi admin sẽ gọi hàm setDeck trên contract để chia bài: 
    command: 'deal-cards', value:{roomNumber: ""}
    
5.getPlayerCards: Người chơi gọi hàm getPlayerCards  lên contract để lấy 3 quân bài mã hóa đã được chia
6.get-key-for-player: Người chơi gọi lên server  "get-key-for-player" để lấy key giải mã 3 lá bài. Truyền tới server:
"encrypted-cards":3 lá bài đã encrypted gọi bằng hàm getPlayerCards,

command: 'get-key-for-player', 
value:{
	"hash":hash,
      "sign":sign,
      "pubKey":pubKey,
      "roomNumber":roomid
}
tất cả là string
7.decrypt-cards: Sau khi có key thì người chơi gửi lên server để giải mã lấy ra lá bài :
command: 'decrypt-cards', 
value:{
	encrypted-cards: ['AW8=', '8NM=', 'XVE=']
	encrypted-keys:["d1d91e6a8e34f5bdb3ac9545b8a64153","2d0e0ee7123200786eb8d2ac4eee9a99","9475ebfe8f03c6658de735453aaeab6c"]

}
8.chosePlayerCompare: Sau khi chia bài thì bidder sẽ được chọn những người chơi trong số những người đã join phòng để so bài 
9.comfirm: sau khi bidder chọn xong người chơi thì FE gửi lên server gọi hàm comfirm để admin gửi kết quả của tất cả các người chơi được chọn bao gồm cả bidder lên contract.
command: 'comfirm', 
value:{
	roomid:""
}
khi đã comfirm đủ để so sánh kết quả với các người chơi. 
11.withdrawFund: người chơi gọi hàm này để rút tiền (chỉ rút được lượng tiền của những phòng đã closed)
12.playerCloseRoom: người chơi trừ bidder có thể gọi hàm này để đóng phòng nếu sau khi chia bài 1 khoảng thời gian limit(đang set mặc định là 5 phút) mà bidder không chọn người so bài
13.balanceOfwithdrawFunds: kiểm tra số dư tiền người chơi có thể rút từ contract
14.setTimeLimit: admin update limit time là thời gian tối đa sau khi chơi bài bidder có thể chọn người so bài. sau thời gian này có thể gọi hàm playerCloseRoom để đóng phòng.
15.balanceOfRegisteredFunds: kiểm tra tổng tiền người chơi đã đăng kí 

16.registeredFunds(số phòng, địa chỉ người chơi): lấy ra số tiền đã cược của người chơi mỗi phòng.

*Note: 
- room active la room chua room chua o trang thai CLOSED. 
- userLive la chi user dang trong nhung room active

*SO KET QUA NHU SAU: 
- Nhung nguoi choi duoc bidder chon thi se so sanh so cuoi voi bidder, neu > thi la thang, < la thua, = thi ca 2 deu thua
- So cuoi duoc tinh: neu la bai la 1 thi tong diem +1, la 2 thi tong diem +2. Tuong tu den het so 9. Cac la 10,J,Q,K se duoc +10
+ neu tong  = 30 thi se kiem tra xem trong 3 la neu co 1 la 10 thi so cuoi = 0 , neu khong thi so cuoi = 10 
+ neu tong <10 thi so cuoi = tong
+ neu tong  = 10 hoac 20 thi so cuoi =0
+ neu 10<tong<20 thi so cuoi = tong - 10
+ con lai(>20 va <30) so cuoi = tong - 20

Cac ham khac trong contract trong ThreeCardsMain:
- totalRoom: tong so room da duoc tao ra
- totalActiveRoom: tong so room active
- convert: lay thong tin 1 phong theo id phong
- getUserLiveRoomIds: lay ra id cac phong userLive tham gia theo dia chi user
- getUserLiveLength: lay ra tong so userLive
- TotalActiveRoomLength: lay ra so luong phong active
- GetHistory: lay ra mang cac phong 1 user da tham gia theo dia chi user
- getRoom: lay ra thong tin 1 phong , khong bao gom bo bai va chi tiet nhung nguoi choi
- getACard: lay ra 1 quan bai duoc chia cua 1 nguoi theo so phong, dia chi nguoi choi, va thu tu quan bai trong mang(tu 0 den 2)
- getPlayerCards: lay ra bo 3 quan bai duoc chia cua nguoi choi theo so phong va dia chi nguoi choi
- getPlayersBytesAroom: lay thong tin chi tiet tat ca nguoi choi cua 1 phong theo so phong
- getAPlayerBytes: lay ra thong tin chi tiet cua 1 nguoi choi theo so phong va dia chi nguoi choi dang bytes(de FE co the parse ra string)
- getPlayer: lay ra thong tin chi tiet cua 1 nguoi choi theo so phong va dia chi nguoi choi 
- getAllPlayersAdd: lay ra dia chi cua tat ca nguoi choi cua 1 phong
- setTimeLimit:set thoi gian bidder duoc chon nguoi choi de so bai sau khi chia bai, sau thoi gian nay bidder khong duoc chon nua
- setNewAdmin: set lai dia chi admin(nguoi chia bai)
- getChosenPlayerAddr: lay ra danh sach dia chi nguoi choi duoc chon de so bai trong 1 phong theo so phong

Cac ham trong contract ThreeCardsMainGet:
- GetActiveRoom: lay danh sach chi tiet cac room dang active
- GetLiveRoom: lay danh sach chi tiet cac room dang active cua 1 user

Cac truong hop room chuyen sang CLOSED:
1.sau khi admin goi comfirm de so bai, sau khi so lan goi comfirm = so luong nguoi bidder chon de so bai thi room chuyen qua la CLOSED
2.sau khi qua register time ma khong nguoi choi nao tham gia tru bidder thi bidder goi closeRoom de dong phong, lay lai tien da coc
3.khi bidder khong chon ai de so bai trong thoi han quy dinh(dang set la 1phut) sau khi da chia bai thi nguoi choi goi playerCloseRoom de dong phong, nhan lai tien coc



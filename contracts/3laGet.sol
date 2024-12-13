//SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.4;
import "./3laMain.sol";
contract ThreeCardsMainGet {
    ThreeCardsMain main;
        uint8 private limit = 5;

    constructor(address _main) payable {
        main=ThreeCardsMain(_main);
    } 
    function SetLimit(uint8 _number) external  {
        limit = _number;
    }

    function GetActiveRoom(
        uint8 _page
    )
        external
        view
        returns (bool isMore, ThreeCardsMain.RoomCopy[] memory arrayRoom)
    {
        uint length = main.TotalActiveRoomLength();
        if (_page * limit > length + limit) {
            return (false, arrayRoom);
        } else {
            if (_page * limit <= length) {
                isMore = true;
                arrayRoom = new ThreeCardsMain.RoomCopy[](limit);
                for (uint i = 0; i < arrayRoom.length; i++) {
                    arrayRoom[i] = main.convert(
                        // totalActiveRoom[_page * limit - limit + i]
                        main.getTotalActiveRoom(_page * limit - limit + i)
                    );
                }
                return (isMore, arrayRoom);
            } else {
                isMore = false;
                arrayRoom = new ThreeCardsMain.RoomCopy[](
                    limit - (_page * limit - length)
                );
                for (uint i = 0; i < arrayRoom.length; i++) {
                    arrayRoom[i] = main.convert(
                        // totalActiveRoom[_page * limit - limit + i]
                        main.getTotalActiveRoom(_page * limit - limit + i)
                    );
                }
                return (isMore, arrayRoom);
            }
        }
    }
        function GetLiveRoom(
        uint8 _page,
        address _user
    )
        external
        view
        returns (bool isMore, ThreeCardsMain.RoomCopy[] memory arrayRoom)
    {
        uint length = main.getUserLiveLength(_user);
        if (_page * limit > length + limit) {
            return (false, arrayRoom);
        } else {
            if (_page * limit <= length) {
                isMore = true;
                arrayRoom = new ThreeCardsMain.RoomCopy[](limit);
                for (uint i = 0; i < arrayRoom.length; i++) {
                    arrayRoom[i] = main.convert(
                        // userLive[_user][_page * limit - limit + i]
                        main.getUserLiveItem(_user,_page * limit - limit + i)
                    );
                }
                return (isMore, arrayRoom);
            } else {
                isMore = false;
                arrayRoom = new ThreeCardsMain.RoomCopy[](
                    limit - (_page * limit - length)
                );
                for (uint i = 0; i < arrayRoom.length; i++) {
                    arrayRoom[i] = main.convert(
                        // userLive[_user][_page * limit - limit + i]
                        main.getUserLiveItem(_user,_page * limit - limit + i)
                    );
                }
                // RoomCopy memory room =convert(userLive[_user][_page * limit - limit + 0]);
                // arrayRoom[0] = room;
                return (isMore, arrayRoom);
            }
        }
    }


}
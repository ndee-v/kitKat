package utils

// HelpMessage ...
const HelpMessage = `[PARAMS USAGE]:
tr=   time for registration // default: 30 seconds // ex: tr=30	
ti=   time for idle timeout // default: 180 seconds // ex: ti=180
r=    count of rooms for chats // default: 4 rooms // ex: r=4
c=    count of maximum connections // default: 10 connections // ex: c=10
g=    gui terminal on/off // default: off // ex: g=on or g=true`

// IncorrectName ...
const IncorrectName = `[INCORRECT NAME]
[USAGE]: any word character
from a-z, A-Z, 0-9, _
min length 1 character
max length 10 characters`

// HelpOptions ...
const HelpOptions = `[HELP OPTIONS]:
--name  : change name 
		example: --name John_Doe		
--online  : list of online users
--room : change chat room
		example: --room 1`

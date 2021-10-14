// Package handler contains an HTTP Cloud Function to handle update from Telegram whenever a users interacts with the
// bot.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Define a few constants and variable to handle different commands
const (
	startCommand string = "/start"
	bksCommand   string = "/bks"
	infoCommand  string = "/info"
	moneyCommand string = "/money"
	helpCommand  string = "/help"
)

const (
	botTag         string = "@xCheckInformationBot"
	welcomeMessage string = `Xin chào!
Số dư tài khoản: %v trong tài khoản.

Dùng /money để nạp tiền, dùng /help để xem hướng dẫn!
Group support (không chính thức): https://t.me/biensoxe`

	helpMessage string = `/bsk - Tra thông tin biển số xe - Giá: $2.00
/info - Xem thông tin tài khoản
/money - Nạp tiền vào tài khoản
/help - Hướng dẫn sử dụng

Ghi chú:
Group support:@sodienthoai https://t.me/biensoxe`

	moneyMessage string = `ĐỊA CHỈ TIỀN ẢO NẠP TIỀN

Các địa chỉ dưới đây DÀNH RIÊNG cho tài khoản có ID: %v, chấp nhận TẤT CẢ các mức nạp.
Giá trị = số coin thực nhận * tỷ giá tại thời điểm tạo địa chỉ nạp

Bitcoin(BTC = %v):
%v

Bitcoin Cash(BCH = %v):
%v

Ethereum(ETH= %v):
%v

Litecoin(LTC= %v):
%v

DAI(DAI= %v):
%v

DOGECOIN(DOGE= %v):
%v

USDC(=USDC %v - mạng lưới ERC20):
%v

Sử dụng: /money<Mã tiền ảo> để lấy địa chỉ + ảnh QR riêng cho từng Coin.
Ví dụ:
/money bitcoin
/money BTC
/money ETH
/money Litecoin`

	bksMessage = `Biển số xe: %v
Địa chỉ: Times City, 458 Minh Khai, Hai Bà Trưng, Hà Nội
Số điện thoại: 0912665660
Số khung: ABC123456XYZ
Số máy: XYZ9876543210ABC
Nhãn hiệu: Honda
Trạng thái xe: Đang lưu hành`

	infoMessage = `Mã ID(Username): %v
Fullname: %v
Số dư tài khoản: $2.00`
)

// Pass token and sensible APIs through environment variables
const (
	telegramApiBaseUrl     string = "https://api.telegram.org/bot"
	telegramApiSendMessage string = "/sendMessage"
	telegramTokenEnv       string = "2013265111:AAEADc9CE21y23vmu6tcKQr1RMm6uTINd5Q"

	demoApiEnv string = "https://dummy.restapiexample.com/api/v1/employees"
)

var lenBotTag = len(botTag)
var lenStartCommand = len(startCommand)
var lenHelpCommand = len(helpCommand)
var lenMoneyCommand = len(moneyMessage)
var lenInfoCommand = len(infoCommand)
var lenBKSCommand = len(bksCommand)

var telegramApi = telegramApiBaseUrl + telegramTokenEnv + telegramApiSendMessage
var rapLyricsApi = demoApiEnv

// Update is a Telegram object that we receive every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Implements the fmt.String interface to get the representation of an Update as a string.
func (u Update) String() string {
	return fmt.Sprintf("(update id: %d, message: %s)", u.UpdateId, u.Message)
}

// Message is a Telegram object that can be found in an update.
// Note that not all Update contains a Message. Update for an Inline Query doesn't.
type Message struct {
	Text     string   `json:"text"`
	Chat     Chat     `json:"chat"`
	Audio    Audio    `json:"audio"`
	Voice    Voice    `json:"voice"`
	Document Document `json:"document"`
}

// Implements the fmt.String interface to get the representation of a Message as a string.
func (m Message) String() string {
	return fmt.Sprintf("(text: %s, chat: %s, audio %s)", m.Text, m.Chat, m.Audio)
}

// Audio message has extra attributes
type Audio struct {
	FileId   string `json:"file_id"`
	Duration int    `json:"duration"`
}

// Implements the fmt.String interface to get the representation of an Audio as a string.
func (a Audio) String() string {
	return fmt.Sprintf("(file id: %s, duration: %d)", a.FileId, a.Duration)
}

// Voice Message can be summarized with similar attribute as an Audio message for our use case.
type Voice Audio

// Document Message refer to a file sent.
type Document struct {
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
}

// Implements the fmt.String interface to get the representation of an Document as a string.
func (d Document) String() string {
	return fmt.Sprintf("(file id: %s, file name: %s)", d.FileId, d.FileName)
}

// A Chat indicates the conversation to which the Message belongs.
type Chat struct {
	Id int `json:"id"`
}

// Implements the fmt.String interface to get the representation of a Chat as a string.
func (c Chat) String() string {
	return fmt.Sprintf("(id: %d)", c.Id)
}

// A Lyric is generated by the RapLyrics service.
type Lyric struct {
	Punch string `json:"output"`
}

//type WalletAccount struct {
//	MakerCommission  int    `json:"makerCommission"`
//	TakerCommission  int    `json:"takerCommission"`
//	BuyerCommission  int    `json:"buyerCommission"`
//	SellerCommission int    `json:"sellerCommission"`
//	CanTrade         bool   `json:"canTrade"`
//	CanWithdraw      bool   `json:"canWithdraw"`
//	CanDeposit       bool   `json:"canDeposit"`
//	UpdateTime       int64  `json:"updateTime"`
//	AccountType      string `json:"accountType"`
//	Balances         []struct {
//		Asset  string `json:"asset"`
//		Free   string `json:"free"`
//		Locked string `json:"locked"`
//	} `json:"balances"`
//	Permissions []string `json:"permissions"`
//}

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {

	// Parse incoming request
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	// Sanitize input
	var sanitizedSeed = sanitize(update.Message.Text)
	//var telegramResponseBody string
	//var errTelegram error

	var message string

	switch sanitizedSeed {
	case "":
		message = welcomeMessage
	case helpCommand:
		message = helpMessage
	case moneyCommand:
		message = moneyMessage
	case infoCommand:
		message = infoMessage
	default:
		message = bksMessage
		//case bksCommand:
		//	message = bksMessage
	}

	telegramResponseBody, errTelegram := sendTextToTelegramChat(update.Message.Chat.Id, message)
	if errTelegram != nil {
		log.Printf("got error %s from telegram, response body is %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("startline %s successfully distributed to chat id %d", welcomeMessage, update.Message.Chat.Id)
	}

	// Start command
	//if sanitizedSeed == "" {
	//	// Init welcome message
	//	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, welcomeMessage)
	//	if errTelegram != nil {
	//		log.Printf("got error %s from telegram, response body is %s", errTelegram.Error(), telegramResponseBody)
	//	} else {
	//		log.Printf("startline %s successfully distributed to chat id %d", welcomeMessage, update.Message.Chat.Id)
	//	}
	//} else {
	//
	//	// Call RapLyrics to get a punchline
	//	var lyric, errRapLyrics = getPunchline(sanitizedSeed)
	//	if errRapLyrics != nil {
	//		log.Printf("got error when calling RapLyrics API %s", errRapLyrics.Error())
	//		return
	//	}
	//
	//	// Send the punchline back to Telegram
	//	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, lyric)
	//	if errTelegram != nil {
	//		log.Printf("got error %s from telegram, response body is %s", errTelegram.Error(), telegramResponseBody)
	//	} else {
	//		log.Printf("punchline %s successfully distributed to chat id %d", lyric, update.Message.Chat.Id)
	//	}
	//}
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	if update.UpdateId == 0 {
		log.Printf("invalid update id, got update id = 0")
		return nil, errors.New("invalid update id of 0 indicates failure to parse incoming update")
	}
	return &update, nil
}

// sanitize remove clutter like /start /punch or the bot name from the string s passed as input
func sanitize(s string) (msg string) {
	if len(s) >= lenStartCommand {
		if s[:lenStartCommand] == startCommand {
			msg = s[lenStartCommand:]
			return msg
		}
	}

	if len(s) >= lenHelpCommand {
		if s[:lenHelpCommand] == helpCommand {
			//msg = s[lenHelpCommand:]
			return s
		}
	}

	if len(s) >= lenMoneyCommand {
		if s[:lenMoneyCommand] == moneyCommand {
			msg = s[lenMoneyCommand:]
			return msg
		}
	}

	if len(s) >= lenInfoCommand {
		if s[:lenInfoCommand] == infoCommand {
			//msg = s[lenInfoCommand:]
			return s
		}
	}

	if len(s) >= lenBKSCommand {
		if s[:lenBKSCommand] == bksCommand {
			msg = s[lenBKSCommand:]
			return msg
		}
	}

	if len(s) >= lenBotTag {
		if s[:lenBotTag] == botTag {
			msg = s[lenBotTag:]
			return msg
		}
	}

	return s
}

// getPunchline calls the RapLyrics API to get a punchline back.
func getPunchline(seed string) (string, error) {
	rapLyricsResp, err := http.PostForm(
		rapLyricsApi,
		url.Values{"input": {seed}})
	if err != nil {
		log.Printf("error while calling raplyrics %s", err.Error())
		return "", err
	}
	var punchline Lyric
	if err := json.NewDecoder(rapLyricsResp.Body).Decode(&punchline); err != nil {
		log.Printf("could not decode incoming punchline %s", err.Error())
		return "", err
	}
	defer rapLyricsResp.Body.Close()
	return punchline.Punch, nil
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatId int, text string) (string, error) {

	log.Printf("Sending %s to chat_id: %d", text, chatId)

	//var telegramApi = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

func main() {
	fmt.Println("start bot API")
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", 8080), http.HandlerFunc(HandleTelegramWebHook))
	if err != nil {
		return
	}
}

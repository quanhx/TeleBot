// Package handler contains an HTTP Cloud Function to handle update from Telegram whenever a users interacts with the
// bot.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"xcheck.info/telebot/pkg/conf"
	"xcheck.info/telebot/pkg/database"
	"xcheck.info/telebot/pkg/dtos"
	"xcheck.info/telebot/pkg/routers"
	"xcheck.info/telebot/pkg/services"
)

func init() {
	err := conf.InitConfig()
	if err != nil {
		log.Fatalln("Error when init config, detail: ", err)
		panic(err)
	}

	database.InitDb()
	services.InitServices()

}

// FilePath is a path to a local file.
type FilePath string

// Define assets accept for payment
const (
	DOGE = "doge"
	XRP  = "xrp"
)

// payment status
const (
	Waiting       = "waiting"
	Confirming    = "confirming"
	Confirmed     = "confirmed"
	Sending       = "sending"
	PartiallyPaid = "partially_paid"
	Finished      = "finished"
	Failed        = "failed"
	Refunded      = "refunded"
	Expired       = "expired"
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
	welcomeMessage string = `Xin chào %v!
Số dư tài khoản: $%v trong tài khoản.

Dùng /money để nạp tiền, dùng /help để xem hướng dẫn!
Group support (không chính thức): https://t.me/biensoxe`

	helpMessage string = `/bks - Tra thông tin biển số xe - Giá: $2.00
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
/money Bitcoin
/money BTC
/money ETH
/money DOGE`

	motoLicenseNoMessage = `Biển số xe: %v
Số khung: %v
Số máy: %v
Ngày đăng ký: %v
Ngày đăng ký lần đầu: %v
Trạng thái: %v
Tên màu biển: %v
Nhãn hiệu: %v`

	carLicenseNoMessage = `Biển số xe: %v
Tên chủ xe: %v
Địa chỉ: %v
Số điện thoại: %v
Số máy: %v
Số khung: %v
Nhãn hiệu: %v
Số loại: %v
Tên màu biển: %v
Mẫu biển: %v
Năm sản xuất: %v
Ngày đăng ký lần đầu: %v
Ngày đăng ký: %v`

	infoMessage = `Mã ID(Username): %v
Fullname: %v
Số dư tài khoản: $%v`

	depositMessage = `Số dư của bạn không đủ. Vui lòng nạp vào tài khoản. Dùng /money để nạp tiền.`

	errMessage = `Thông tin bạn tìm không thấy. Vui lòng quay lại sau. Xin cảm ơn.`
)

// Pass token and sensible APIs through environment variables
const (
	telegramApiBaseUrl     string = "https://api.telegram.org/bot"
	telegramApiSendMessage string = "/sendMessage"
	//telegramTokenEnv       string = "2013265111:AAEADc9CE21y23vmu6tcKQr1RMm6uTINd5Q"
	telegramTokenEnv string = "2013265111:AAEDEd-jySxjLcJTtdpuOsMRe53rzoojEqY"

	//motoAPI string = "https://g32932006cda407-moto.adb.ap-tokyo-1.oraclecloudapps.com/ords/data/moto/search/%s"
	motoAPI string = "http://45.32.113.203/api/moto/search/%s"
	carAPI  string = "http://45.32.113.203/api/car/search/%s"
	//carAPI  string = "https://g32932006cda407-bks.adb.ap-tokyo-1.oraclecloudapps.com/ords/oto/search/bks/%s"
)

var lenBotTag = len(botTag)
var lenStartCommand = len(startCommand)
var lenHelpCommand = len(helpCommand)
var lenMoneyCommand = len(moneyCommand)
var lenInfoCommand = len(infoCommand)
var lenBKSCommand = len(bksCommand)

var telegramApi = telegramApiBaseUrl + telegramTokenEnv + telegramApiSendMessage

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
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
}

// Implements the fmt.String interface to get the representation of a Chat as a string.
func (c Chat) String() string {
	return fmt.Sprintf("(id: %d)", c.Id)
}

type Car struct {
	TenChuXe     string `json:"ten_chu_xe"`
	DiaChi       string `json:"dia_chi"`
	SoDienThoai  string `json:"so_dien_thoai"`
	SoMay        string `json:"so_may"`
	SoKhung      string `json:"so_khung"`
	NhanHieu     string `json:"nhan_hieu"`
	SoLoai       string `json:"so_loai"`
	TenMau       string `json:"ten_mau"`
	MauBien      string `json:"mau_bien"`
	NamSx        int    `json:"nam_sx"`
	NgayDkLandau string `json:"ngay_dk_landau"`
	NgayDangky   string `json:"ngay_dangky"`
}

type Moto struct {
	SoKhung      string `json:"so_khung"`
	SoMay        string `json:"so_may"`
	NgayDangky   string `json:"ngay_dangky"`
	NgayDkLandau string `json:"ngay_dk_landau"`
	TrangThaiHs  string `json:"trang_thai_hs"`
	TenMauBien   string `json:"ten_mau_bien"`
	NhanHieu     string `json:"nhan_hieu"`
}

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	// Parse incoming request
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}
	// Init env
	db := database.InitDb()
	defer db.Close()

	userAPI := routers.InitUserAPI(db)
	transactionAPI := routers.InitTransactionAPI(db)

	// Get userID from ChatID and register TeleBot's user if not exist
	profile := userAPI.UserService.FindByUserID(uint(update.Message.Chat.Id))
	if profile == nil || (profile != nil && profile.ID == 0) {
		var user = dtos.UserRequest{
			ChatID:     update.Message.Chat.Id,
			FirstName:  update.Message.Chat.FirstName,
			LastName:   update.Message.Chat.LastName,
			DeviceName: "",
			Phone:      "",
			Status:     "ACTIVE",
			Token:      "",
			UserID:     int64(update.Message.Chat.Id),
			UserName:   update.Message.Chat.UserName,
		}

		err := userAPI.UserService.CreateUser(user)
		if err != nil {
			log.Fatalln("Can not create user")
		}
	}

	// Get Balance
	//var needDeposit bool
	balance, err := transactionAPI.TransactionService.GetBalance(uint(update.Message.Chat.Id))
	if err != nil {
		//needDeposit = true
		log.Fatalln("balance is empty")
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
		message = fmt.Sprintf(infoMessage, update.Message.Chat.UserName, update.Message.Chat.UserName, balance)
		//default:
		//	message = motoLicenseNoMessage
		//case bksCommand:
		//	message = motoLicenseNoMessage
	}

	//	Call get information
	sanitizedSeed = strings.Trim(sanitizedSeed, " ")
	if len(message) == 0 {
		if strings.Contains(sanitizedSeed, "moto") {
			sanitizedSeed = strings.Replace(sanitizedSeed, "moto", "", -1)
			if balance > 0 {
				message, err = getMotoInformation(sanitizedSeed)
				if err != nil {
					log.Printf("got error when calling Moto API %s", err.Error())
					message = errMessage
				}
			} else {
				message = depositMessage
			}
		} else if strings.Contains(sanitizedSeed, "oto") {
			sanitizedSeed = strings.Replace(sanitizedSeed, "oto", "", -1)
			if balance > 0 {
				message, err = getCarInformation(sanitizedSeed)
				if err != nil {
					log.Printf("got error when calling Car API %s", err.Error())
					message = errMessage
				}
				err := paymentForCheckService(int64(update.Message.Chat.Id))
				if err != nil {
					log.Printf("got error payment for check service %s", err.Error())
				}
			} else {
				message = depositMessage
			}
		} else if strings.Contains(sanitizedSeed, "doge") {
			//sanitizedSeed = strings.Replace(sanitizedSeed, "oto", "", -1)
			//if needDeposit {
			deposit, err := deposit(update.Message.Chat.Id, "doge", 100)
			if err != nil || deposit == nil {
				log.Fatal(err)
			}

			// Check payment status
			var paymentOrder *services.PaymentStatusResponse
			for i := 0; i <= 3; i++ {
				// Sleep 3 minute and get status payment
				time.Sleep(300 * time.Second)

				paymentStatus := services.PaymentStatus(deposit.PaymentId)
				if paymentStatus == nil {
					log.Printf("Payment not found")
				}

				if paymentStatus.PaymentStatus == Finished || paymentStatus.PaymentStatus == Confirmed ||
					paymentStatus.PaymentStatus == Confirming || paymentStatus.PaymentStatus == PartiallyPaid {
					paymentOrder = paymentStatus
					break
				} else {
					continue
				}
			}

			var transactionRequest = dtos.TransactionRequest{
				UserID:           int64(update.Message.Chat.Id),
				Status:           paymentOrder.PaymentStatus,
				OrderDescription: fmt.Sprintf("Deposit from Now Payment with %s coint", paymentOrder.PayCurrency),
				OrderID:          paymentOrder.OrderId,
				PayCurrency:      paymentOrder.PayCurrency,
				PayAmount:        paymentOrder.ActuallyPaid,
				PriceCurrency:    paymentOrder.PriceCurrency,
				PriceAmount:      paymentOrder.PriceAmount,
				PaymentAddress:   paymentOrder.PayAddress,
			}
			transaction := transactionAPI.TransactionService.CreateTransaction(transactionRequest)
			if err != nil {
				log.Println("can not create transaction", transaction)
			}
			fullName := fmt.Sprintf("%s %s", update.Message.Chat.FirstName, update.Message.Chat.LastName)
			message = fmt.Sprintf(infoMessage, update.Message.Chat.Id, fullName, balance)

			//}
		}
	}

	if message == welcomeMessage {
		// Create user from telegram information
		message = fmt.Sprintf(welcomeMessage, update.Message.Chat.UserName, balance)
	}

	telegramResponseBody, errTelegram := sendTextToTelegramChat(update.Message.Chat.Id, message)
	if errTelegram != nil {
		log.Printf("got error %s from telegram, response body is %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("startline %s successfully distributed to chat id %d", welcomeMessage, update.Message.Chat.Id)
	}
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
			//msg = s[lenMoneyCommand:]
			return s
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

// getMotoInformation calls the Moto API to get a moto back.
func getMotoInformation(licenseMotoNo string) (string, error) {
	licenseMotoNo = strings.Trim(licenseMotoNo, " ")
	searchURL := fmt.Sprintf(motoAPI, licenseMotoNo)
	resp, err := http.Get(searchURL)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var info Moto
	//Convert the body to type string
	if err := json.Unmarshal(body, &info); err != nil {
		return "", err
	}
	return fmt.Sprintf(motoLicenseNoMessage, licenseMotoNo, info.SoKhung, info.SoMay, info.NgayDangky, info.NgayDkLandau, info.TrangThaiHs, info.TenMauBien, info.NhanHieu), nil
}

// getCarInformation calls the Car API to get a moto back.
func getCarInformation(licenseCarNo string) (string, error) {
	licenseCarNo = strings.Trim(licenseCarNo, " ")
	searchURL := fmt.Sprintf(carAPI, licenseCarNo)
	resp, err := http.Get(searchURL)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var info Car
	if err := json.Unmarshal(body, &info); err != nil {
		return "", err
	}
	return fmt.Sprintf(carLicenseNoMessage, licenseCarNo, info.TenChuXe, info.DiaChi, info.SoDienThoai, info.SoMay,
		info.SoKhung, info.NhanHieu, info.SoLoai, info.TenMau, info.MauBien, info.NamSx, info.NgayDkLandau, info.NgayDangky), nil
}

//
func paymentForCheckService(chatID int64) error {
	// Init env
	db := database.InitDb()
	defer db.Close()

	transactionAPI := routers.InitTransactionAPI(db)

	var transactionRequest = dtos.TransactionRequest{
		UserID:           chatID,
		Status:           "finished",
		OrderDescription: "Payment for check information service x 1",
		OrderID:          strconv.FormatInt(time.Now().Unix(), 10),
		PayCurrency:      "usd",
		PayAmount:        -2,
		PriceCurrency:    "usd",
		PriceAmount:      -2,
		PaymentAddress:   "",
		PaymentType:      1,
	}
	//err := services.CreatePayment(transactionRequest)
	err := transactionAPI.TransactionService.CreateTransaction(transactionRequest)
	if err != nil {
		log.Println("can not create transaction", transactionRequest)
	}
	return nil
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

// sendPhotoToTelegramChat sends a photo to the Telegram chat identified by its chat Id
func sendPhotoToTelegramChat(chatId int, photo string) (string, error) {

	log.Printf("Sending %s to chat_id: %d", photo, chatId)

	//var telegramApi = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"photo":   {photo},
			"caption": {"QR Code asset cryptocurrency for deposit"},
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

var asset = map[string]string{
	"XRP":  "XXRP",
	"NANO": "NANO",
	"USDT": "USDT",
	"XLM":  "XLM",
	"DOGE": "XXDG",
}

var assetMethod = map[string]string{
	"XXRP": "Ripple XRP",
	"NANO": "NANO",
	"USDT": "Tether USD (TRC20)",
	"XLM":  "Stellar XLM",
	"XXDG": "Dogecoin",
}

func deposit(chatID int, asset string, amount int) (*services.PaymentResponse, error) {
	// Get address for user
	address := services.Payment(amount, asset)
	if address == nil {
		log.Fatalln("Can not get address")
	}

	// Gen QR code from address
	qr := services.GenQR(address.PayAddress, address.PayCurrency)
	fmt.Println(qr)
	msg := botapi.NewPhoto(int64(chatID), botapi.FilePath(qr))
	msg.Caption = address.PayAddress
	bot, err := botapi.NewBotAPI(telegramTokenEnv)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	// Send asset address as QR Code
	bot.Send(msg)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return address, nil
}

func main() {
	fmt.Println("start bot API")

	db := database.InitDb()
	defer db.Close()

	userAPI := routers.InitUserAPI(db)

	r := gin.Default()

	r.GET("/users/:id", userAPI.FindByID)
	r.POST("/users", userAPI.Create)

	r.GET("/transaction/:id", userAPI.FindByID)
	r.POST("/transaction", userAPI.Create)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", 8888), http.HandlerFunc(HandleTelegramWebHook))
	if err != nil {
		return
	}
}

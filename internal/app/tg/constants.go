package tg

import tele "gopkg.in/telebot.v3"

// media assets paths
const (
	assetsBasePath = "./assets/"
	logoAssetPath  = assetsBasePath + "logo.jpg"
)

// media assets as tele photos
var (
	logoImage tele.Inputtable = &tele.Photo{File: tele.FromDisk(logoAssetPath), Caption: startMessageTemplate}
)

// buttons
var (
	btnManualBuy = tele.InlineButton{
		Text:   "Join Hanzo Academy 😎",
		Unique: "manual_buy",
	}
	btnSupport = tele.InlineButton{
		Text: "Support 🤓",
		URL:  "https://t.me/defihanzo",
	}
	btnMainMenu = tele.InlineButton{
		Text:   "Main Menu 🉐",
		Unique: "menu",
	}
	btnAcceptPayment = tele.InlineButton{
		Text:   "✅",
		Unique: "accept_payment",
	}
)

// inline selectors for messages
var (
	startSelector = &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{btnManualBuy},
			{btnSupport},
		},
	}
	startPaidSelector = &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{btnSupport},
		},
	}
	buySelector = &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{btnSupport},
			{btnMainMenu},
		},
	}
	txSelector = &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{btnSupport},
			{btnMainMenu},
		},
	}
	paymentAcceptedSelector = &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{btnSupport},
			{btnMainMenu},
		},
	}
)

// message text templates
const (
	startMessageTemplate = `
<b>Unlock Your Potential with Hanzo Academy!</b>

🏮<b>Join Now, gain enough knowledge to start and insights that will change your life. Some of the public cases:</b>

➱ $100,000 on Starknet.
➱ $100,000 on Tensor.
➱ $20,000 on Backpack.
➱ Called $BALLZ before its 100x surge.
➱ $100,000 on TON live.
➱ Gave away 3 whitelist spots to Entangle, valued at $6,000 each among my private chat members.
➱ and much more

🏮<b>Exclusive Member Benefits:</b>

➱ Connect with like-minded individuals and potential partners
➱ Get direct access to exclusive insider insights.
➱ Get support from me and have your questions answered.
➱ Get access to an educational platform, which will help you to make $100,000+ in 1 year or find a web3 job

🏮<b>Special Offer:</b>

Price: Only $192/year, instead of $1,000

👇<b>Click on the button below and join within 2 clicks</b>👇
`
	startPaidMessageTemplate = `
Hello, <b>%s</b>
You have already purchased the course

Your credentials:
<b>Login</b>: <span class="tg-spoiler">%s</span>
<b>Password</b>: <span class="tg-spoiler">%s</span>
<b>https://www.hanzo.academy/</b>

You can also join the private chat via this <a href="%s">LINK</a>
`
	buyManualTemplate = `
Send <b>%d$</b> in any token to one of those wallets:

<b>TON</b>
<pre><code>%s</code></pre>

<b>SOLANA</b>
<pre><code>%s</code></pre>

<b>TRON</b>
<pre><code>%s</code></pre>

<b>Ethereum/BSC/Polygon</b>
<pre><code>%s</code></pre>

Once the transaction is complete, send tx hash in the following format:
<pre><code>/tx [your tx hash]</code></pre>

<b>Access will be given to you almost instantly!</b>
`

	txFailTemplate = `
You've already submitted tx hash and it is not being reviewed by our support manager, please wait.
If you believe that's a mistake, please contact the support
`

	txSuccessTemplate = `
<b>Thanks for you trust!</b>
Bot will send you all you need to access the course as soon as support manager verifies the payment.
If something is wrong, I will contact you (the tx hash sent by you will be a proof)
`

	paymentNotifTemplate = `
<b>New Payment Submission Received</b>

<b>User</b>: @%s
<b>Tx Hash</b>: <pre><code>%s</code></pre>
`
	paymentAcceptedTemplate = `
<b>Your payment submission has been accepted!</b>

<b>Login</b>: <span class="tg-spoiler">%s</span>
<b>Password</b>: <span class="tg-spoiler">%s</span>
<b>https://www.hanzo.academy/</b>

You can also join the private chat via this <a href="%s">LINK</a>

Save these credentials somewhere
`
	globalErrorTemplate = `
Bot was unable to handle your request.
Please contact the support manager.
`
)

// pieces of text to edit the message
const (
	piecePaymentAccepted = `

<b>Accepted ✅</b>
`
	piecePaymentAcceptanceError = `
	
<b>Unable to accept the payment, contact dev ☎️</b>
`
	pieceUserAlreadyHasPassword = `

<b>User was already accepted ⚠️</b>
`
)

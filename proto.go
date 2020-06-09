package skype

/**
 * all need api url
 *
 */
const (
	API_LOGIN = "https://login.skype.com/login"
	API_MSACC = "https://login.live.com"
	API_USER = "https://api.skype.com"
	API_LOGIN_WEB="https://web.skype.com/login"
	API_PROFILE = "https://profile.skype.com/profile/v1"
	API_OPTIONS = "https://options.skype.com/options/v1/users/self/options"
	API_JOIN = "https://join.skype.com"
	API_JOIN_CREATE = "https://api.join.skype.com/v1"
	API_BOT = "https://api.aps.skype.com/v1"
	API_FLAGS = "https://flagsapi.skype.com/flags/v1"
	API_ENTITLEMENT = "https://consumer.entitlement.skype.com"
	API_TRANSLATE = "https://dev.microsofttranslator.com/api"
	API_ASM = "https://api.asm.skype.com/v1/objects"
	API_ASM_LOCAL = "https://{0}1-api.asm.skype.com/v1/objects"
	API_URL = "https://urlp.asm.skype.com/v1/url/info"
	API_CONTACTS = "https://contacts.skype.com/contacts/v2"
	API_MSGSHOST = "https://client-s.gateway.messenger.live.com"
	API_DIRECTORY = "https://skypegraph.skype.com/search/v1.1/namesearch/swx/"
	API_CONFIG = "https://a.config.skype.com/config/v1"
	API_JOIN_URL = "https://api.scheduler.skype.com/threads"
)

const (
	HTTPS_SCHEME                   string = "https"
	SKYPEWEB_LOCKANDKEY_APPID      string = "msmsgs@msnmsgr.com"
	SKYPEWEB_LOCKANDKEY_SECRET     string = "Q1P7W2E4J9R8U3S5"
	SKYPEWEB_CLIENTINFO_NAME       string = "skype.com"
	SKYPEWEB_CLIENTINFO_VERSION    string = "908/1.30.0.128"
	SKYPEWEB_API_SKYPE_HOST        string = "api.skype.com"
	SKYPEWEB_CONTACTS_HOST         string = "contacts.skype.com"
	SKYPEWEB_DEFAULT_MESSAGES_HOST string = "client-s.gateway.messenger.live.com"

	SKYPEWEB_LOGIN_URL            string = "https://login.skype.com/login?method=skype&client_id=578134&redirect_uri=https%3A%2F%2Fweb.skype.com"
	SKYPEWEB_LOGIN_OAUTH          string = "https://login.skype.com/login/oauth/microsoft"
	SKYPEWEB_SELF_DISPLAYNAME_URL string = "/users/self/displayname"
)
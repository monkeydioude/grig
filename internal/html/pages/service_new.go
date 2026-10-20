package pages

// ServiceNew is the page gathering, in one place, all the config
// a new service needs.
type ServiceNew struct {
	// Forges are the git hosts a repository can be created on.
	Forges []string
	// WebhookURL is where the new repository's webhook will point.
	WebhookURL string
	// WebhookErr tells why no webhook can be added, if so.
	WebhookErr error
	// DatabaseURL is the admin DB_URL, password redacted. Empty when unset.
	DatabaseURL string
	// Password prefills the new database user's password.
	Password string
}

func NewServiceNew(forges []string, webhookURL string, webhookErr error) ServiceNew {
	return ServiceNew{Forges: forges, WebhookURL: webhookURL, WebhookErr: webhookErr}
}

func (p ServiceNew) WithDatabase(redactedURL, password string) ServiceNew {
	p.DatabaseURL = redactedURL
	p.Password = password
	return p
}

func (ServiceNew) Title() string {
	return "New service"
}

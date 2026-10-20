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
	// EnvBasePath is where the project's .env is written. Empty when unset.
	EnvBasePath string
	// ComposeTypes are the docker-compose.<type>.yml files a service can go in.
	ComposeTypes []string
	// ComposeErr tells why they could not be listed.
	ComposeErr error
}

func NewServiceNew(forges []string, webhookURL string, webhookErr error) ServiceNew {
	return ServiceNew{Forges: forges, WebhookURL: webhookURL, WebhookErr: webhookErr}
}

func (p ServiceNew) WithDatabase(redactedURL, password, envBasePath string) ServiceNew {
	p.DatabaseURL = redactedURL
	p.Password = password
	p.EnvBasePath = envBasePath
	return p
}

func (p ServiceNew) WithCompose(types []string, err error) ServiceNew {
	p.ComposeTypes = types
	p.ComposeErr = err
	return p
}

func (p ServiceNew) composeUnavailable() string {
	if p.ComposeErr != nil {
		return "Could not list the docker compose files: " + p.ComposeErr.Error()
	}
	return "No docker compose file: set docker_compose_base_path in grig's config, to a folder of docker-compose.<type>.yml files."
}

func (ServiceNew) Title() string {
	return "New service"
}

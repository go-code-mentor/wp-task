package app

func New(cfg Config) *App {
	return &App{cfg: cfg}
}

type App struct {
	cfg Config
}

func (a *App) Build() error {

	return nil
}

func (a *App) Run() error {

	return nil
}

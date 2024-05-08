package main

import (
	"context"
	"mekano-react-wails/backend/mekano"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) onDomReady(ctx context.Context) {
	//Check if API is alive
	_, err := http.Get(os.Getenv("MONGO_SERVER"))
	if err != nil {
		ok, _ := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Message:       "Error al realizar la peticion al servidor. Compruebe su estado",
			DefaultButton: "Cerrar",
			CancelButton:  "Cerrar",
		})
		if ok == "OK" {
			runtime.Quit(ctx)
		}
	}

}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OpenFile() string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{DisplayName: "*.xlsx", Pattern: "*.xlsx"},
		},
	})

	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Message: err.Error(),
		})
	}

	return path
}

func (a *App) MekanoPayment(path string) mekano.PaymentStatistics {
	database, err := mekano.NewDatabase()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Message: err.Error(),
		})
	}

	mekano := mekano.NewMekano(database)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Message: err.Error(),
		})
	}
	data, err := mekano.ProcessPaymentFile(path)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Message: err.Error(),
		})
	}

	return data
}

func (a *App) MekanoBilling(path string, extrasPath string) mekano.BillingStatistics {
	database, err := mekano.NewDatabase()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Message: err.Error(),
		})
	}

	mekano := mekano.NewMekano(database)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Message: err.Error(),
		})
	}

	data, err := mekano.ProcessBillFile(path, extrasPath)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Message: err.Error(),
		})
	}

	return data
}

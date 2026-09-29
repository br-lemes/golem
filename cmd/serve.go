package cmd

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/br-lemes/golem/assets"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/br-lemes/golem/ui/components/navigation"
	_ "github.com/br-lemes/golem/ui/pages"
	"github.com/br-lemes/golem/ui/routes"
	"github.com/spf13/cobra"
)

type serveFlags struct {
	Host string `flag:"host" default:"127.0.0.1" desc:"Host address to listen on"`
	Port int    `flag:"port" default:"8080" desc:"Port to listen on"`
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		navigation.SetVersion(rootCmd.Version)
		flags, err := utils.ReadFlags[serveFlags](cmd)
		if err != nil {
			return err
		}
		if flags.Port < 1 || flags.Port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535: %d", flags.Port)
		}
		mux := http.NewServeMux()
		setupAssetsRoutes(mux)
		routes.RegisterRoutes(mux)
		address := net.JoinHostPort(flags.Host, fmt.Sprintf("%d", flags.Port))
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Server is running on http://%s\n", address)
		return http.ListenAndServe(address, mux)
	},
}

func setupAssetsRoutes(mux *http.ServeMux) {
	development := !strings.HasPrefix(rootCmd.Version, "v")
	files := http.FileServer(http.FS(assets.Assets))
	if development {
		files = http.FileServer(http.Dir("./assets"))
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if development {
			w.Header().Set("Cache-Control", "no-store")
		} else if strings.HasPrefix(strings.TrimPrefix(r.URL.Path, "/"), "js/shadcn-templ-") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		files.ServeHTTP(w, r)
	})
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", handler))
}

func init() {
	rootCmd.AddCommand(serveCmd)
	err := utils.RegisterFlags[serveFlags](serveCmd)
	if err != nil {
		panic(err)
	}
}

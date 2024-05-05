package kitty

// [
//
//	{
//	  "background_opacity": 1.0,
//	  "id": 1,
//	  "is_active": true,
//	  "is_focused": true,
//	  "last_focused": true,
//	  "platform_window_id": 41222,
//	  "tabs": [
//	    {
//	      "active_window_history": [
//	        1
//	      ],
//	      "enabled_layouts": [
//	        "fat",
//	        "grid",
//	        "horizontal",
//	        "splits",
//	        "stack",
//	        "tall",
//	        "vertical"
//	      ],
//	      "groups": [
//	        {
//	          "id": 1,
//	          "windows": [
//	            1
//	          ]
//	        }
//	      ],
//	      "id": 1,
//	      "is_active": true,
//	      "is_focused": true,
//	      "layout": "fat",
//	      "layout_opts": {
//	        "bias": 50,
//	        "full_size": 1,
//	        "mirrored": false
//	      },
//	      "layout_state": {
//	        "biased_map": {},
//	        "main_bias": [
//	          0.5,
//	          0.5
//	        ],
//	        "num_full_size_windows": 1
//	      },
//	      "title": "kitten @ ls ~/l/tix",
//	      "windows": [
//	        {
//	          "at_prompt": false,
//	          "cmdline": [
//	            "/opt/homebrew/bin/fish"
//	          ],
//	          "columns": 144,
//	          "created_at": 1714073325124355000,
//	          "cwd": "/Users/agravelot/lab/tix/\n0\t\u0003\u0001",
//	          "env": {
//	            "COLORTERM": "truecolor",
//	            "COMMAND_MODE": "unix2003",
//	            "EDITOR": "nvim",
//	            "HOME": "/Users/agravelot",
//	            "HOMEBREW_CELLAR": "/opt/homebrew/Cellar",
//	            "HOMEBREW_PREFIX": "/opt/homebrew",
//	            "HOMEBREW_REPOSITORY": "/opt/homebrew",
//	            "INFOPATH": "/opt/homebrew/share/info:",
//	            "KITTY_INSTALLATION_DIR": "/Applications/kitty.app/Contents/Resources/kitty",
//	            "KITTY_PID": "16775",
//	            "KITTY_PUBLIC_KEY": "1:&P$iAjW7TCNI=@iP6~aqPOgkWK6hnK3p&ZLZ#J)7",
//	            "KITTY_WINDOW_ID": "1",
//	            "LANG": "en_US.UTF-8",
//	            "LOGNAME": "agravelot",
//	            "LaunchInstanceID": "5CEB4B31-4505-4CF3-9D47-736A3D112E54",
//	            "MANPATH": "/Applications/kitty.app/Contents/Resources/man:/opt/homebrew/share/man:/usr/share/man:/usr/local/share/man:/Applications/kitty.app/Contents/Resources/man:",
//	            "PATH": "/opt/homebrew/sbin:/Users/agravelot/.orbstack/bin:/Users/agravelot/.volta/bin:/Users/agravelot/go/bin:/opt/homebrew/Cellar/go/1.22.2/libexec/bin:/opt/homebrew/bin:/usr/local/bin:/System/Cryptexes/App/usr/bin:/usr/bin:/bin:/usr/sbin:/sbin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/appleinternal/bin:/Library/Apple/usr/bin:/Applications/kitty.app/Contents/MacOS",
//	            "PWD": "/Users/agravelot/lab/tix",
//	            "SECURITYSESSIONID": "186a2",
//	            "SHELL": "/opt/homebrew/bin/fish",
//	            "SHLVL": "1",
//	            "SSH_AUTH_SOCK": "/private/tmp/com.apple.launchd.hJWqf3asPZ/Listeners",
//	            "STARSHIP_SESSION_KEY": "2854870651306262",
//	            "STARSHIP_SHELL": "fish",
//	            "TERM": "xterm-kitty",
//	            "TERMINFO": "/Applications/kitty.app/Contents/Resources/kitty/terminfo",
//	            "TMPDIR": "/var/folders/59/tlgd2lp913z60gkgz1nrbcvm0000gp/T/",
//	            "USER": "agravelot",
//	            "VOLTA_HOME": "/Users/agravelot/.volta",
//	            "WINDOWID": "41222",
//	            "XPC_FLAGS": "0x0",
//	            "XPC_SERVICE_NAME": "0",
//	            "__CFBundleIdentifier": "net.kovidgoyal.kitty",
//	            "__CF_USER_TEXT_ENCODING": "0x1F6:0x0:0x0",
//	            "_tide_color_separator_same_color": "\u001b[38;2;148;148;148m",
//	            "_tide_location_color": "\u001b[38;2;95;215;0m"
//	          },
//	          "foreground_processes": [
//	            {
//	              "cmdline": [
//	                "/Applications/kitty.app/Contents/MacOS/kitten",
//	                "@",
//	                "ls"
//	              ],
//	              "cwd": "/Users/agravelot/lab/tix",
//	              "pid": 18040
//	            }
//	          ],
//	          "id": 1,
//	          "is_active": true,
//	          "is_focused": true,
//	          "is_self": true,
//	          "lines": 71,
//	          "pid": 16778,
//	          "title": "kitten @ ls ~/l/tix",
//	          "user_vars": {}
//	        }
//	      ]
//	    }
//	  ],
//	  "wm_class": "kitty",
//	  "wm_name": "kitty"
//	}
//
// ]
type Window struct {
	BackgroundOpacity float64 `json:"background_opacity,omitempty"`
	ID                float64 `json:"id,omitempty"`
	IsActive          bool    `json:"is_active,omitempty"`
	IsFocused         bool    `json:"is_focused,omitempty"`
	LastFocused       bool    `json:"last_focused,omitempty"`
	PlatformWindowID  float64 `json:"platform_window_id,omitempty"`
	Tabs              []struct {
		ActiveWindowHistory []float64 `json:"active_window_history,omitempty"`
		EnabledLayouts      []string  `json:"enabled_layouts,omitempty"`
		Groups              []struct {
			ID      float64   `json:"id,omitempty"`
			Windows []float64 `json:"windows,omitempty"`
		} `json:"groups,omitempty"`
		ID         float64 `json:"id,omitempty"`
		IsActive   bool    `json:"is_active,omitempty"`
		IsFocused  bool    `json:"is_focused,omitempty"`
		Layout     string  `json:"layout,omitempty"`
		LayoutOpts struct {
			Bias     float64 `json:"bias,omitempty"`
			FullSize float64 `json:"full_size,omitempty"`
			Mirrored bool    `json:"mirrored,omitempty"`
		} `json:"layout_opts,omitempty"`
		LayoutState struct {
			BiasedMap          struct{}  `json:"biased_map,omitempty"`
			MainBias           []float64 `json:"main_bias,omitempty"`
			NumFullSizeWindows float64   `json:"num_full_size_windows,omitempty"`
		} `json:"layout_state,omitempty"`
		Title   string `json:"title,omitempty"`
		Windows []struct {
			AtPrompt  bool     `json:"at_prompt,omitempty"`
			Cmdline   []string `json:"cmdline,omitempty"`
			Columns   float64  `json:"columns,omitempty"`
			CreatedAt float64  `json:"created_at,omitempty"`
			Cwd       string   `json:"cwd,omitempty"`
			Env       struct {
				COLORTERM                   string
				COMMANDMODE                 string `json:"COMMAND_MODE,omitempty"`
				EDITOR                      string
				HOME                        string
				HOMEBREWCELLAR              string `json:"HOMEBREW_CELLAR,omitempty"`
				HOMEBREWPREFIX              string `json:"HOMEBREW_PREFIX,omitempty"`
				HOMEBREWREPOSITORY          string `json:"HOMEBREW_REPOSITORY,omitempty"`
				INFOPATH                    string
				KITTYINSTALLATIONDIR        string `json:"KITTY_INSTALLATION_DIR,omitempty"`
				KITTYPID                    string `json:"KITTY_PID,omitempty"`
				KITTYPUBLICKEY              string `json:"KITTY_PUBLIC_KEY,omitempty"`
				KITTYWINDOWID               string `json:"KITTY_WINDOW_ID,omitempty"`
				LANG                        string
				LOGNAME                     string
				LaunchInstanceID            string
				MANPATH                     string
				PATH                        string
				PWD                         string
				SECURITYSESSIONID           string
				SHELL                       string
				SHLVL                       string
				SSHAUTHSOCK                 string `json:"SSH_AUTH_SOCK,omitempty"`
				STARSHIPSESSIONKEY          string `json:"STARSHIP_SESSION_KEY,omitempty"`
				STARSHIPSHELL               string `json:"STARSHIP_SHELL,omitempty"`
				TERM                        string
				TERMINFO                    string
				TMPDIR                      string
				USER                        string
				VOLTAHOME                   string `json:"VOLTA_HOME,omitempty"`
				WINDOWID                    string
				XPCFLAGS                    string `json:"XPC_FLAGS,omitempty"`
				XPCSERVICENAME              string `json:"XPC_SERVICE_NAME,omitempty"`
				CFBundleIdentifier          string `json:"__CFBundleIdentifier,omitempty"`
				CFUSERTEXTENCODING          string `json:"__CF_USER_TEXT_ENCODING,omitempty"`
				TideColorSeparatorSameColor string `json:"_tide_color_separator_same_color,omitempty"`
				TideLocationColor           string `json:"_tide_location_color,omitempty"`
			} `json:"env,omitempty"`
			ForegroundProcesses []struct {
				Cmdline []string `json:"cmdline,omitempty"`
				Cwd     string   `json:"cwd,omitempty"`
				Pid     float64  `json:"pid,omitempty"`
			} `json:"foreground_processes,omitempty"`
			ID        float64  `json:"id,omitempty"`
			IsActive  bool     `json:"is_active,omitempty"`
			IsFocused bool     `json:"is_focused,omitempty"`
			IsSelf    bool     `json:"is_self,omitempty"`
			Lines     float64  `json:"lines,omitempty"`
			Pid       float64  `json:"pid,omitempty"`
			Title     string   `json:"title,omitempty"`
			UserVars  struct{} `json:"user_vars,omitempty"`
		} `json:"windows,omitempty"`
	} `json:"tabs,omitempty"`
	WmClass string `json:"wm_class,omitempty"`
	WmName  string `json:"wm_name,omitempty"`
}

package config

import "time"

type (
	Config struct {
		Application Application
		Server      Server
		Cron        Cron
		Gymstation  Gymstation
	}

	Application struct {
		Development bool
		LogLevel    string `yaml:"log_level"`
	}

	Server struct {
		API   *HTTP
		Debug *HTTP
	}

	HTTP struct {
		Host    string
		Port    string
		Network string

		ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
		ReadTimeout       time.Duration `yaml:"read_timeout"`
		WriteTimeout      time.Duration `yaml:"write_timeout"`
		IdleTimeout       time.Duration `yaml:"idle_timeout"`

		MaxHeaderBytes      int `yaml:"max_header_bytes"`
		MaxHeaderValueCount int `yaml:"max_header_value_count"`

		DisableGeneralOptionsHandler bool `yaml:"disable_general_options_handler"`
		DisableKeepAlives            bool `yaml:"disable_keep_alives"`
		DisableClientPriority        bool `yaml:"disable_client_priority"`

		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`

		Protocols    *HTTPProtocols
		HTTP2        *HTTP2
		TCPKeepAlive *TCPKeepAlive `yaml:"tcp_keep_alive"`
	}

	HTTPProtocols struct {
		HTTP1            bool
		HTTP2            bool
		UnencryptedHTTP2 bool `yaml:"unencrypted_http2"`
	}

	HTTP2 struct {
		MaxConcurrentStreams          int           `yaml:"max_concurrent_streams"`
		MaxDecoderHeaderTableSize     int           `yaml:"max_decoder_header_table_size"`
		MaxEncoderHeaderTableSize     int           `yaml:"max_encoder_header_table_size"`
		MaxReadFrameSize              int           `yaml:"max_read_frame_size"`
		MaxReceiveBufferPerConnection int           `yaml:"max_receive_buffer_per_connection"`
		MaxReceiveBufferPerStream     int           `yaml:"max_receive_buffer_per_stream"`
		SendPingTimeout               time.Duration `yaml:"send_ping_timeout"`
		PingTimeout                   time.Duration `yaml:"ping_timeout"`
		WriteByteTimeout              time.Duration `yaml:"write_byte_timeout"`
	}

	TCPKeepAlive struct {
		Enable   bool
		Idle     time.Duration
		Interval time.Duration
		Count    int
	}

	Cron struct {
		Location    string
		StopTimeout time.Duration `yaml:"stop_timeout"`
		Jobs        []*CronJob
	}

	CronJob struct {
		Name string
		Mask string
	}

	Gymstation struct {
		Host                  string
		Timeout               time.Duration `yaml:"timeout"`
		DialTimeout           time.Duration `yaml:"dial_timeout"`
		KeepAlive             time.Duration `yaml:"keep_alive"`
		TLSHandshakeTimeout   time.Duration `yaml:"tls_handshake_timeout"`
		ResponseHeaderTimeout time.Duration `yaml:"response_header_timeout"`
		MaxIdleConns          int           `yaml:"max_idle_conns"`
		IdleConnTimeout       time.Duration `yaml:"idle_conn_timeout"`
		ExpectContinueTimeout time.Duration `yaml:"expect_continue_timeout"`
		ForceAttemptHTTP2     bool          `yaml:"force_attempt_http2"`
		UserAgent             string        `yaml:"user_agent"`
	}
)
